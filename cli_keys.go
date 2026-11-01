package cli

import (
	"io"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

// Key codes
const (
	keyUnknown = iota
	keyEnter
	keyShiftEnter
	keyBackspace
	keyTab
	keyEscape
	keySpace
	keyUp
	keyDown
	keyLeft
	keyRight
	keyHome
	keyEnd
	keyDelete
	keyCtrlA
	keyCtrlC
	keyCtrlD
	keyCtrlE
	keyCtrlJ
	keyCtrlK
	keyCtrlU
	keyCtrlW
	keyAltEnter
)

// pending holds bytes read from the terminal but not yet returned as keys. A
// single read can carry several keys when input is pasted, typed quickly or
// arrives over a slow connection, and a key's bytes can be split across reads.
var pending []byte

// readKey returns the next key press and its rune, reading from r only when
// no complete key is pending.
func readKey(r io.Reader) (key int, char rune, err error) {
	for {
		if len(pending) > 0 {
			key, char, n := parseKey(pending)
			if n > 0 {
				pending = pending[n:]
				return key, char, nil
			}
		}
		buf := make([]byte, 256)
		n, err := r.Read(buf)
		if err != nil {
			return keyUnknown, 0, err
		}
		if n == 0 {
			continue
		}
		// A lone ESC is the Escape key; with more bytes it starts a sequence.
		if len(pending) == 0 && n == 1 && buf[0] == 27 {
			return keyEscape, 0, nil
		}
		pending = append(pending, buf[:n]...)
	}
}

// parseKey decodes the first key in buf and reports how many bytes it used;
// 0 means buf holds only the start of a key.
func parseKey(buf []byte) (key int, char rune, n int) {
	b := buf[0]
	switch b {
	case 1: // Ctrl+A
		return keyCtrlA, 0, 1
	case 3: // Ctrl+C
		return keyCtrlC, 0, 1
	case 4: // Ctrl+D
		return keyCtrlD, 0, 1
	case 5: // Ctrl+E
		return keyCtrlE, 0, 1
	case 8, 127: // Backspace: DEL on Unix terminals, BS on some Windows consoles
		return keyBackspace, 0, 1
	case 9: // Tab
		return keyTab, '\t', 1
	case 10: // Ctrl+J (line feed)
		return keyCtrlJ, '\n', 1
	case 11: // Ctrl+K
		return keyCtrlK, 0, 1
	case 13: // Enter (CR)
		return keyEnter, '\n', 1
	case 21: // Ctrl+U
		return keyCtrlU, 0, 1
	case 23: // Ctrl+W
		return keyCtrlW, 0, 1
	case 27:
		return parseEscapeSequence(buf)
	case 32: // Space
		return keySpace, ' ', 1
	}
	if b < 32 {
		return keyUnknown, 0, 1
	}
	if b < 0x80 {
		return keyUnknown, rune(b), 1
	}
	if !utf8.FullRune(buf) {
		return keyUnknown, 0, 0
	}
	r, size := utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		return keyUnknown, 0, size
	}
	return keyUnknown, r, size
}

// parseEscapeSequence decodes an ANSI escape sequence at the start of buf.
func parseEscapeSequence(buf []byte) (key int, char rune, n int) {
	if len(buf) < 2 {
		return keyEscape, 0, 1
	}

	switch buf[1] {
	case '[': // CSI: ESC [ parameters final, with a final byte from @ to ~
		end := -1
		for i := 2; i < len(buf); i++ {
			if buf[i] >= 0x40 && buf[i] <= 0x7E {
				end = i
				break
			}
		}
		if end < 0 {
			return keyUnknown, 0, 0
		}
		params := splitCSI(string(buf[2:end]))
		switch buf[end] {
		case 'A':
			return keyUp, 0, end + 1
		case 'B':
			return keyDown, 0, end + 1
		case 'C':
			return keyRight, 0, end + 1
		case 'D':
			return keyLeft, 0, end + 1
		case 'H':
			return keyHome, 0, end + 1
		case 'F':
			return keyEnd, 0, end + 1
		case '~':
			switch params[0] {
			case 1, 7:
				return keyHome, 0, end + 1
			case 3:
				return keyDelete, 0, end + 1
			case 4, 8:
				return keyEnd, 0, end + 1
			}
		case 'u': // CSI u: ESC [ code ; modifier u, e.g. Shift+Enter
			if len(params) >= 2 && params[0] == 13 {
				switch params[1] {
				case 2:
					return keyShiftEnter, '\n', end + 1
				case 3:
					return keyAltEnter, '\n', end + 1
				}
			}
		}
		return keyUnknown, 0, end + 1

	case 'O': // SS3: ESC O final, an alternate form of arrow keys
		if len(buf) < 3 {
			return keyUnknown, 0, 0
		}
		switch buf[2] {
		case 'A':
			return keyUp, 0, 3
		case 'B':
			return keyDown, 0, 3
		case 'C':
			return keyRight, 0, 3
		case 'D':
			return keyLeft, 0, 3
		case 'H':
			return keyHome, 0, 3
		case 'F':
			return keyEnd, 0, 3
		}
		return keyUnknown, 0, 3

	case 13: // Alt+Enter in terminals that prefix ESC
		return keyAltEnter, '\n', 2
	}

	// Escape followed by an ordinary key: Escape now, the key next.
	return keyEscape, 0, 1
}

// splitCSI splits a CSI parameter string like "13;2" into integers
func splitCSI(s string) []int {
	var result []int
	var current int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			current = current*10 + int(c-'0')
		} else if c == ';' {
			result = append(result, current)
			current = 0
		}
	}
	result = append(result, current)
	return result
}

// withRawMode executes a function with the terminal in raw mode
func withRawMode(fn func() error) error {
	fd := int(os.Stdin.Fd())

	// Save current terminal state
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	return fn()
}

// isTerminal returns true if stdout is a terminal
func isTerminalCheck() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
