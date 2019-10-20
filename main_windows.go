// +build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"

	"github.com/luhring/reach/cmd"
)

func main() {
	var originalMode uint32
	stdout := windows.Handle(os.Stdout.Fd())

	_ = windows.GetConsoleMode(stdout, &originalMode)
	_ = windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	defer windows.SetConsoleMode(stdout, originalMode)

	cmd.Execute()
}
