package python

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"os/exec"
	"strings"
)

type PipExecutor struct {
	command    string
	extraFlags string
}

func NewPipExecutor(command string, extraFlags string) PipExecutor {
	return PipExecutor{
		command:    command,
		extraFlags: extraFlags,
	}
}

func (p PipExecutor) Execute(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, p.command, "install", "--help")
	fmt.Printf("executing command: %s\n", cmd.String())
	tflog.Debug(ctx, "Executing command", map[string]interface{}{
		"command": cmd.String(),
	})

	var out strings.Builder
	cmd.Stdout = &out

	//stdoutPipe, err := cmd.StdoutPipe()
	//if err != nil {
	//	tflog.Error(ctx, "unable to get command stdout", map[string]interface{}{
	//		"error": err,
	//	})
	//	return err
	//}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		tflog.Error(ctx, "unable to get command stderr", map[string]interface{}{
			"error": err,
		})
		return err
	}

	//defer logOutput(ctx, stdoutPipe, stdErrPipe)

	err = cmd.Run()
	if err != nil {
		tflog.Error(ctx, "error when running command", map[string]interface{}{
			"command": cmd.String(),
			"error":   err,
		})
		return err
	}
	cmd.Wait()

	//stdout, _ := io.ReadAll(stdoutPipe)
	stderr, _ := io.ReadAll(stderrPipe)

	fmt.Printf("error: %s\n", stderr)
	fmt.Printf("output: %s\n", out.String())

	return nil
}

func logOutput(ctx context.Context, stdoutPipe io.ReadCloser, stderrPipe io.ReadCloser) func() {
	return func() {
		defer stdoutPipe.Close()
		defer stderrPipe.Close()

		stdout, _ := io.ReadAll(stdoutPipe)
		stderr, _ := io.ReadAll(stderrPipe)

		fmt.Printf("error: %s\n", stderr)
		fmt.Printf("output: %s\n", stdout)
		tflog.Debug(ctx, "command output", map[string]interface{}{
			"stdout": stdout,
			"stderr": stderr,
		})
	}
}
