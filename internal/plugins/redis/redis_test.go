package redis

import (
	"testing"

	"mw-sre-platform/internal/core"
)

func TestRedisSlaveLinkDownFinding(t *testing.T) {
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "redis.info-replication", Stdout: "role:slave\nmaster_link_status:down\n"}})
	assertRule(t, findings, "redis.slave_link_down")
}

func TestRedisPortDownFinding(t *testing.T) {
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "redis.info-replication", Stdout: "Could not connect to Redis at 127.0.0.1:6379: Connection refused"}})
	assertRule(t, findings, "redis.port_down")
}

func assertRule(t *testing.T, findings []core.Finding, rule string) {
	t.Helper()
	for _, f := range findings {
		if f.RuleID == rule {
			return
		}
	}
	t.Fatalf("rule %s not found in %#v", rule, findings)
}
