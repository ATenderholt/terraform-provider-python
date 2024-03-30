package python

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	tflog.Debug(ctx, "Executing command", map[string]interface{}{
		"command": cmd.String(),
	})

	output, err := cmd.CombinedOutput()
	if err != nil {
		tflog.Error(ctx, "error when running command", map[string]interface{}{
			"command": cmd.String(),
			"error":   err,
		})
		return err
	}

	tflog.Debug(ctx, "output from running command", map[string]interface{}{
		"command": cmd.String(),
		"output":  output,
	})

	fmt.Printf("output: %s\n", output)
	return nil
}
