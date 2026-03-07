package editor

import (
	"unicode"

	"fyne.io/fyne/v2"
)

// Tapped moves the cursor to the tapped position.
func (e *NoteEditor) Tapped(ev *fyne.PointEvent) {
	row, col := e.grid.CursorLocationForPosition(ev.Position)
	e.cursor = e.rowColToIndex(row, col)
	e.selStart = -1 // clear selection
	if c := fyne.CurrentApp().Driver().CanvasForObject(e); c != nil {
		c.Focus(e)
	}
	e.Refresh()
}

// Dragged handles selection.
func (e *NoteEditor) Dragged(ev *fyne.DragEvent) {
	// For prototype, we'll treat any drag as selection (simplified).
	// In a real app, you'd differentiate scroll vs selection.
	row, col := e.grid.CursorLocationForPosition(ev.Position)
	newPos := e.rowColToIndex(row, col)

	if e.selStart < 0 {
		// Start a new selection
		e.selStart = e.cursor
	}
	e.selEnd = newPos
	if e.selEnd < e.selStart {
		e.selStart, e.selEnd = e.selEnd, e.selStart
	}
	e.Refresh()
}

// DragEnd completes the selection.
func (e *NoteEditor) DragEnd() {
	// Nothing needed; selection remains.
}

// TappedSecondary (long press) selects the word under the finger and shows the context menu.
func (e *NoteEditor) TappedSecondary(ev *fyne.PointEvent) {
	// Get the character position under the tap
	row, col := e.grid.CursorLocationForPosition(ev.Position)
	pos := e.rowColToIndex(row, col)

	// Expand to word boundaries
	start, end := e.expandToWord(pos)
	if start < end {
		e.selStart, e.selEnd = start, end
		e.cursor = end // optional: place cursor at end of selection
	} else {
		// If no word, just place cursor (should not happen on long press, but fallback)
		e.selStart, e.selEnd = -1, -1
		e.cursor = pos
	}

	e.Refresh()

	// Show context menu
	e.showContextMenu(ev.AbsolutePosition)
}

// expandToWord returns the start and end indices of the word containing pos.
// Returns pos, pos if no word found.
func (e *NoteEditor) expandToWord(pos int) (int, int) {
	lines := e.doc.Lines()
	if len(lines) == 0 {
		return pos, pos
	}

	lineIdx, col := e.indexToLineCol(pos)
	if lineIdx < 0 || lineIdx >= len(lines) {
		return pos, pos
	}

	line := lines[lineIdx]
	if col < 0 || col > len(line) {
		return pos, pos
	}

	// Find word start (backwards until non-letter/digit or line start)
	start := col
	for start > 0 && (unicode.IsLetter(rune(line[start-1])) || unicode.IsDigit(rune(line[start-1]))) {
		start--
	}

	// Find word end (forwards until non-letter/digit or line end)
	end := col
	for end < len(line) && (unicode.IsLetter(rune(line[end])) || unicode.IsDigit(rune(line[end]))) {
		end++
	}

	// Convert to absolute indices
	absStart := e.lineColToIndex(lineIdx, start)
	absEnd := e.lineColToIndex(lineIdx, end)

	return absStart, absEnd
}

// Helper: convert grid row/col to absolute character index.
func (e *NoteEditor) rowColToIndex(row, col int) int {
	lines := e.doc.Lines()
	if len(lines) == 0 {
		return 0
	}
	if row < 0 {
		row = 0
	}
	if row >= len(lines) {
		// return total length
		total := 0
		for _, line := range lines {
			total += len(line) + 1
		}
		if total > 0 {
			total-- // last newline? but we want last char index
		}
		return total
	}
	idx := 0
	for i := 0; i < row; i++ {
		idx += len(lines[i]) + 1
	}
	if col < 0 {
		col = 0
	}
	if col > len(lines[row]) {
		col = len(lines[row])
	}
	idx += col
	return idx
}

// Helper: convert absolute index to row/col for cursor.
func (e *NoteEditor) indexToRowCol() (row, col int) {
	return e.indexToRowColFor(e.cursor)
}

// Helper: convert absolute index to row/col for any position.
func (e *NoteEditor) indexToRowColFor(pos int) (row, col int) {
	lines := e.doc.Lines()
	if len(lines) == 0 {
		return 0, 0
	}
	if pos < 0 {
		pos = 0
	}
	remaining := pos
	for i, line := range lines {
		lineLen := len(line)
		if remaining <= lineLen {
			return i, remaining
		}
		remaining -= lineLen + 1
	}
	lastRow := len(lines) - 1
	return lastRow, len(lines[lastRow])
}

// Helper: convert line index and column to absolute index.
func (e *NoteEditor) lineColToIndex(line, col int) int {
	lines := e.doc.Lines()
	if len(lines) == 0 {
		return 0
	}
	if line < 0 {
		line = 0
	}
	if line >= len(lines) {
		line = len(lines) - 1
	}
	idx := 0
	for i := 0; i < line; i++ {
		idx += len(lines[i]) + 1
	}
	if col < 0 {
		col = 0
	}
	if col > len(lines[line]) {
		col = len(lines[line])
	}
	idx += col
	return idx
}

// Helper: get line index and column from absolute index.
func (e *NoteEditor) indexToLineCol(pos int) (line, col int) {
	lines := e.doc.Lines()
	remaining := pos
	for i, lineStr := range lines {
		lineLen := len(lineStr)
		if remaining <= lineLen {
			return i, remaining
		}
		remaining -= lineLen + 1
	}
	if len(lines) == 0 {
		return 0, 0
	}
	return len(lines) - 1, len(lines[len(lines)-1])
}
