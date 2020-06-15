// +build linux
package modern

import (
	"os/user"
	"path/filepath"
)

func DefaultDirectory() string {
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	return filepath.Join(usr.HomeDir, ".config", walletApp)
}
