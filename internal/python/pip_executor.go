package python

import (
	"context"
	"os/exec"
)

type PipExecutor struct {
	command      string
	requirements string
	installPath  string
	extraFlags   []string
}

func NewPipExecutor(command string, requirements string, installPath string, extraFlags []string) PipExecutor {
	return PipExecutor{
		command:      command,
		requirements: requirements,
		installPath:  installPath,
		extraFlags:   extraFlags,
	}
}

func (p PipExecutor) Execute(ctx context.Context) error {
	cmd := exec.CommandContext(ctx,
		p.command,
		"install", "-r", p.requirements,
		"-t", p.installPath)

	LogDebug(ctx, "Executing command", map[string]interface{}{
		"command": cmd.String(),
	})

	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError(ctx, "error when running command", map[string]interface{}{
			"command": cmd.String(),
			"error":   err,
		})
		return err
	}

	LogDebug(ctx, "output from running command", map[string]interface{}{
		"command": cmd.String(),
		"output":  string(output),
	})

	return nil
}
