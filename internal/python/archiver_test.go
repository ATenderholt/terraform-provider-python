package python_test

import (
	"github.com/ATenderholt/terraform-provider-python/internal/python"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiver_ArchiveFile_NoError(t *testing.T) {
	d := os.TempDir()
	path := filepath.Join(d, "example_archive.zip")
	t.Logf("Creating archive at %s", path)

	a := python.NewArchiver(path)
	err := a.Open()
	require.NoError(t, err)

	t.Cleanup(func() {
		a.Close()
		checksum, err := python.Checksum(path)
		assert.NoError(t, err)
		assert.Equal(t, "sP8C+nqT8APB7XoK49rWyYezKA+VKUVU4rGum0WxcQY=", checksum)
	})

	err = a.ArchiveFile("test-fixtures/example/main.py", "main.py")
	assert.NoError(t, err)

	err = a.ArchiveFile("test-fixtures/example/example/__init__.py", "example/__init__.py")
	assert.NoError(t, err)

	err = a.ArchiveFile("test-fixtures/example/example/message.py", "example/message.py")
	assert.NoError(t, err)
}

func TestArchiver_ArchiveFile_FileNotFoundError(t *testing.T) {
	d := os.TempDir()
	path := filepath.Join(d, "example_archive_missing_file.zip")
	t.Logf("Creating archive at %s", path)

	a := python.NewArchiver(path)
	err := a.Open()
	require.NoError(t, err)
	defer a.Close()

	err = a.ArchiveFile("example/main.py", "main.py")
	assert.ErrorContains(t, err, "unable to archive missing file=example/main.py")
}

func TestArchiver_ArchiveDir_NoError(t *testing.T) {
	d := os.TempDir()
	path := filepath.Join(d, "example_archive_dir.zip")
	t.Logf("Creating archive at %s", path)

	a := python.NewArchiver(path)
	err := a.Open()
	require.NoError(t, err)

	t.Cleanup(func() {
		a.Close()
		checksum, err := python.Checksum(path)
		assert.NoError(t, err)
		assert.Equal(t, "mXDSpyuc5PJwKzQ3NfwfNfWIpIoo7l7kZa6JTZmXR+k=", checksum)
	})

	err = a.ArchiveDir("test-fixtures/example", "/opt/python", []string{"requirements.txt"})
	assert.NoError(t, err)
}

func TestArchiver_ArchiveDir_WithoutRoot_NoError(t *testing.T) {
	d := os.TempDir()
	path := filepath.Join(d, "example_archive_dir_without_root.zip")
	t.Logf("Creating archive at %s", path)

	a := python.NewArchiver(path)
	err := a.Open()
	require.NoError(t, err)

	t.Cleanup(func() {
		a.Close()
		checksum, err := python.Checksum(path)
		assert.NoError(t, err)
		assert.Equal(t, "G99BHkVKpQzS5ica7MToKLXlch8IDdRj7we2pVbMGwg=", checksum)
	})

	err = a.ArchiveDir("test-fixtures/example", "", nil)
	assert.NoError(t, err)
}

func TestArchiver_ArchiveDir_WithoutRootExcludePy_NoError(t *testing.T) {
	d := os.TempDir()
	path := filepath.Join(d, "example_archive_dir_without_root_exclude_py.zip")
	t.Logf("Creating archive at %s", path)

	a := python.NewArchiver(path)
	err := a.Open()
	require.NoError(t, err)

	t.Cleanup(func() {
		a.Close()
		checksum, err := python.Checksum(path)
		assert.NoError(t, err)
		assert.Equal(t, "6z7sd7Qi85I+wYtrtuKB++0k3boV4I2NNn3VD0QgN8o=", checksum)
	})

	err = a.ArchiveDir("test-fixtures/example", "", []string{"*.py", "**/*.py"})
	assert.NoError(t, err)
}
