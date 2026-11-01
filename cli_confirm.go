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
			fmt.Print(markAsk() + " " + themeAccent(bold(label)) + " " + themeSubtle(hint) + "  ")
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
				result = defaultValue
				answered(label, yesNo(result))
				return nil

			default:
				if char == 'y' || char == 'Y' {
					result = true
					answered(label, yesNo(result))
					return nil
				}
				if char == 'n' || char == 'N' {
					result = false
					answered(label, yesNo(result))
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

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
