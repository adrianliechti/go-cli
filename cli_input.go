package cli

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func Input(label, placeholder string) (string, error) {
	var result string

	err := withRawMode(func() error {
		buffer := ""

		redraw := func() {
			clearLine()
			prompt := markAsk() + " " + themeAccent(bold(label)) + "  "
			if placeholder != "" && buffer == "" {
				fmt.Print(prompt + themeSubtle(placeholder))
			} else {
				fmt.Print(prompt + themeText(buffer))
			}
		}

		redraw()

		for {
			key, char, err := readKey(os.Stdin)
			if err != nil {
				return err
			}

			switch key {
			case keyCtrlC:
				fmt.Print("\r\033[K")
				return ErrUserAborted

			case keyEnter:
				if buffer == "" && placeholder != "" {
					buffer = placeholder
				}
				result = buffer
				answered(label, buffer)
				return nil

			case keyBackspace:
				if len(buffer) > 0 {
					_, size := utf8.DecodeLastRuneInString(buffer)
					buffer = buffer[:len(buffer)-size]
				}

			case keyCtrlU:
				buffer = ""

			default:
				if char != 0 && char >= 32 {
					buffer += string(char)
				}
			}

			redraw()
		}
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

func MustInput(label, placeholder string) string {
	value, err := Input(label, placeholder)

	if err != nil {
		Fatal(err)
	}

	return value
}
