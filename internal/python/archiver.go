package python

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Archiver struct {
	outputPath string
	file       *os.File
	writer     *zip.Writer
}

func NewArchiver(outputPath string) Archiver {
	return Archiver{
		outputPath: outputPath,
	}
}

func (a *Archiver) Open() error {
	dir := filepath.Dir(a.outputPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("unable to create dir '%s': %v", dir, err)
	}

	file, err := os.Create(a.outputPath)
	if err != nil {
		return fmt.Errorf("unable to create archive file: %v", err)
	}
	a.file = file
	a.writer = zip.NewWriter(a.file)

	return nil
}

func (a *Archiver) Close() {
	if a.writer != nil {
		a.writer.Close()
		a.writer = nil
	}

	if a.file != nil {
		a.file.Close()
		a.file = nil
	}
}

func (a *Archiver) ArchiveFile(in string, out string) error {
	fi, err := os.Stat(in)
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("unable to archive missing file=%s", in)
	case err != nil:
		return fmt.Errorf("unable to archive file=%s: %v", in, err)
	}

	file, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("unable to open file=%s: %v", in, err)
	}
	defer file.Close()

	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return fmt.Errorf("unable to generate zip header for file=%s: %v", in, err)
	}
	fh.Name = filepath.ToSlash(out)
	fh.Method = zip.Deflate
	// Need to provide a fixed date so the zip checksum doesn't change
	fh.Modified = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	f, err := a.writer.CreateHeader(fh)
	if err != nil {
		return fmt.Errorf("unable to create zip header for file=%s: %v", in, err)
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return fmt.Errorf("unable to archive content for file=%s: %v", in, err)
	}

	return nil
}

func (a *Archiver) ArchiveDir(dir string, root string, excludeGlobs []string) error {
	excludes := make(map[string]interface{})
	for _, excludeGlob := range excludeGlobs {
		matches, err := filepath.Glob(filepath.Join(dir, excludeGlob))
		if err != nil {
			return fmt.Errorf("invalid exclude format=%s: %v", excludeGlob, err)
		}

		for _, match := range matches {
			excludes[match] = true
		}
	}

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error encountered walking dir=%s: %v", dir, err)
		}

		// If current path is a directory and matches the exclude glob, we should NOT continue walking
		// If current path is a directory and doesn't match the exclude glob, we should continue walking
		// If current path is NOT a directory and matches the exclude glob, we should continue walking
		// otherwise, add the current path to the archive.
		_, exclude := excludes[path]
		switch {
		case d.IsDir() && exclude:
			return filepath.SkipDir
		case d.IsDir() && !exclude:
			return nil
		case !d.IsDir() && exclude:
			return nil
		}

		base, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("unable to relativize path=%s: %v", path, err)
		}

		joined := filepath.Join(root, base)
		return a.ArchiveFile(path, filepath.ToSlash(joined))
	})
}
