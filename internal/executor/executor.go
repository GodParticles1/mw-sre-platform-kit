package executor

import (
	"context"
	"mw-sre-platform/internal/core"
)

type Command struct {
	Name    string
	Host    string
	Args    []string
	Timeout int
	Env     map[string]string
}

type Executor interface {
	Run(ctx context.Context, cmd Command) core.CommandResult
}
