package opengauss

import (
	"strings"

	"mw-sre-platform/internal/core"
	"mw-sre-platform/internal/executor"
	base "mw-sre-platform/internal/plugins"
)

type Plugin struct{}

func (Plugin) Name() string      { return "opengauss" }
func (Plugin) Aliases() []string { return []string{"og", "open-gauss", "gauss", "gaussdb"} }

func (Plugin) Commands(cfg base.RuntimeConfig) []executor.Command {
	env := map[string]string{}
	if cfg.OpenGaussDataDir != "" {
		env["OPENGAUSS_DATA_DIR"] = cfg.OpenGaussDataDir
	}
	return []executor.Command{
		{
			Name: "opengauss.gsctl-query",
			Args: []string{"bash", "-lc", `set +e
found=""
for d in "$OPENGAUSS_DATA_DIR" "$RUNTIME_ROOTFS/opt/data/opengauss/data/dn" "$HOME/opt/data/opengauss/data/dn" "/opt/data/opengauss/data/dn"; do
  [ -n "$d" ] && [ -d "$d" ] && found="$d" && break
done
if [ -z "$found" ]; then
  echo "OPENGAUSS_DATA_DIR_NOT_FOUND"
  exit 0
fi
echo "DATA_DIR=$found"
if command -v gs_ctl >/dev/null 2>&1; then
  gs_ctl query -D "$found" 2>&1 || true
else
  echo "GS_CTL_NOT_FOUND"
fi`},
			Env: env,
		},
		{
			Name: "opengauss.process-port",
			Args: []string{"bash", "-lc", `pgrep -a -f 'gaussdb|gs_ctl|gsql' 2>/dev/null | grep -v -E 'mwctl check --module opengauss|bash -lc' || true; (ss -lntup 2>/dev/null || netstat -ntlp 2>/dev/null || true) | grep -E '(:5432|:26000|:15432|gauss)' || true`},
		},
		{
			Name: "opengauss.log-tail",
			Args: []string{"bash", "-lc", `for d in "$RUNTIME_ROOTFS/opt/data/opengauss" "$HOME/opt/data/opengauss" "/opt/data/opengauss"; do [ -d "$d" ] && find "$d" -maxdepth 4 -type f \( -name '*.log' -o -name 'postgresql-*' \) -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -3 | cut -d' ' -f2-; done | while read -r f; do echo "FILE=$f"; tail -n 40 "$f" 2>/dev/null | grep -Ei 'panic|fatal|error|abnormal|fail|denied|permission|could not|No space|disk' | tail -20; done | head -160 || true`},
		},
	}
}

func (Plugin) Findings(commands []core.CommandResult) []core.Finding {
	commands = base.FilterCommands(commands, "opengauss")
	joined := base.JoinEvidence(commands)
	lower := strings.ToLower(joined)
	var out []core.Finding

	if strings.Contains(joined, "OPENGAUSS_DATA_DIR_NOT_FOUND") {
		out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.data_dir_not_found", Severity: core.SeverityWarning, Summary: "OpenGauss data directory was not found in common paths", Evidence: []string{"OPENGAUSS_DATA_DIR_NOT_FOUND"}, Recommendation: "Set OPENGAUSS_DATA_DIR or verify root/non-root path prefix before judging database state."})
	}
	if strings.Contains(joined, "GS_CTL_NOT_FOUND") {
		out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.gs_ctl_missing", Severity: core.SeverityWarning, Summary: "gs_ctl command is not available in current shell", Evidence: []string{"GS_CTL_NOT_FOUND"}, Recommendation: "Source the product/openGauss environment profile or switch to the database runtime user before running HA checks."})
	}
	if strings.Contains(lower, "no space left") || strings.Contains(lower, "disk full") {
		out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.disk_full", Severity: core.SeverityCritical, Summary: "OpenGauss evidence contains disk full symptoms", Evidence: []string{"log contains no space left/disk full"}, Recommendation: "Stop recovery attempts that may worsen data risk; confirm disk usage and clean only approved log/temp files."})
	}
	if strings.Contains(lower, "permission denied") {
		out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.permission_denied", Severity: core.SeverityCritical, Summary: "OpenGauss evidence contains permission denied", Evidence: []string{"log or gs_ctl output contains permission denied"}, Recommendation: "Check data directory owner, mode and runtime user; do not chown recursively without backup and approval."})
	}
	if strings.Contains(joined, "HA state:") {
		if !containsFieldValue(joined, "db_state", "Normal") || !containsFieldValue(joined, "detail_information", "Normal") {
			out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.ha_not_normal", Severity: core.SeverityCritical, Summary: "OpenGauss HA state is not Normal", Evidence: []string{"gs_ctl query did not show db_state/detail_information Normal"}, Recommendation: "Inspect sender/receiver state and database logs before restarting upper-layer services."})
		}
		if containsFieldValue(joined, "local_role", "Standby") && strings.Contains(joined, "Receiver info:\nNo information") {
			out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.standby_receiver_missing", Severity: core.SeverityCritical, Summary: "OpenGauss standby has no receiver information", Evidence: []string{"local_role Standby with Receiver info: No information"}, Recommendation: "Check primary connectivity, replication slot/sender state and HA configuration from both primary and standby."})
		}
	}
	if strings.Contains(lower, "could not start server") || strings.Contains(lower, "database system is starting up") {
		out = append(out, core.Finding{Service: "opengauss", RuleID: "opengauss.startup_unhealthy", Severity: core.SeverityCritical, Summary: "OpenGauss startup is not healthy", Evidence: []string{"startup output contains could not start server or still starting"}, Recommendation: "Collect recent database logs and confirm disk/permission/environment before retrying service start."})
	}
	return out
}

func containsFieldValue(s, field, value string) bool {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, field) && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[1]) == value {
				return true
			}
		}
	}
	return false
}
