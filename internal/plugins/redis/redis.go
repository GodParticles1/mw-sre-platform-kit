package redis

import (
	"fmt"
	"strings"

	"mw-sre-platform/internal/core"
	"mw-sre-platform/internal/executor"
	base "mw-sre-platform/internal/plugins"
)

type Plugin struct{}

func (Plugin) Name() string      { return "redis" }
func (Plugin) Aliases() []string { return []string{"redis-replication", "redis-ha"} }

func (Plugin) Commands(cfg base.RuntimeConfig) []executor.Command {
	env := map[string]string{}
	authPrefix := ""
	authOpt := ""
	if cfg.RedisPass != "" {
		env["REDISCLI_AUTH"] = cfg.RedisPass
		authPrefix = "REDISCLI_AUTH is set; "
		authOpt = " --no-auth-warning"
	}
	return []executor.Command{
		{
			Name: "redis.info-replication",
			Args: []string{"bash", "-lc", fmt.Sprintf(`echo %q; redis-cli -h 127.0.0.1 -p 6379%s INFO replication 2>&1 || true`, authPrefix+"collecting local Redis replication", authOpt)},
			Env:  env,
		},
		{
			Name: "redis.process-port",
			Args: []string{"bash", "-lc", `pgrep -a -f 'redis-server|redis-cli' 2>/dev/null | grep -v -E 'mwctl check --module redis|bash -lc' || true; (ss -lntup 2>/dev/null || netstat -ntlp 2>/dev/null || true) | grep -E '(:6379|redis)' || true`},
		},
		{
			Name: "redis.config-summary",
			Args: []string{"bash", "-lc", `for f in /etc/redis.conf /etc/redis/redis.conf /opt/data/redis/redis.conf "$HOME"/opt/data/redis/redis.conf; do [ -f "$f" ] && echo "FILE=$f" && grep -nE '^[[:space:]]*(port|bind|requirepass|masterauth|replicaof|slaveof)[[:space:]]+' "$f" | sed -E 's#(requirepass|masterauth)[[:space:]]+.*#\1 ***#'; done 2>/dev/null | head -120 || true`},
		},
	}
}

func (Plugin) Findings(commands []core.CommandResult) []core.Finding {
	commands = base.FilterCommands(commands, "redis")
	joined := base.JoinEvidence(commands)
	lower := strings.ToLower(joined)
	var out []core.Finding

	if strings.Contains(joined, "NOAUTH") || strings.Contains(lower, "authentication required") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.noauth", Severity: core.SeverityWarning, Summary: "Redis requires authentication", Evidence: []string{"INFO replication returned NOAUTH/authentication required"}, Recommendation: "Re-run with REDIS_PASS or source the product Redis password from the approved profile secret."})
	}
	if strings.Contains(lower, "redis-cli: command not found") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.cli_missing", Severity: core.SeverityWarning, Summary: "redis-cli is not available in current shell", Evidence: []string{"redis-cli: command not found"}, Recommendation: "Source the product profile or run from the Redis runtime environment before judging Redis health."})
	}
	if strings.Contains(lower, "connection refused") || strings.Contains(lower, "could not connect to redis") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.port_down", Severity: core.SeverityCritical, Summary: "Redis local port is not reachable", Evidence: []string{"redis-cli could not connect to 127.0.0.1:6379"}, Recommendation: "Check redis process, monit status and local bind/port before changing replication."})
	}
	if strings.Count(joined, "role:master") >= 2 && !strings.Contains(joined, "connected_slaves:1") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.both_master", Severity: core.SeverityCritical, Summary: "More than one Redis evidence block appears to be master", Evidence: []string{"role:master appears more than once without connected slave evidence"}, Recommendation: "Confirm VIP owner and physical master/slave roles, then use the Redis runbook in dry-run mode before applying replica changes."})
	}
	if strings.Contains(joined, "role:slave") && strings.Contains(joined, "master_link_status:down") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.slave_link_down", Severity: core.SeverityCritical, Summary: "Redis replica link is down", Evidence: []string{"role:slave with master_link_status:down"}, Recommendation: "Check TCP 6379 from replica to master, masterauth/requirepass consistency, current master_host and VIP owner."})
	}
	if strings.Contains(joined, "masterauth ***") && !strings.Contains(joined, "requirepass ***") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.auth_asymmetric", Severity: core.SeverityWarning, Summary: "Redis config shows masterauth without visible requirepass", Evidence: []string{"config summary contains masterauth but not requirepass"}, Recommendation: "Verify the effective Redis password policy; do not change config until roles and security requirements are confirmed."})
	}
	return out
}
