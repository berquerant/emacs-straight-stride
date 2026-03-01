package temp

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type Dir struct {
	path string
}

func NewDir(prefix string) (*Dir, error) {
	path := filepath.Join(os.TempDir(), prefix)
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("%w: failed to create tempdir, prefix = %s", err, prefix)
	}
	return &Dir{
		path: path,
	}, nil
}

func (d *Dir) Close() error {
	slog.Debug("TempDir: Close", slog.String("path", d.path))
	if err := os.RemoveAll(d.path); err != nil {
		return fmt.Errorf("%w: failed to close tempdir, path = %s", err, d.path)
	}
	return nil
}

func (d *Dir) Join(path string) string {
	x := filepath.Join(d.path, path)
	slog.Debug("TempDir: Join", slog.String("path", x))
	return x
}
