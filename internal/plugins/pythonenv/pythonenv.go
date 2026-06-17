package pythonenv

import (
	"strings"

	"mw-sre-platform/internal/core"
	"mw-sre-platform/internal/executor"
	base "mw-sre-platform/internal/plugins"
)

type Plugin struct{}

func (Plugin) Name() string      { return "python-env" }
func (Plugin) Aliases() []string { return []string{"python", "pythonenv", "profile", "env"} }

func (Plugin) Commands(cfg base.RuntimeConfig) []executor.Command {
	return []executor.Command{
		{
			Name: "python-env.profile-source",
			Args: []string{"bash", "-lc", `set +e
out="$({ source /etc/profile; } 2>&1)"
rc=$?
echo "PROFILE_SOURCE_EXIT=$rc"
if [ -n "$out" ]; then
  echo "PROFILE_SOURCE_OUTPUT_BEGIN"
  printf '%s\n' "$out"
  echo "PROFILE_SOURCE_OUTPUT_END"
fi
command -v python3 >/dev/null 2>&1 && echo "PYTHON3=$(command -v python3)" || echo "PYTHON3_NOT_FOUND"
python3 -V 2>&1 || true
printf 'ENV_SELECTED_BEGIN\n'
env | grep -E '^(PATH|PYTHONPATH|LD_LIBRARY_PATH|LANG|LC_|HOME|RUNTIME_ROOTFS)=' | sort
printf 'ENV_SELECTED_END\n'`},
		},
		{
			Name: "python-env.profile-risk-lines",
			Args: []string{"bash", "-lc", `for f in /etc/profile /etc/bashrc "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile"; do
  [ -f "$f" ] || continue
  echo "FILE=$f"
  grep -nE '^[[:space:]]*(TIMEOUT|TMOUT|PYTHONPATH|PATH|LD_LIBRARY_PATH|ulimit|echo|printf|return|exit)[[:space:]=]' "$f" 2>/dev/null | head -40
done | head -160 || true`},
		},
		{
			Name: "python-env.import-smoke",
			Args: []string{"bash", "-lc", `python3 - <<'PY' 2>&1 || true
import os, sys
print('PYTHON_EXECUTABLE=' + sys.executable)
print('PYTHON_VERSION=' + sys.version.replace('\n', ' '))
print('SYS_PATH_LEN=' + str(len(sys.path)))
for name in ('ssl', 'sqlite3', 'ctypes'):
    try:
        __import__(name)
        print('IMPORT_OK=' + name)
    except Exception as exc:
        print('IMPORT_FAIL=' + name + ':' + type(exc).__name__ + ':' + str(exc))
PY`},
		},
	}
}

func (Plugin) Findings(commands []core.CommandResult) []core.Finding {
	commands = base.FilterCommands(commands, "python-env")
	joined := base.JoinEvidence(commands)
	lower := strings.ToLower(joined)
	var out []core.Finding

	if strings.Contains(joined, "PROFILE_SOURCE_EXIT=") && !strings.Contains(joined, "PROFILE_SOURCE_EXIT=0") {
		out = append(out, core.Finding{Service: "python-env", RuleID: "python_env.profile_source_nonzero", Severity: core.SeverityCritical, Summary: "/etc/profile returns non-zero when sourced", Evidence: []string{"PROFILE_SOURCE_EXIT is not 0"}, Recommendation: "Fix /etc/profile or the sourced file that returns non-zero before installing drivers or starting Console scripts."})
	}
	if hasProfileOutput(joined) {
		out = append(out, core.Finding{Service: "python-env", RuleID: "python_env.profile_source_output", Severity: core.SeverityWarning, Summary: "/etc/profile prints output when sourced", Evidence: []string{"PROFILE_SOURCE_OUTPUT_BEGIN/END block is not empty"}, Recommendation: "Remove or guard echo/printf/banner output from /etc/profile and sourced scripts; product installers expect a quiet source."})
	}
	if strings.Contains(joined, "PYTHON3_NOT_FOUND") {
		out = append(out, core.Finding{Service: "python-env", RuleID: "python_env.python3_missing", Severity: core.SeverityWarning, Summary: "python3 command is not available", Evidence: []string{"PYTHON3_NOT_FOUND"}, Recommendation: "Confirm the expected Python runtime for this baseline and source the correct product profile before running Python-based tools."})
	}
	if strings.Contains(joined, "IMPORT_FAIL=") {
		out = append(out, core.Finding{Service: "python-env", RuleID: "python_env.import_fail", Severity: core.SeverityWarning, Summary: "Python standard-library import smoke test failed", Evidence: importFailures(joined), Recommendation: "Check PYTHONPATH/LD_LIBRARY_PATH pollution and missing OS shared libraries before blaming the business script."})
	}
	if strings.Contains(lower, "time") && strings.Contains(joined, "TIMEOUT=") && !strings.Contains(joined, "export TIMEOUT") {
		out = append(out, core.Finding{Service: "python-env", RuleID: "python_env.timeout_unexported", Severity: core.SeverityInfo, Summary: "Profile contains TIMEOUT assignment", Evidence: []string{"profile risk lines contain TIMEOUT="}, Recommendation: "Some security baselines add TIMEOUT/TMOUT; confirm whether it breaks source /etc/profile and export semantics before changing it."})
	}
	if strings.Contains(joined, "PYTHONPATH=") || strings.Contains(joined, "LD_LIBRARY_PATH=") {
		out = append(out, core.Finding{Service: "python-env", RuleID: "python_env.path_pollution_possible", Severity: core.SeverityInfo, Summary: "Python-related environment variables are set", Evidence: []string{"PYTHONPATH or LD_LIBRARY_PATH is present"}, Recommendation: "If Python/driver behavior is abnormal, compare these variables with a healthy node before modifying them."})
	}
	return out
}

func hasProfileOutput(s string) bool {
	start := strings.Index(s, "PROFILE_SOURCE_OUTPUT_BEGIN")
	end := strings.Index(s, "PROFILE_SOURCE_OUTPUT_END")
	if start < 0 || end < 0 || end <= start {
		return false
	}
	content := s[start+len("PROFILE_SOURCE_OUTPUT_BEGIN") : end]
	return strings.TrimSpace(content) != ""
}

func importFailures(s string) []string {
	var ev []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "IMPORT_FAIL=") {
			ev = append(ev, strings.TrimSpace(line))
		}
	}
	if len(ev) == 0 {
		return []string{"IMPORT_FAIL detected"}
	}
	return ev
}
