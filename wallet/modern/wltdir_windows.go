// +build windows
package modern

import (
	"os"
	"path/filepath"
)

func DefaultDirectory() string {
	return filepath.Join(os.Getenv("localappdata"), walletApp)
}
