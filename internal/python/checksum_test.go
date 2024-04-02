package python_test

import (
	"context"
	"github.com/ATenderholt/terraform-provider-python-package/internal/python"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateDirChecksum(t *testing.T) {
	checksum, err := python.CalculateDirChecksum(context.TODO(), "./test-fixtures/content")

	assert.NoError(t, err)
	assert.Equal(t, "d466b73379759351df08f2668cfb6066", checksum)
}
