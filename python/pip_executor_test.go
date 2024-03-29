package python_test

import (
	"context"
	"github.com/ATenderholt/terraform-provider-python-package/python"
	"testing"
)

func TestExecute(t *testing.T) {
	pip := python.NewPipExecutor("pip3", "")
	err := pip.Execute(context.TODO())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
