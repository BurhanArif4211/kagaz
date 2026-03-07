package editor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// NoteEditor is the custom text editing widget.
type NoteEditor struct {
	widget.BaseWidget
	doc      *Document
	grid     *widget.TextGrid
	cursor   int // absolute character index
	selStart int // -1 if no selection
	selEnd   int

	// Cached row/col for rendering
	cursorRow, cursorCol     int
	selStartRow, selStartCol int
	selEndRow, selEndCol     int
}

// NewNoteEditor creates a new editor with empty content.
func NewNoteEditor() *NoteEditor {
	e := &NoteEditor{
		doc:      NewDocument(),
		grid:     widget.NewTextGrid(),
		cursor:   0,
		selStart: -1,
		selEnd:   -1,
	}
	e.grid.Scroll = fyne.ScrollBoth // enable scrolling
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

// CreateRenderer implements fyne.Widget.
func (e *NoteEditor) CreateRenderer() fyne.WidgetRenderer {
	return &noteEditorRenderer{editor: e, grid: e.grid}
}
