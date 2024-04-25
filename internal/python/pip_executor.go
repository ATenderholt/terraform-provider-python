package python

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	err := os.RemoveAll(p.installPath)
	if err != nil {
		LogError(ctx, "unable to cleanup directory", map[string]interface{}{
			"installPath": p.installPath,
		})
	}

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
			"output":  string(output),
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

func (p PipExecutor) GetPythonVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, p.command, "--version")

	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		LogError(ctx, "error when running command", map[string]interface{}{
			"command": cmd.String(),
			"output":  output,
			"error":   err,
		})
		return "", err
	}

	LogDebug(ctx, "output from running command", map[string]interface{}{
		"command": cmd.String(),
		"output":  output,
	})

	re := regexp.MustCompile("\\(python (.*)\\)")
	matches := re.FindStringSubmatch(output)
	if len(matches) != 2 {
		return "", fmt.Errorf("unable to find python version from pip: %s", output)
	}

	return matches[1], nil
}
