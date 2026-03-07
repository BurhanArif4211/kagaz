package editor

import (
	"fyne.io/fyne/v2"
)

// Tapped moves the cursor to the tapped position.
func (e *NoteEditor) Tapped(ev *fyne.PointEvent) {
	row, col := e.grid.CursorLocationForPosition(ev.Position)
	e.cursor = e.rowColToIndex(row, col)
	e.selStart = -1 // clear selection
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

// TappedSecondary (long press) shows a context menu.
func (e *NoteEditor) TappedSecondary(ev *fyne.PointEvent) {
	// Find word under tap
	row, col := e.grid.CursorLocationForPosition(ev.Position)
	pos := e.rowColToIndex(row, col)

	// Expand to word boundaries (simplified: select the whole line for demo)
	lines := e.doc.Lines()
	lineIdx, _ := e.indexToLineCol(pos)
	if lineIdx >= 0 && lineIdx < len(lines) {
		line := lines[lineIdx]
		start := e.lineColToIndex(lineIdx, 0)
		end := e.lineColToIndex(lineIdx, len(line))
		e.selStart, e.selEnd = start, end
		e.Refresh()
	}

	// Show a popup (stub)
	// In a real app, use widget.NewPopUpMenu(...)
}

// Helper: convert grid row/col to absolute character index.
func (e *NoteEditor) rowColToIndex(row, col int) int {
	lines := e.doc.Lines()
	if row < 0 || row >= len(lines) {
		// Out of bounds – return last valid position
		return e.doc.len()
	}
	// Sum lengths of previous lines + col
	idx := 0
	for i := 0; i < row; i++ {
		idx += len(lines[i]) + 1 // +1 for newline character
	}
	idx += col
	// Ensure col not beyond line length
	if col > len(lines[row]) {
		idx = idx - col + len(lines[row]) // clamp to line end
	}
	return idx
}

// Helper: convert absolute index to row/col for cursor.
func (e *NoteEditor) indexToRowCol() (row, col int) {
	return e.indexToRowColFor(e.cursor)
}

// Helper: convert absolute index to row/col for any position.
func (e *NoteEditor) indexToRowColFor(pos int) (row, col int) {
	lines := e.doc.Lines()
	remaining := pos
	for i, line := range lines {
		lineLen := len(line)
		if remaining <= lineLen {
			return i, remaining
		}
		remaining -= lineLen + 1 // +1 for newline
	}
	// Beyond end: return last line end
	if len(lines) == 0 {
		return 0, 0
	}
	lastLine := lines[len(lines)-1]
	return len(lines) - 1, len(lastLine)
}

// Helper: convert line index and column to absolute index.
func (e *NoteEditor) lineColToIndex(line, col int) int {
	lines := e.doc.Lines()
	if line < 0 || line >= len(lines) {
		return e.doc.len()
	}
	idx := 0
	for i := 0; i < line; i++ {
		idx += len(lines[i]) + 1
	}
	idx += col
	if col > len(lines[line]) {
		idx = idx - col + len(lines[line])
	}
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
