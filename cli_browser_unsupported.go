//go:build !linux && !windows && !darwin

package cli

import (
	"fmt"
	"runtime"
)

func openBrowser(url string) error {
	return fmt.Errorf("openBrowser: unsupported operating system: %v", runtime.GOOS)
}
