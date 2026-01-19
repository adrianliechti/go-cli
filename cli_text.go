package cli

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func Text(label, placeholder string) (string, error) {
	var result string

	err := withRawMode(func() error {
		lines := []string{""}
		currentLine := 0
		lastLineCount := 0

		// Hide cursor during editing
		fmt.Print(escHideCursor)
		defer fmt.Print(escShowCursor)

		clearPrevious := func() {
			for i := 0; i < lastLineCount; i++ {
				fmt.Print("\033[A")   // Move up
				fmt.Print("\r\033[K") // Clear line
			}
		}

		redraw := func() {
			lineCount := 0

			// Print label and hint
			if label != "" {
				fmt.Print("\r\033[K" + themeAccent(bold(label)) + " " + themeSubtle("(Ctrl+D to submit)") + "\r\n")
				lineCount++
			}

			// Print lines
			for i, line := range lines {
				lineNum := themeMuted(fmt.Sprintf("%2d │ ", i+1))
				fmt.Print("\r\033[K")
				if i == currentLine {
					fmt.Print(lineNum + themeText(line) + themeSubtle("█"))
				} else {
					fmt.Print(lineNum + themeText(line))
				}
				fmt.Print("\r\n")
				lineCount++
			}

			lastLineCount = lineCount
		}

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

			case keyCtrlD:
				result = strings.Join(lines, "\n")
				clearPrevious()
				if label != "" {
					fmt.Print("\r\033[K" + themeAccent(bold(label)) + "\r\n")
				}
				preview := strings.Join(lines, " ")
				if len(preview) > 60 {
					preview = preview[:57] + "..."
				}
				fmt.Print("\r\033[K" + themeSuccess("> ") + themeText(preview) + "\r\n")
				return nil

			case keyEnter:
				lines = append(lines, "")
				currentLine = len(lines) - 1

			case keyBackspace:
				if len(lines[currentLine]) > 0 {
					_, size := utf8.DecodeLastRuneInString(lines[currentLine])
					lines[currentLine] = lines[currentLine][:len(lines[currentLine])-size]
				} else if currentLine > 0 {
					// Join with previous line
					lines = append(lines[:currentLine], lines[currentLine+1:]...)
					currentLine--
				}

			case keyUp:
				if currentLine > 0 {
					currentLine--
				}

			case keyDown:
				if currentLine < len(lines)-1 {
					currentLine++
				}

			default:
				if char != 0 && char >= 32 {
					lines[currentLine] += string(char)
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

func MustText(label, placeholder string) string {
	value, err := Text(label, placeholder)

	if err != nil {
		Fatal(err)
	}

	return value
}
