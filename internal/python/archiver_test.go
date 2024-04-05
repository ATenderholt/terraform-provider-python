package python_test

import (
	"github.com/ATenderholt/terraform-provider-python-package/internal/python"
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
		assert.Equal(t, "c019027725aeb616daa7d7e588e512ad", checksum)
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
		assert.Equal(t, "9ca3aaba9f972d745aa835ba5326a0bc", checksum)
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
		assert.Equal(t, "169e013551fb6da1e7f3e9dbcdd8449c", checksum)
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
		assert.Equal(t, "47edba4bed71177ee35ac1d5c2d2dcb5", checksum)
	})

	err = a.ArchiveDir("test-fixtures/example", "", []string{"*.py", "**/*.py"})
	assert.NoError(t, err)
}
