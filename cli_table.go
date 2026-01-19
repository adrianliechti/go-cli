package cli

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Box drawing characters
const (
	boxTopLeft     = "┌"
	boxTopRight    = "┐"
	boxBottomLeft  = "└"
	boxBottomRight = "┘"
	boxHorizontal  = "─"
	boxVertical    = "│"
	boxTopTee      = "┬"
	boxBottomTee   = "┴"
	boxLeftTee     = "├"
	boxRightTee    = "┤"
	boxCross       = "┼"
)

func Table(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = utf8.RuneCountInString(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				// Sanitize: replace newlines for width calculation
				sanitized := strings.ReplaceAll(cell, "\n", " ")
				sanitized = strings.ReplaceAll(sanitized, "\r", "")
				w := utf8.RuneCountInString(sanitized)
				if w > colWidths[i] {
					colWidths[i] = w
				}
			}
		}
	}

	// Add padding
	for i := range colWidths {
		colWidths[i] += 2
	}

	// Build horizontal lines
	buildLine := func(left, mid, right, fill string) string {
		var sb strings.Builder
		sb.WriteString(left)
		for i, w := range colWidths {
			sb.WriteString(strings.Repeat(fill, w))
			if i < len(colWidths)-1 {
				sb.WriteString(mid)
			}
		}
		sb.WriteString(right)
		return sb.String()
	}

	topLine := buildLine(boxTopLeft, boxTopTee, boxTopRight, boxHorizontal)
	midLine := buildLine(boxLeftTee, boxCross, boxRightTee, boxHorizontal)
	bottomLine := buildLine(boxBottomLeft, boxBottomTee, boxBottomRight, boxHorizontal)

	// Build row - calculate padding before styling to handle ANSI codes correctly
	buildRow := func(cells []string, isHeader bool) string {
		var sb strings.Builder
		sb.WriteString(themeMuted(boxVertical))
		for i, w := range colWidths {
			cell := ""
			if i < len(cells) {
				// Sanitize cell: replace newlines with spaces for single-line display
				cell = strings.ReplaceAll(cells[i], "\n", " ")
				cell = strings.ReplaceAll(cell, "\r", "")
			}
			// Calculate padding based on visible width (before styling)
			visibleWidth := utf8.RuneCountInString(cell)
			padding := w - visibleWidth - 1
			if padding < 0 {
				padding = 0
			}

			// Build the cell: space + styled content + padding spaces
			sb.WriteString(" ")
			if isHeader {
				sb.WriteString(bold(themeAccent(cell)))
			} else {
				sb.WriteString(themeText(cell))
			}
			sb.WriteString(strings.Repeat(" ", padding))
			sb.WriteString(themeMuted(boxVertical))
		}
		return sb.String()
	}

	// Print table
	fmt.Println(themeMuted(topLine))
	fmt.Println(buildRow(headers, true))
	fmt.Println(themeMuted(midLine))
	for _, row := range rows {
		fmt.Println(buildRow(row, false))
	}
	fmt.Println(themeMuted(bottomLine))
}
