package core

import "time"

// Severity is intentionally small and stable so findings can be consumed by
// CLI users, AIOps systems and Kubernetes CRDs without translation glue.
type Severity string

const (
	SeverityOK       Severity = "ok"
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
	SeverityUnknown  Severity = "unknown"
)

func (s Severity) Valid() bool {
	switch s {
	case SeverityOK, SeverityInfo, SeverityWarning, SeverityCritical, SeverityUnknown:
		return true
	default:
		return false
	}
}

// CommandResult is the raw evidence unit. Keep it close to command execution;
// do not mix diagnosis into it.
type CommandResult struct {
	Name       string        `json:"name"`
	Host       string        `json:"host"`
	Command    []string      `json:"command"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	ExitCode   int           `json:"exit_code"`
	Duration   time.Duration `json:"duration_ns"`
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Error      string        `json:"error,omitempty"`
}

// Finding is a normalized diagnostic statement derived from evidence.
type Finding struct {
	Service       string   `json:"service"`
	RuleID        string   `json:"rule_id"`
	Severity      Severity `json:"severity"`
	Summary       string   `json:"summary"`
	Evidence      []string `json:"evidence"`
	Recommendation string  `json:"recommendation,omitempty"`
}

// CheckReport is the stable machine-readable output consumed by platform/AIOps.
type CheckReport struct {
	RunID     string          `json:"run_id"`
	Profile   string          `json:"profile"`
	Target    string          `json:"target"`
	Module    string          `json:"module"`
	StartedAt time.Time       `json:"started_at"`
	EndedAt   time.Time       `json:"ended_at"`
	Commands  []CommandResult `json:"commands"`
	Findings  []Finding       `json:"findings"`
}
