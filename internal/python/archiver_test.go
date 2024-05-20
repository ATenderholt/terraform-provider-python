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
		assert.Equal(t, hexToBase64("90bae60c82474b7ee708e3be0e62f7d606e5a51b2e34858ee316b9bc56b401a9"), checksum)
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
		assert.Equal(t, hexToBase64("6811b35f81094553d5b81ed2adb3ad9416a503368d15aeee5af51f55805fe600"), checksum)
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
		assert.Equal(t, hexToBase64("cf0edaa1c8777629bfd7b9b47d279720039ec5ffb5c510704c4276e49bb5003d"), checksum)
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
		assert.Equal(t, hexToBase64("9444b186108fb0c2efee4d5914fedcc5c47a03334c6d94ba46226744857b5a40"), checksum)
	})

	err = a.ArchiveDir("test-fixtures/example", "", []string{"*.py", "**/*.py"})
	assert.NoError(t, err)
}
