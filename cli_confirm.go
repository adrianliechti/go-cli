package cli

import (
	"fmt"
	"os"
)

func Confirm(label string, defaultValue bool) (bool, error) {
	var result bool

	err := withRawMode(func() error {
		redraw := func() {
			fmt.Print("\r\033[K")
			hint := "(y/N)"
			if defaultValue {
				hint = "(Y/n)"
			}
			fmt.Print(themeAccent(bold(label)) + " " + themeSubtle(hint) + themeAccent(": "))
		}

		redraw()

		for {
			key, char, err := readKey(os.Stdin)
			if err != nil {
				return err
			}

			switch key {
			case keyCtrlC:
				fmt.Print("\r\n")
				return ErrUserAborted

			case keyEnter:
				result = defaultValue
				if result {
					fmt.Print(themeSuccess("yes") + "\r\n")
				} else {
					fmt.Print(themeError("no") + "\r\n")
				}
				return nil

			default:
				if char == 'y' || char == 'Y' {
					result = true
					fmt.Print(themeSuccess("yes") + "\r\n")
					return nil
				}
				if char == 'n' || char == 'N' {
					result = false
					fmt.Print(themeError("no") + "\r\n")
					return nil
				}
			}
		}
	})

	if err != nil {
		return false, err
	}

	return result, nil
}

func MustConfirm(label string, defaultValue bool) bool {
	value, err := Confirm(label, defaultValue)

	if err != nil {
		Fatal(err)
	}

	return value
}
