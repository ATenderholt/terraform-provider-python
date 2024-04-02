package python

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CalculateDirChecksum(ctx context.Context, dir string) (string, error) {
	var buf bytes.Buffer
	h := md5.New()
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		LogDebug(ctx, "processing path", map[string]interface{}{
			"path":  path,
			"isDir": d.IsDir(),
			"error": err,
		})

		if err != nil {
			LogError(ctx, "error walking directory", map[string]interface{}{
				"dir":   dir,
				"path":  path,
				"error": err,
			})
			return err
		}

		if d.IsDir() {
			return nil
		}

		buf.WriteString(path)
		f, err := os.Open(path)
		if err != nil {
			LogError(ctx, "unable to open file for reading", map[string]interface{}{
				"path":  path,
				"error": err,
			})
			return err
		}
		defer f.Close()

		_, err = io.Copy(h, f)
		if err != nil {
			LogError(ctx, "unable to read file", map[string]interface{}{
				"path":  path,
				"error": err,
			})
			return err
		}

		hash := h.Sum(nil)
		buf.WriteString(hex.EncodeToString(hash))
		h.Reset()

		return nil
	})

	if err != nil {
		LogError(ctx, "unable to walk directory", map[string]interface{}{
			"dir":   dir,
			"error": err,
		})
		return "", err
	}

	LogDebug(ctx, "result", map[string]interface{}{
		"buffer": buf.String(),
	})

	// add a LF character so can compare with md5 of file in terminal
	buf.WriteByte(10)

	h.Reset()
	h.Write(buf.Bytes())

	hash := h.Sum(nil)
	return hex.EncodeToString(hash), err
}
