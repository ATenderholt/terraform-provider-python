package tools_test

import (
	"context"
	"github.com/ATenderholt/terraform-provider-python-package/python/tools"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateDirChecksum(t *testing.T) {
	checksum, err := tools.CalculateDirChecksum(context.TODO(), "./testdata/checksums")

	assert.NoError(t, err)
	assert.Equal(t, "74cdc5a52e3f87dcc817127a8b25dc06", checksum)
}
