// +build windows
package modern

func DefaultDirectory() string {
	return filepath.Join(os.Getenv("localappdata"), walletApp)
}
