package checks

import (
	"testing"

	"mw-sre-platform/internal/core"
)

func TestRedisBothMasterFinding(t *testing.T) {
	cmds := []core.CommandResult{{Stdout: "role:master\nconnected_slaves:0"}, {Stdout: "role:master\nconnected_slaves:0"}}
	findings := RedisFindings(cmds)
	if len(findings) != 1 || findings[0].RuleID != "redis.both_master" {
		t.Fatalf("unexpected findings: %#v", findings)
	}
}

func TestMongoReplSetFinding(t *testing.T) {
	cmds := []core.CommandResult{{Stdout: "lastHeartbeatMessage: not running with --replSet"}}
	findings := MongoFindings(cmds)
	if len(findings) != 1 || findings[0].Severity != core.SeverityCritical {
		t.Fatalf("unexpected findings: %#v", findings)
	}
}
