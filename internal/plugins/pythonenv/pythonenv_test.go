package pythonenv

import (
	"testing"

	"mw-sre-platform/internal/core"
)

func TestProfileSourceNonZero(t *testing.T) {
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "python-env.profile-source", Stdout: "PROFILE_SOURCE_EXIT=1\n"}})
	assertRule(t, findings, "python_env.profile_source_nonzero")
}

func TestProfileSourceOutput(t *testing.T) {
	out := "PROFILE_SOURCE_EXIT=0\nPROFILE_SOURCE_OUTPUT_BEGIN\nhello\nPROFILE_SOURCE_OUTPUT_END\n"
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "python-env.profile-source", Stdout: out}})
	assertRule(t, findings, "python_env.profile_source_output")
}

func TestImportFail(t *testing.T) {
	findings := (Plugin{}).Findings([]core.CommandResult{{Name: "python-env.import-smoke", Stdout: "IMPORT_FAIL=ssl:ImportError:libssl missing\n"}})
	assertRule(t, findings, "python_env.import_fail")
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
