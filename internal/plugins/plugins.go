package plugins

import (
	"strings"
	"time"

	"mw-sre-platform/internal/core"
	"mw-sre-platform/internal/executor"
)

// RuntimeConfig is the small, stable input surface passed from CLI/platform
// control planes into service plugins. Keep secrets in Env-backed fields and
// keep plugin logic read-only unless a runbook explicitly moves to apply mode.
type RuntimeConfig struct {
	Profile          string
	Target           string
	Timeout          time.Duration
	RedisPass        string
	OpenGaussDataDir string
}

// Plugin is the code-side contract matching docs/PLUGIN_SPEC.md.
// In v0.2 the CLI wires discover/collect/diagnose. plan/apply/verify/rollback
// remain expressed as Runbook-as-Code files and are not auto-executed.
type Plugin interface {
	Name() string
	Aliases() []string
	Commands(RuntimeConfig) []executor.Command
	Findings([]core.CommandResult) []core.Finding
}

func Match(p Plugin, module string) bool {
	m := Normalize(module)
	if m == "all" || m == Normalize(p.Name()) {
		return true
	}
	for _, a := range p.Aliases() {
		if m == Normalize(a) {
			return true
		}
	}
	return false
}

func Normalize(s string) string {
	r := strings.NewReplacer("_", "-", " ", "-", ".", "-")
	return strings.ToLower(r.Replace(strings.TrimSpace(s)))
}

func FilterCommands(commands []core.CommandResult, service string) []core.CommandResult {
	service = Normalize(service)
	prefixDash := service + "-"
	prefixDot := service + "."
	out := make([]core.CommandResult, 0, len(commands))
	for _, c := range commands {
		trimmed := strings.ToLower(strings.TrimSpace(c.Name))
		n := Normalize(c.Name)
		if n == service || strings.HasPrefix(n, prefixDash) || strings.HasPrefix(trimmed, prefixDot) {
			out = append(out, c)
		}
	}
	return out
}

func JoinEvidence(commands []core.CommandResult) string {
	var b strings.Builder
	for _, c := range commands {
		b.WriteString(c.Name)
		b.WriteString("\n")
		b.WriteString(c.Stdout)
		b.WriteString("\n")
		b.WriteString(c.Stderr)
		b.WriteString("\n")
		b.WriteString(c.Error)
		b.WriteString("\n")
	}
	return b.String()
}
