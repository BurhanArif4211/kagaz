package editor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type NoteEditor struct {
	widget.BaseWidget
	doc      *Document
	grid     *widget.TextGrid
	cursor   int
	selStart int
	selEnd   int

	// Cached row/col for rendering
	cursorRow, cursorCol     int
	selStartRow, selStartCol int
	selEndRow, selEndCol     int

	// Handle positions (relative to editor)
	handleStartPos fyne.Position
	handleEndPos   fyne.Position
	showHandles    bool

	draggingHandle int // 0=none, 1=start, 2=end
}

const handleRadius = 20 // pixels for hit detection

func NewNoteEditor() *NoteEditor {
	e := &NoteEditor{
		doc:            NewDocument(),
		grid:           widget.NewTextGrid(),
		cursor:         0,
		selStart:       -1,
		selEnd:         -1,
		draggingHandle: 0,
	}
	e.grid.Scroll = fyne.ScrollBoth
	e.ExtendBaseWidget(e)
	return e
}

// SetContent replaces the document with new segments.
func (e *NoteEditor) SetContent(segments []TextSegment) {
	e.doc.segments = segments
	e.cursor = 0
	e.selStart = -1
	e.selEnd = -1
	e.Refresh()
}

func (e *NoteEditor) CreateRenderer() fyne.WidgetRenderer {
	handleStart := canvas.NewCircle(theme.PrimaryColor())
	handleStart.StrokeWidth = 2
	handleStart.StrokeColor = theme.BackgroundColor()
	handleStart.Resize(fyne.NewSize(handleRadius, handleRadius))
	handleStart.Hide()

	handleEnd := canvas.NewCircle(theme.PrimaryColor())
	handleEnd.StrokeWidth = 2
	handleEnd.StrokeColor = theme.BackgroundColor()
	handleEnd.Resize(fyne.NewSize(handleRadius, handleRadius))
	handleEnd.Hide()

	return &noteEditorRenderer{
		editor:      e,
		grid:        e.grid,
		handleStart: handleStart,
		handleEnd:   handleEnd,
	}
}

// Focus handling
func (e *NoteEditor) FocusGained() {
	e.Refresh() // ensure cursor is visible
}

func (e *NoteEditor) FocusLost() {
	e.Refresh() // maybe remove cursor? But we'll keep it for now.
}

// TypedRune handles character input.
func (e *NoteEditor) TypedRune(r rune) {
	// If there's a selection, delete it first.
	if e.selStart >= 0 && e.selEnd > e.selStart {
		e.doc.DeleteRange(e.selStart, e.selEnd)
		e.cursor = e.selStart
		e.selStart = -1
		e.selEnd = -1
	}

	// Insert the character with default style.
	// For now, use an empty style (will be applied later).
	style := TextStyle{}
	text := string(r)
	newCursor := e.doc.InsertText(e.cursor, text, style)
	e.cursor = newCursor
	e.Refresh()
}

// TypedKey handles special keys.
func (e *NoteEditor) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyBackspace:
		if e.selStart >= 0 && e.selEnd > e.selStart {
			// Delete selected range
			e.doc.DeleteRange(e.selStart, e.selEnd)
			e.cursor = e.selStart
			e.selStart = -1
			e.selEnd = -1
		} else if e.cursor > 0 {
			// Delete character before cursor
			e.doc.DeleteRange(e.cursor-1, e.cursor)
			e.cursor--
		}
		e.Refresh()

	case fyne.KeyDelete:
		if e.selStart >= 0 && e.selEnd > e.selStart {
			e.doc.DeleteRange(e.selStart, e.selEnd)
			e.cursor = e.selStart
			e.selStart = -1
			e.selEnd = -1
		} else if e.cursor < e.doc.len() {
			e.doc.DeleteRange(e.cursor, e.cursor+1)
			// cursor stays
		}
		e.Refresh()

	case fyne.KeyReturn, fyne.KeyEnter:
		// Insert newline
		if e.selStart >= 0 && e.selEnd > e.selStart {
			e.doc.DeleteRange(e.selStart, e.selEnd)
			e.cursor = e.selStart
			e.selStart = -1
			e.selEnd = -1
		}
		style := TextStyle{}
		newCursor := e.doc.InsertText(e.cursor, "\n", style)
		e.cursor = newCursor
		e.Refresh()

	case fyne.KeyLeft:
		e.moveCursor(-1, ev)
	case fyne.KeyRight:
		e.moveCursor(1, ev)
	case fyne.KeyUp:
		e.moveCursorLine(-1, ev)
	case fyne.KeyDown:
		e.moveCursorLine(1, ev)
	}
}

// moveCursor moves the cursor horizontally by delta.
func (e *NoteEditor) moveCursor(delta int, ev *fyne.KeyEvent) {
	newPos := e.cursor + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos > e.doc.len() {
		newPos = e.doc.len()
	}

	// If shift is held, adjust selection.
	if fyne.KeyModifier(ev.Physical.ScanCode) == fyne.KeyModifierShift {
		if e.selStart < 0 {
			// Start new selection
			e.selStart = e.cursor
			e.selEnd = newPos
		} else {
			// Extend selection
			if newPos < e.selStart {
				e.selStart = newPos
			} else if newPos > e.selEnd {
				e.selEnd = newPos
			} else {
				// If moving inside selection, maybe shrink? We'll keep simple: just move cursor and keep selection.
				// For now, we'll just update cursor without changing selection.
			}
		}
	} else {
		// No shift: clear selection
		e.selStart = -1
		e.selEnd = -1
	}
	e.cursor = newPos
	e.Refresh()
}

// moveCursorLine moves the cursor up/down by one line.
func (e *NoteEditor) moveCursorLine(delta int, ev *fyne.KeyEvent) {
	lines := e.doc.Lines()
	if len(lines) == 0 {
		return
	}
	// Get current line and column
	line, col := e.indexToLineCol(e.cursor)

	newLine := line + delta
	if newLine < 0 {
		newLine = 0
	}
	if newLine >= len(lines) {
		newLine = len(lines) - 1
	}
	// Clamp column to new line length
	lineLen := len(lines[newLine])
	if col > lineLen {
		col = lineLen
	}
	newPos := e.lineColToIndex(newLine, col)

	// Handle shift similar to moveCursor
	if fyne.KeyModifier(ev.Physical.ScanCode) == fyne.KeyModifierShift {
		if e.selStart < 0 {
			e.selStart = e.cursor
			e.selEnd = newPos
		} else {
			// Extend selection appropriately; we'll keep simple for now
			if newPos < e.selStart {
				e.selStart = newPos
			} else if newPos > e.selEnd {
				e.selEnd = newPos
			}
		}
	} else {
		e.selStart = -1
		e.selEnd = -1
	}
	e.cursor = newPos
	e.Refresh()
}
