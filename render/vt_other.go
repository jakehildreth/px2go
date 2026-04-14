//go:build !windows

package render

func enableWindowsVT() error {
	return nil
}
