package tools_test

import (
	"context"
	"github.com/ATenderholt/terraform-provider-python-package/python/tools"
	"testing"
)

func TestExecute(t *testing.T) {
	pip := tools.NewPipExecutor("pip3",
		"./testdata/requirements.txt", "./testdata/output/opt/python",
		[]string{})
	err := pip.Execute(context.TODO())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
