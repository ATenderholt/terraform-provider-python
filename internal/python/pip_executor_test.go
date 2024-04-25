package python_test

import (
	"context"
	"github.com/ATenderholt/terraform-provider-python/internal/python"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestPipExecutor_Execute(t *testing.T) {
	td := os.TempDir()
	output := filepath.Join(td, "pip_output")
	t.Logf("Using pip to install to %s", output)

	pip := python.NewPipExecutor("pip3",
		"./test-fixtures/example/requirements.txt", output,
		[]string{})
	err := pip.Execute(context.TODO())

	assert.NoError(t, err, "unexpected error when running pip")
	assert.FileExists(t, filepath.Join(output, "dataclasses_avroschema", "__init__.py"))
	assert.FileExists(t, filepath.Join(output, "fastavro", "__init__.py"))
}

func TestPipExecutor_GetPythonVersion(t *testing.T) {
	pip := python.NewPipExecutor("pip3.10",
		"./test-fixtures/example/requirements.txt", "",
		[]string{})

	version, err := pip.GetPythonVersion(context.TODO())

	assert.NoError(t, err, "unexpected error when getting python version")
	assert.Equal(t, "3.10", version)
}
