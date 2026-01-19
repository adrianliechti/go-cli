package cli

import (
	"io"
	"os"

	"golang.org/x/term"
)

// Key codes
const (
	keyUnknown = iota
	keyEnter
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
	keyCtrlK
	keyCtrlU
	keyCtrlW
)

// readKey reads a single key press and returns the key code and rune
func readKey(r io.Reader) (key int, char rune, err error) {
	buf := make([]byte, 4)
	n, err := r.Read(buf)
	if err != nil {
		return keyUnknown, 0, err
	}

	if n == 0 {
		return keyUnknown, 0, nil
	}

	b := buf[0]

	// Control characters
	switch b {
	case 1: // Ctrl+A
		return keyCtrlA, 0, nil
	case 3: // Ctrl+C
		return keyCtrlC, 0, nil
	case 4: // Ctrl+D
		return keyCtrlD, 0, nil
	case 5: // Ctrl+E
		return keyCtrlE, 0, nil
	case 9: // Tab
		return keyTab, '\t', nil
	case 11: // Ctrl+K
		return keyCtrlK, 0, nil
	case 13, 10: // Enter (CR or LF)
		return keyEnter, '\n', nil
	case 21: // Ctrl+U
		return keyCtrlU, 0, nil
	case 23: // Ctrl+W
		return keyCtrlW, 0, nil
	case 27: // Escape sequence
		if n == 1 {
			return keyEscape, 0, nil
		}
		return parseEscapeSequence(buf[:n])
	case 32: // Space
		return keySpace, ' ', nil
	case 127: // Backspace (DEL)
		return keyBackspace, 0, nil
	}

	// Regular ASCII characters
	if b >= 32 && b < 127 {
		return keyUnknown, rune(b), nil
	}

	// UTF-8 multi-byte characters
	if b >= 0xC0 {
		r := decodeUTF8(buf[:n])
		return keyUnknown, r, nil
	}

	return keyUnknown, 0, nil
}

// parseEscapeSequence parses ANSI escape sequences
func parseEscapeSequence(buf []byte) (key int, char rune, err error) {
	if len(buf) < 2 {
		return keyEscape, 0, nil
	}

	// CSI sequences (ESC [)
	if buf[1] == '[' {
		if len(buf) < 3 {
			return keyEscape, 0, nil
		}

		switch buf[2] {
		case 'A':
			return keyUp, 0, nil
		case 'B':
			return keyDown, 0, nil
		case 'C':
			return keyRight, 0, nil
		case 'D':
			return keyLeft, 0, nil
		case 'H':
			return keyHome, 0, nil
		case 'F':
			return keyEnd, 0, nil
		case '3':
			if len(buf) > 3 && buf[3] == '~' {
				return keyDelete, 0, nil
			}
		case '1':
			if len(buf) > 3 && buf[3] == '~' {
				return keyHome, 0, nil
			}
		case '4':
			if len(buf) > 3 && buf[3] == '~' {
				return keyEnd, 0, nil
			}
		}
	}

	// SS3 sequences (ESC O) - alternate arrow keys
	if buf[1] == 'O' && len(buf) >= 3 {
		switch buf[2] {
		case 'A':
			return keyUp, 0, nil
		case 'B':
			return keyDown, 0, nil
		case 'C':
			return keyRight, 0, nil
		case 'D':
			return keyLeft, 0, nil
		case 'H':
			return keyHome, 0, nil
		case 'F':
			return keyEnd, 0, nil
		}
	}

	return keyEscape, 0, nil
}

// decodeUTF8 decodes a UTF-8 sequence from bytes
func decodeUTF8(buf []byte) rune {
	if len(buf) == 0 {
		return 0
	}

	b := buf[0]

	// 1-byte (ASCII)
	if b < 0x80 {
		return rune(b)
	}

	// 2-byte
	if b < 0xE0 && len(buf) >= 2 {
		return rune(b&0x1F)<<6 | rune(buf[1]&0x3F)
	}

	// 3-byte
	if b < 0xF0 && len(buf) >= 3 {
		return rune(b&0x0F)<<12 | rune(buf[1]&0x3F)<<6 | rune(buf[2]&0x3F)
	}

	// 4-byte
	if len(buf) >= 4 {
		return rune(b&0x07)<<18 | rune(buf[1]&0x3F)<<12 | rune(buf[2]&0x3F)<<6 | rune(buf[3]&0x3F)
	}

	return 0
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
