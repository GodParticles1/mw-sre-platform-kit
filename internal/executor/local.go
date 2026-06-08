package executor

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"

	"mw-sre-platform/internal/core"
)

type LocalExecutor struct{}

func NewLocalExecutor() *LocalExecutor { return &LocalExecutor{} }

func (e *LocalExecutor) Run(ctx context.Context, cmd Command) core.CommandResult {
	start := time.Now()
	res := core.CommandResult{
		Name:      cmd.Name,
		Host:      cmd.Host,
		Command:   cmd.Args,
		StartedAt: start,
		ExitCode:  0,
	}
	if len(cmd.Args) == 0 {
		res.ExitCode = 127
		res.Error = "empty command"
		res.FinishedAt = time.Now()
		res.Duration = res.FinishedAt.Sub(start)
		return res
	}

	c := exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	res.Stdout = stdout.String()
	res.Stderr = stderr.String()
	if err != nil {
		res.Error = err.Error()
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			res.ExitCode = ee.ExitCode()
		} else {
			res.ExitCode = 1
		}
	}
	res.FinishedAt = time.Now()
	res.Duration = res.FinishedAt.Sub(start)
	return res
}
