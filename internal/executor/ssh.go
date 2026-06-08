package executor

import (
	"context"
	"mw-sre-platform/internal/core"
)

// SSHExecutor intentionally shells out to ssh in v0.1. The interface keeps the
// door open for a native SSH implementation later.
type SSHExecutor struct {
	local *LocalExecutor
}

func NewSSHExecutor() *SSHExecutor {
	return &SSHExecutor{local: NewLocalExecutor()}
}

func (e *SSHExecutor) Run(ctx context.Context, cmd Command) core.CommandResult {
	args := []string{"ssh", cmd.Host}
	args = append(args, cmd.Args...)
	wrapped := cmd
	wrapped.Args = args
	return e.local.Run(ctx, wrapped)
}
