package opengauss

import (
	"testing"

	"mw-sre-platform/internal/core"
)

func TestOpenGaussHANotNormal(t *testing.T) {
	out := `HA state:
        local_role                     : Standby
        db_state                       : Need repair
        detail_information             : Abnormal`
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "opengauss.gsctl-query", Stdout: out}})
	assertRule(t, findings, "opengauss.ha_not_normal")
}

func TestOpenGaussStandbyReceiverMissing(t *testing.T) {
	out := "HA state:\n        local_role                     : Standby\n        db_state                       : Normal\n        detail_information             : Normal\n\n Receiver info:\nNo information\n"
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "opengauss.gsctl-query", Stdout: out}})
	assertRule(t, findings, "opengauss.standby_receiver_missing")
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
