package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/stretchr/testify/assert"
)

// tableLineWidth returns the number of runes in the first line of rendered output.
// Box-drawing characters (─, │, ┌, etc.) each occupy one terminal column.
func tableLineWidth(rendered string) int {
	line, _, _ := strings.Cut(rendered, "\n")
	return len([]rune(line))
}

// TestSummaryColumnsClampedToTerminalWidth verifies that the plain-text summary
// table produced by renderSummaryTablePlain never exceeds the requested terminal
// width.  Prior to the fix, clampColumns underestimated the per-column rendering
// overhead (2 instead of 3 chars, missing the leading '│'), so the table could
// be up to len(columns)+1 = 9 characters wider than the terminal.
func TestSummaryColumnsClampedToTerminalWidth(t *testing.T) {
	// Width chosen so the table CAN fit after correct shrinkage, but
	// overflows by 9 chars when the padding formula is wrong.
	const terminalWidth = 145

	rows := []table.Row{emptySummaryRow()}
	columns := summaryColumns(terminalWidth, rows)
	rendered := renderSummaryTablePlain(columns, rows)

	lineWidth := tableLineWidth(rendered)
	assert.LessOrEqual(t, lineWidth, terminalWidth,
		"table is %d chars wide, exceeds terminal width %d", lineWidth, terminalWidth)
}
