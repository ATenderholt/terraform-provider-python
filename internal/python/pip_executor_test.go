package python_test

import (
	"context"
	"github.com/ATenderholt/terraform-provider-python/internal/python"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestPipExecutor_Install(t *testing.T) {
	td := os.TempDir()
	output := filepath.Join(td, "pip_output")
	t.Logf("Using pip to install to %s", output)

	pip := python.NewPipExecutor("pip3")
	err := pip.Install(context.TODO(), "./test-fixtures/example/requirements.txt", output)

	assert.NoError(t, err, "unexpected error when running pip")
	assert.FileExists(t, filepath.Join(output, "requests", "__init__.py"))
	assert.FileExists(t, filepath.Join(output, "urllib3", "__init__.py"))
}

func TestPipExecutor_GetPythonVersion(t *testing.T) {
	pip := python.NewPipExecutor("pip3.10")
	version, err := pip.GetPythonVersion(context.TODO())

	assert.NoError(t, err, "unexpected error when getting python version")
	assert.Equal(t, "3.10", version)
}
