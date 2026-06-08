package checks

import (
	"strings"

	"mw-sre-platform/internal/core"
)

func RedisFindings(commands []core.CommandResult) []core.Finding {
	var out []core.Finding
	joined := join(commands)
	if strings.Contains(joined, "NOAUTH") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.noauth", Severity: core.SeverityWarning, Summary: "Redis requires authentication", Evidence: []string{"INFO replication returned NOAUTH"}, Recommendation: "Re-run with the correct Redis password or source it from the profile secret."})
	}
	if strings.Count(joined, "role:master") >= 2 && !strings.Contains(joined, "connected_slaves:1") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.both_master", Severity: core.SeverityCritical, Summary: "Both Redis nodes appear to be master", Evidence: []string{"role:master appears on more than one node"}, Recommendation: "Make the non-VIP node replicate from the VIP-holder physical IP after confirmation."})
	}
	if strings.Contains(joined, "master_link_status:down") {
		out = append(out, core.Finding{Service: "redis", RuleID: "redis.slave_link_down", Severity: core.SeverityCritical, Summary: "Redis slave link is down", Evidence: []string{"master_link_status:down"}, Recommendation: "Check TCP 6379, masterauth, requirepass and current master host."})
	}
	return out
}

func MongoFindings(commands []core.CommandResult) []core.Finding {
	joined := join(commands)
	if strings.Contains(joined, "not running with --replSet") {
		return []core.Finding{{Service: "mongo", RuleID: "mongo.not_running_with_replset", Severity: core.SeverityCritical, Summary: "Mongo secondary is not running with replSet enabled", Evidence: []string{"lastHeartbeatMessage contains not running with --replSet"}, Recommendation: "Patch mongod.conf replication.replSetName to the primary replica set name, then restart mongod."}}
	}
	return nil
}

func MySQLFindings(commands []core.CommandResult) []core.Finding {
	joined := join(commands)
	var out []core.Finding
	if strings.Contains(joined, "equal MySQL server ids") {
		out = append(out, core.Finding{Service: "mysql", RuleID: "mysql.equal_server_id", Severity: core.SeverityCritical, Summary: "MySQL replication has equal server_id on both nodes", Evidence: []string{"Last_IO_Error reports equal MySQL server ids"}, Recommendation: "Change the replica server-id, restart MySQL and then restart replication."})
	}
	if strings.Contains(joined, "fatal error 1236") || strings.Contains(joined, "Could not find first log file") {
		out = append(out, core.Finding{Service: "mysql", RuleID: "mysql.binlog_1236", Severity: core.SeverityCritical, Summary: "MySQL replica points to a missing binlog", Evidence: []string{"Last_IO_Error contains error 1236 or missing binlog"}, Recommendation: "After data-risk confirmation, align replication to the source current SHOW MASTER STATUS position."})
	}
	return out
}

func join(commands []core.CommandResult) string {
	var b strings.Builder
	for _, c := range commands {
		b.WriteString(c.Stdout)
		b.WriteString("\n")
		b.WriteString(c.Stderr)
		b.WriteString("\n")
		b.WriteString(c.Error)
		b.WriteString("\n")
	}
	return b.String()
}
