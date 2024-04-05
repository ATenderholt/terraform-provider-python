package python

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
