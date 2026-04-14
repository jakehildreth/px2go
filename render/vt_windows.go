//go:build windows

package render

import "golang.org/x/sys/windows"

func enableWindowsVT() error {
	h, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return err
	}
	const enableVirtualTerminalProcessing uint32 = 0x0004
	return windows.SetConsoleMode(h, mode|enableVirtualTerminalProcessing)
}
