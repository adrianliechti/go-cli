package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Select lets the user pick one item, filtering as they type. An item may
// carry a description after a tab ("Wingman\tTCP 4242"); it is shown muted
// and aligned, and left out of the returned value and the summary line.
func Select(label string, items []string) (int, string, error) {
	if len(items) == 0 {
		return 0, "", errors.New("no items to select")
	}

	var result int
	var filter string

	width := 0
	for _, item := range items {
		value, _ := splitItem(item)
		width = max(width, utf8.RuneCountInString(value))
	}
	render := func(item string, selected bool) string {
		value, description := splitItem(item)
		if description != "" {
			value = pad(value, width)
		}
		// The selector sits in the mark column; items line up with the label.
		line := "  " + themeText(value)
		if selected {
			line = themeSuccess("› " + value)
		}
		if description != "" {
			line += "  " + themeMuted(description)
		}
		return line
	}

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
			fmt.Print("\r\033[K" + markAsk() + " " + themeAccent(bold(label)))
			if filter != "" {
				fmt.Print("  " + themeMuted("filter: ") + themeText(filter))
			}
			fmt.Print("\r\n")
			lineCount++

			// Print options
			for i, item := range filteredItems {
				fmt.Print("\r\033[K" + render(item, i == selectedIdx) + "\r\n")
				lineCount++
			}
			if len(filteredItems) == 0 {
				fmt.Print("\r\033[K  " + themeMuted("no matches") + "\r\n")
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
					value, _ := splitItem(items[result])
					answered(label, value)
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

	value, _ := splitItem(items[result])
	return result, value, nil
}

func MustSelect(label string, items []string) (int, string) {
	index, value, err := Select(label, items)

	if err != nil {
		Fatal(err)
	}

	return index, value
}
