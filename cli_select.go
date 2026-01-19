package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func Select(label string, items []string) (int, string, error) {
	if len(items) == 0 {
		return 0, "", errors.New("no items to select")
	}

	var result int
	var filter string

	err := withRawMode(func() error {
		selectedIdx := 0
		filteredItems := items
		filteredIndices := make([]int, len(items))
		for i := range items {
			filteredIndices[i] = i
		}
		lastLineCount := 0

		// Hide cursor during selection
		fmt.Print(escHideCursor)
		defer fmt.Print(escShowCursor)

		clearPrevious := func() {
			// Move up and clear each line
			for i := 0; i < lastLineCount; i++ {
				fmt.Print("\033[A")   // Move up
				fmt.Print("\r\033[K") // Clear line
			}
		}

		redraw := func() {
			// Filter items first
			if filter != "" {
				filteredItems = nil
				filteredIndices = nil
				for i, item := range items {
					if strings.Contains(strings.ToLower(item), strings.ToLower(filter)) {
						filteredItems = append(filteredItems, item)
						filteredIndices = append(filteredIndices, i)
					}
				}
				if selectedIdx >= len(filteredItems) {
					selectedIdx = len(filteredItems) - 1
				}
				if selectedIdx < 0 {
					selectedIdx = 0
				}
			} else {
				filteredItems = items
				filteredIndices = make([]int, len(items))
				for i := range items {
					filteredIndices[i] = i
				}
			}

			lineCount := 0

			// Print label
			if label != "" {
				fmt.Print("\r\033[K" + themeAccent(bold(label)) + "\r\n")
				lineCount++
			}

			// Print filter line if active
			if filter != "" {
				fmt.Print("\r\033[K" + themeMuted("Filter: ") + themeText(filter) + "\r\n")
				lineCount++
			}

			// Print options
			for i, item := range filteredItems {
				fmt.Print("\r\033[K")
				if i == selectedIdx {
					fmt.Print(themeSuccess("> ") + themeSuccess(item))
				} else {
					fmt.Print(themeSubtle("  ") + themeText(item))
				}
				fmt.Print("\r\n")
				lineCount++
			}

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
				if len(filteredItems) > 0 {
					result = filteredIndices[selectedIdx]
					clearPrevious()
					if label != "" {
						fmt.Print("\r\033[K" + themeAccent(bold(label)) + "\r\n")
					}
					fmt.Print("\r\033[K" + themeSuccess("> ") + themeText(items[result]) + "\r\n")
					return nil
				}

			case keyUp:
				if selectedIdx > 0 {
					selectedIdx--
				}

			case keyDown:
				if selectedIdx < len(filteredItems)-1 {
					selectedIdx++
				}

			case keyBackspace:
				if len(filter) > 0 {
					filter = filter[:len(filter)-1]
				}

			case keyEscape:
				filter = ""
				selectedIdx = 0

			default:
				if char != 0 && char >= 32 {
					filter += string(char)
				}
			}

			clearPrevious()
			redraw()
		}
	})

	if err != nil {
		return 0, "", err
	}

	return result, items[result], nil
}

func MustSelect(label string, items []string) (int, string) {
	index, value, err := Select(label, items)

	if err != nil {
		Fatal(err)
	}

	return index, value
}
