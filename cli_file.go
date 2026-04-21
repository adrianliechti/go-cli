package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unicode/utf8"
)

// fileEntry represents a file or directory in the browser
type fileEntry struct {
	name  string
	path  string
	isDir bool
}

// isRootDir checks if a path is a root directory (cross-platform)
func isRootDir(path string) bool {
	if runtime.GOOS == "windows" {
		// On Windows, root looks like "C:\" or "C:/"
		if len(path) == 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/') {
			return true
		}
		// Also handle "C:" without trailing slash
		if len(path) == 2 && path[1] == ':' {
			return true
		}
		return false
	}
	return path == "/"
}

func File(label string, types []string) (string, error) {
	var result string

	err := withRawMode(func() error {
		// Start in current directory
		currentDir, err := os.Getwd()
		if err != nil {
			currentDir = "."
		}

		selectedIdx := 0
		scrollOffset := 0
		maxVisible := 12
		filter := ""
		lastLineCount := 0

		// Hide cursor during selection
		fmt.Print(escHideCursor)
		defer fmt.Print(escShowCursor)

		// Helper to navigate to a new directory
		navigateToDir := func(dir string) {
			currentDir = dir
			filter = ""
			selectedIdx = 0
			scrollOffset = 0
		}

		// Read directory contents
		readDir := func(dir string) []fileEntry {
			entries := []fileEntry{}

			// Add parent directory option if not at root
			if !isRootDir(dir) {
				entries = append(entries, fileEntry{
					name:  "..",
					path:  filepath.Dir(dir),
					isDir: true,
				})
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				return entries
			}

			// Separate directories and files for sorting
			var dirs, regularFiles []fileEntry
			for _, f := range files {
				// Skip hidden files
				if strings.HasPrefix(f.Name(), ".") {
					continue
				}

				fullPath := filepath.Join(dir, f.Name())
				entry := fileEntry{
					name:  f.Name(),
					path:  fullPath,
					isDir: f.IsDir(),
				}

				if f.IsDir() {
					dirs = append(dirs, entry)
				} else {
					// Filter by extension if types specified
					if len(types) > 0 {
						ext := strings.ToLower(filepath.Ext(f.Name()))
						match := false
						for _, t := range types {
							if ext == strings.ToLower(t) {
								match = true
								break
							}
						}
						if !match {
							continue
						}
					}
					regularFiles = append(regularFiles, entry)
				}
			}

			// Sort directories and files alphabetically
			sort.Slice(dirs, func(i, j int) bool {
				return strings.ToLower(dirs[i].name) < strings.ToLower(dirs[j].name)
			})
			sort.Slice(regularFiles, func(i, j int) bool {
				return strings.ToLower(regularFiles[i].name) < strings.ToLower(regularFiles[j].name)
			})

			entries = append(entries, dirs...)
			entries = append(entries, regularFiles...)
			return entries
		}

		allEntries := readDir(currentDir)
		filteredEntries := allEntries

		// Filter entries by search term
		filterEntries := func() {
			if filter == "" {
				filteredEntries = allEntries
			} else {
				filteredEntries = []fileEntry{}
				lowerFilter := strings.ToLower(filter)
				for _, e := range allEntries {
					if strings.Contains(strings.ToLower(e.name), lowerFilter) {
						filteredEntries = append(filteredEntries, e)
					}
				}
			}
			if selectedIdx >= len(filteredEntries) {
				selectedIdx = len(filteredEntries) - 1
			}
			if selectedIdx < 0 {
				selectedIdx = 0
			}
		}

		clearPrevious := func() {
			for i := 0; i < lastLineCount; i++ {
				fmt.Print("\033[A")
				fmt.Print("\r\033[K")
			}
		}

		redraw := func() {
			lineCount := 0

			// Print label
			prompt := themeAccent(bold(label))
			if len(types) > 0 {
				prompt += " " + themeSubtle("("+strings.Join(types, ", ")+")")
			}
			fmt.Print("\r\033[K" + prompt + "\r\n")
			lineCount++

			// Print current path
			displayPath := currentDir
			home, _ := os.UserHomeDir()
			if home != "" && strings.HasPrefix(displayPath, home) {
				displayPath = "~" + displayPath[len(home):]
			}
			fmt.Print("\r\033[K" + themeMuted("▸ ") + themeText(displayPath) + "\r\n")
			lineCount++

			// Print filter line if active
			if filter != "" {
				fmt.Print("\r\033[K" + themeMuted("/ ") + themeText(filter) + "\r\n")
				lineCount++
			}

			// Handle empty directory
			if len(filteredEntries) == 0 {
				fmt.Print("\r\033[K" + themeMuted("  (empty)") + "\r\n")
				lineCount++
			} else {
				// Adjust scroll offset
				if selectedIdx < scrollOffset {
					scrollOffset = selectedIdx
				}
				if selectedIdx >= scrollOffset+maxVisible {
					scrollOffset = selectedIdx - maxVisible + 1
				}

				// Print entries
				visibleEnd := scrollOffset + maxVisible
				if visibleEnd > len(filteredEntries) {
					visibleEnd = len(filteredEntries)
				}

				// Show scroll indicator at top
				if scrollOffset > 0 {
					fmt.Print("\r\033[K" + themeMuted("  ↑ more items above") + "\r\n")
					lineCount++
				}

				for i := scrollOffset; i < visibleEnd; i++ {
					entry := filteredEntries[i]
					fmt.Print("\r\033[K")

					prefix := "  "
					if i == selectedIdx {
						prefix = themeSuccess("> ")
					}

					icon := "  "
					if entry.isDir {
						icon = "▸ "
					}

					name := entry.name
					if entry.isDir && entry.name != ".." {
						name += "/"
					}

					if i == selectedIdx {
						fmt.Print(prefix + icon + themeSuccess(name))
					} else {
						if entry.isDir {
							fmt.Print(prefix + icon + themeAccent(name))
						} else {
							fmt.Print(prefix + icon + themeText(name))
						}
					}
					fmt.Print("\r\n")
					lineCount++
				}

				// Show scroll indicator at bottom
				if visibleEnd < len(filteredEntries) {
					fmt.Print("\r\033[K" + themeMuted("  ↓ more items below") + "\r\n")
					lineCount++
				}
			}

			// Print help
			fmt.Print("\r\033[K" + themeMuted("↑/↓ navigate • Enter select • ← parent • → enter dir • Type to filter • Esc clear") + "\r\n")
			lineCount++

			lastLineCount = lineCount
		}

		// Initial draw
		redraw()

		for {
			key, char, err := readKey(os.Stdin)
			if err != nil {
				return err
			}

			switch key {
			case keyCtrlC:
				clearPrevious()
				return ErrUserAborted

			case keyEnter:
				if len(filteredEntries) > 0 {
					entry := filteredEntries[selectedIdx]
					if entry.isDir {
						// Navigate into directory
						navigateToDir(entry.path)
						allEntries = readDir(currentDir)
						filteredEntries = allEntries
					} else {
						// Select file
						result = entry.path
						clearPrevious()
						// Show final selection
						prompt := themeAccent(bold(label))
						if len(types) > 0 {
							prompt += " " + themeSubtle("("+strings.Join(types, ", ")+")")
						}
						fmt.Print("\r\033[K" + prompt + "\r\n")
						fmt.Print("\r\033[K" + themeSuccess("> ") + themeText(result) + "\r\n")
						return nil
					}
				}

			case keyUp:
				if selectedIdx > 0 {
					selectedIdx--
				}

			case keyDown:
				if len(filteredEntries) > 0 && selectedIdx < len(filteredEntries)-1 {
					selectedIdx++
				}

			case keyLeft:
				// Go to parent directory
				if !isRootDir(currentDir) {
					navigateToDir(filepath.Dir(currentDir))
					allEntries = readDir(currentDir)
					filteredEntries = allEntries
				}

			case keyRight:
				// Enter directory if selected
				if len(filteredEntries) > 0 {
					entry := filteredEntries[selectedIdx]
					if entry.isDir {
						navigateToDir(entry.path)
						allEntries = readDir(currentDir)
						filteredEntries = allEntries
					}
				}

			case keyBackspace:
				if len(filter) > 0 {
					_, size := utf8.DecodeLastRuneInString(filter)
					filter = filter[:len(filter)-size]
					filterEntries()
				}

			case keyEscape:
				filter = ""
				filterEntries()

			case keyHome:
				// Go to home directory
				home, err := os.UserHomeDir()
				if err == nil {
					navigateToDir(home)
					allEntries = readDir(currentDir)
					filteredEntries = allEntries
				}

			default:
				if char != 0 && char >= 32 {
					filter += string(char)
					filterEntries()
				}
			}

			clearPrevious()
			redraw()
		}
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

func MustFile(label string, types []string) string {
	value, err := File(label, types)

	if err != nil {
		Fatal(err)
	}

	return value
}
