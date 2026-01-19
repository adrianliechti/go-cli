package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Theme represents a color theme
type Theme struct {
	// Base colors
	Rosewater RGB
	Flamingo  RGB
	Pink      RGB
	Mauve     RGB
	Red       RGB
	Maroon    RGB
	Peach     RGB
	Yellow    RGB
	Green     RGB
	Teal      RGB
	Sky       RGB
	Sapphire  RGB
	Blue      RGB
	Lavender  RGB

	// Text colors
	Text     RGB
	Subtext1 RGB
	Subtext0 RGB

	// Overlay colors
	Overlay2 RGB
	Overlay1 RGB
	Overlay0 RGB

	// Surface colors
	Surface2 RGB
	Surface1 RGB
	Surface0 RGB

	// Background colors
	Base   RGB
	Mantle RGB
	Crust  RGB
}

// RGB represents an RGB color
type RGB struct {
	R, G, B uint8
}

// Hex creates an RGB from a hex string like "#1e1e2e"
func Hex(hex string) RGB {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return RGB{}
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return RGB{uint8(r), uint8(g), uint8(b)}
}

// Catppuccin Mocha (dark theme)
var Mocha = Theme{
	Rosewater: Hex("#f5e0dc"),
	Flamingo:  Hex("#f2cdcd"),
	Pink:      Hex("#f5c2e7"),
	Mauve:     Hex("#cba6f7"),
	Red:       Hex("#f38ba8"),
	Maroon:    Hex("#eba0ac"),
	Peach:     Hex("#fab387"),
	Yellow:    Hex("#f9e2af"),
	Green:     Hex("#a6e3a1"),
	Teal:      Hex("#94e2d5"),
	Sky:       Hex("#89dceb"),
	Sapphire:  Hex("#74c7ec"),
	Blue:      Hex("#89b4fa"),
	Lavender:  Hex("#b4befe"),
	Text:      Hex("#cdd6f4"),
	Subtext1:  Hex("#bac2de"),
	Subtext0:  Hex("#a6adc8"),
	Overlay2:  Hex("#9399b2"),
	Overlay1:  Hex("#7f849c"),
	Overlay0:  Hex("#6c7086"),
	Surface2:  Hex("#585b70"),
	Surface1:  Hex("#45475a"),
	Surface0:  Hex("#313244"),
	Base:      Hex("#1e1e2e"),
	Mantle:    Hex("#181825"),
	Crust:     Hex("#11111b"),
}

// Catppuccin Latte (light theme)
var Latte = Theme{
	Rosewater: Hex("#dc8a78"),
	Flamingo:  Hex("#dd7878"),
	Pink:      Hex("#ea76cb"),
	Mauve:     Hex("#8839ef"),
	Red:       Hex("#d20f39"),
	Maroon:    Hex("#e64553"),
	Peach:     Hex("#fe640b"),
	Yellow:    Hex("#df8e1d"),
	Green:     Hex("#40a02b"),
	Teal:      Hex("#179299"),
	Sky:       Hex("#04a5e5"),
	Sapphire:  Hex("#209fb5"),
	Blue:      Hex("#1e66f5"),
	Lavender:  Hex("#7287fd"),
	Text:      Hex("#4c4f69"),
	Subtext1:  Hex("#5c5f77"),
	Subtext0:  Hex("#6c6f85"),
	Overlay2:  Hex("#7c7f93"),
	Overlay1:  Hex("#8c8fa1"),
	Overlay0:  Hex("#9ca0b0"),
	Surface2:  Hex("#acb0be"),
	Surface1:  Hex("#bcc0cc"),
	Surface0:  Hex("#ccd0da"),
	Base:      Hex("#eff1f5"),
	Mantle:    Hex("#e6e9ef"),
	Crust:     Hex("#dce0e8"),
}

// Current theme (default to Mocha)
var currentTheme = Mocha

// Color capability detection
type colorMode int

const (
	colorModeNone colorMode = iota
	colorMode256
	colorModeTrueColor
)

var detectedColorMode colorMode

func init() {
	detectedColorMode = detectColorMode()
}

func detectColorMode() colorMode {
	// Check for true color support
	colorterm := os.Getenv("COLORTERM")
	if colorterm == "truecolor" || colorterm == "24bit" {
		return colorModeTrueColor
	}

	// Check TERM for common true color terminals
	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") || strings.Contains(term, "24bit") {
		return colorModeTrueColor
	}

	// Check for known terminals that support true color
	termProgram := os.Getenv("TERM_PROGRAM")
	switch termProgram {
	case "iTerm.app", "Apple_Terminal", "Hyper", "vscode":
		return colorModeTrueColor
	}

	// Windows Terminal supports true color
	if os.Getenv("WT_SESSION") != "" {
		return colorModeTrueColor
	}

	// Check if any color is supported
	if term != "" && term != "dumb" {
		return colorMode256
	}

	return colorModeNone
}

