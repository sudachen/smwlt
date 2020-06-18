package fu

import (
	"github.com/sudachen/smwlt/fu/errstr"
	"io"
	"os"
	"path/filepath"
)

func SaveWithBackup(path string, writer func(io.Writer) error) (err error) {
	if _, e := os.Stat(path); e == nil {
		_ = os.Remove(path + "~")
		if err = os.Rename(path, path+"~"); err != nil {
			return errstr.Wrapf(1, err, "failed to backup file: %v", err.Error())
		}
	}

	defer func() {
		if err != nil {
			if _, e := os.Stat(path); e != os.ErrNotExist {
				_ = os.Rename(path+"~", path)
			}
		}
	}()

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	if err = writer(f); err != nil {
		return
	}
	if err = f.Close(); err != nil {
		return
	}
	_ = os.Remove(path + "~")
	return

}
