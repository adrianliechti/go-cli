package cli

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Prompts and status lines start with a one-character mark in the same
// column: ? while asking, a spinner while working, ✓ when answered or done,
// ✗ on failure. Text after the mark, and Fields below it, line up.

func markAsk() string  { return themeAccent("?") }
func markDone() string { return themeSuccess("✓") }
func markFail() string { return themeError("✗") }

// answered replaces the current prompt line with its summary.
func answered(label, value string) {
	line := "\r\033[K" + markDone() + " "
	if label != "" {
		line += bold(label) + "  "
	}
	fmt.Print(line + themeText(value) + "\r\n")
}

// Done prints a finished step: "✓ text".
func Done(text string) {
	fmt.Println(markDone() + " " + themeText(text))
}

// Failed prints a failed step to standard error: "✗ text".
func Failed(text string) {
	os.Stderr.WriteString(markFail() + " " + themeError(text) + "\n")
}

// Fields prints label and value pairs aligned under the text of the line
// above, e.g. the details of a step printed with Done.
func Fields(rows [][2]string) {
	width := 0
	for _, r := range rows {
		width = max(width, utf8.RuneCountInString(r[0]))
	}
	for _, r := range rows {
		fmt.Println("  " + themeMuted(pad(r[0], width)) + "  " + themeText(r[1]))
	}
}

func pad(s string, width int) string {
	return s + strings.Repeat(" ", max(0, width-utf8.RuneCountInString(s)))
}

// splitItem separates a Select item's value from its optional description,
// given after a tab.
func splitItem(item string) (value, description string) {
	value, description, _ = strings.Cut(item, "\t")
	return value, description
}