// SetTheme sets the current color theme
func SetTheme(theme Theme) {
	currentTheme = theme
}

// SetDarkTheme sets the theme to Catppuccin Mocha
func SetDarkTheme() {
	currentTheme = Mocha
}

// SetLightTheme sets the theme to Catppuccin Latte
func SetLightTheme() {
	currentTheme = Latte
}

// GetTheme returns the current theme
func GetTheme() Theme {
	return currentTheme
}

// ANSI escape codes
const (
	escReset     = "\033[0m"
	escBold      = "\033[1m"
	escDim       = "\033[2m"
	escItalic    = "\033[3m"
	escUnderline = "\033[4m"
	escBlink     = "\033[5m"
	escReverse   = "\033[7m"
	escHidden    = "\033[8m"
	escStrike    = "\033[9m"
)

// Cursor control
const (
	escHideCursor = "\033[?25l"
	escShowCursor = "\033[?25h"
	escClearLine  = "\033[2K"
	escClearRight = "\033[K"
	escMoveUp     = "\033[%dA"
	escMoveDown   = "\033[%dB"
	escMoveRight  = "\033[%dC"
	escMoveLeft   = "\033[%dD"
	escMoveTo     = "\033[%d;%dH"
	escSaveCursor = "\033[s"
	escRestCursor = "\033[u"
)

// rgbToAnsi256 converts RGB to the nearest 256-color ANSI code
func rgbToAnsi256(r, g, b uint8) int {
	// Check for grayscale
	if r == g && g == b {
		if r < 8 {
			return 16
		}
		if r > 248 {
			return 231
		}
		return int((float64(r)-8)/247*24) + 232
	}

	// Convert to 6x6x6 color cube
	ri := int(float64(r) / 255 * 5)
	gi := int(float64(g) / 255 * 5)
	bi := int(float64(b) / 255 * 5)
	return 16 + 36*ri + 6*gi + bi
}

// Color applies foreground color to text
func (c RGB) Color(text string) string {
	if detectedColorMode == colorModeNone {
		return text
	}

	var colorCode string
	if detectedColorMode == colorModeTrueColor {
		colorCode = fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
	} else {
		colorCode = fmt.Sprintf("\033[38;5;%dm", rgbToAnsi256(c.R, c.G, c.B))
	}
	return colorCode + text + escReset
}

// Bg applies background color to text
func (c RGB) Bg(text string) string {
	if detectedColorMode == colorModeNone {
		return text
	}

	var colorCode string
	if detectedColorMode == colorModeTrueColor {
		colorCode = fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
	} else {
		colorCode = fmt.Sprintf("\033[48;5;%dm", rgbToAnsi256(c.R, c.G, c.B))
	}
	return colorCode + text + escReset
}

// Style helpers using current theme
func bold(text string) string {
	return escBold + text + escReset
}

func dim(text string) string {
	return escDim + text + escReset
}

func italic(text string) string {
	return escItalic + text + escReset
}

func underline(text string) string {
	return escUnderline + text + escReset
}

// Themed style helpers
func themeText(text string) string {
	return currentTheme.Text.Color(text)
}

func themeSubtle(text string) string {
	return currentTheme.Subtext0.Color(text)
}

func themeMuted(text string) string {
	return currentTheme.Overlay0.Color(text)
}

func themeAccent(text string) string {
	return currentTheme.Blue.Color(text)
}

func themeSuccess(text string) string {
	return currentTheme.Green.Color(text)
}

func themeWarning(text string) string {
	return currentTheme.Yellow.Color(text)
}

func themeError(text string) string {
	return currentTheme.Red.Color(text)
}

func themeHighlight(text string) string {
	return currentTheme.Mauve.Color(text)
}

// Cursor helpers
func hideCursor() {
	fmt.Print(escHideCursor)
}

func showCursor() {
	fmt.Print(escShowCursor)
}

func clearLine() {
	fmt.Print("\r" + escClearLine)
}

func clearRight() {
	fmt.Print(escClearRight)
}

func moveUp(n int) {
	if n > 0 {
		fmt.Printf(escMoveUp, n)
	}
}

func moveDown(n int) {
	if n > 0 {
		fmt.Printf(escMoveDown, n)
	}
}

func moveRight(n int) {
	if n > 0 {
		fmt.Printf(escMoveRight, n)
	}
}

func moveLeft(n int) {
	if n > 0 {
		fmt.Printf(escMoveLeft, n)
	}
}

func saveCursor() {
	fmt.Print(escSaveCursor)
}

func restoreCursor() {
	fmt.Print(escRestCursor)
}
