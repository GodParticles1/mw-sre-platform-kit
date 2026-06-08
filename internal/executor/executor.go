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
}

type Executor interface {
	Run(ctx context.Context, cmd Command) core.CommandResult
}
