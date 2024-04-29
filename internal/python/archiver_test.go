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
		assert.Equal(t, hexToBase64("5af38573cd4dd654a61731be4a9d19826b3335577d5e9a0c96bdf39e7b65ae58"), checksum)
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
		assert.Equal(t, hexToBase64("914fb57b71e661695db1b35a89a158a8a51e4e723033c493179a1355e55efbca"), checksum)
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
		assert.Equal(t, hexToBase64("842611c6d40cc437abda689b68204416172152e5b70072d7a681e510ca08f40f"), checksum)
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
		assert.Equal(t, hexToBase64("11100ecdbc7f1b563e806c9c1bee7448d7b021b17509763bca284e03f75264f0"), checksum)
	})

	err = a.ArchiveDir("test-fixtures/example", "", []string{"*.py", "**/*.py"})
	assert.NoError(t, err)
}
