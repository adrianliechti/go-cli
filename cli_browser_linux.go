//go:build linux

package cli

import (
	"os/exec"
	"strings"
)

func openBrowser(url string) error {
	providers := []string{"xdg-open", "x-www-browser", "www-browser", "wslview"}

	// There are multiple possible providers to open a browser on linux
	// One of them is xdg-open, another is x-www-browser, then there's www-browser, etc.
	// wslview is used for Windows Subsystem for Linux (WSL)
	// Look for one that exists and run it
	for _, provider := range providers {
		if _, err := exec.LookPath(provider); err == nil {
			return exec.Command(provider, url).Run()
		}
	}

	return &exec.Error{Name: strings.Join(providers, ","), Err: exec.ErrNotFound}
}
