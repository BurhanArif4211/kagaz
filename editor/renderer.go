package editor

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type noteEditorRenderer struct {
	editor *NoteEditor
	grid   *widget.TextGrid
}

// Define styles for cursor and selection
type gridStyle struct {
	fg color.Color
	bg color.Color
}

func (s *gridStyle) TextColor() color.Color       { return s.fg }
func (s *gridStyle) BackgroundColor() color.Color { return s.bg }

// var cursorStyle = &gridStyle{
// 	fg: theme.BackgroundColor(),
// 	bg: theme.ForegroundColor(),
// }

// var selectionStyle = &gridStyle{
// 	fg: theme.ForegroundColor(),
// 	bg: theme.PrimaryColor(),
// }

func (r *noteEditorRenderer) Layout(size fyne.Size) {
	r.grid.Resize(size)
}

func (r *noteEditorRenderer) MinSize() fyne.Size {
	return r.grid.MinSize()
}

func (r *noteEditorRenderer) Refresh() {
	// Update the grid's plain text from document lines.
	lines := r.editor.doc.Lines()
	r.grid.SetText(strings.Join(lines, "\n"))

	// Convert cursor and selection to row/col
	if r.editor.cursor >= 0 {
		r.editor.cursorRow, r.editor.cursorCol = r.editor.indexToRowColFor(r.editor.cursor)
	} else {
		r.editor.cursorRow, r.editor.cursorCol = -1, -1
	}
	if r.editor.selStart >= 0 && r.editor.selEnd > r.editor.selStart {
		r.editor.selStartRow, r.editor.selStartCol = r.editor.indexToRowColFor(r.editor.selStart)
		r.editor.selEndRow, r.editor.selEndCol = r.editor.indexToRowColFor(r.editor.selEnd)
	} else {
		r.editor.selStartRow, r.editor.selStartCol = -1, -1
		r.editor.selEndRow, r.editor.selEndCol = -1, -1
	}

	// Apply selection style first (cursor may override part of it)
	// if r.editor.selStartRow >= 0 && r.editor.selEndRow >= 0 {
	// 	r.grid.SetStyleRange(
	// 		r.editor.selStartRow, r.editor.selStartCol,
	// 		r.editor.selEndRow, r.editor.selEndCol,
	// 		selectionStyle,
	// 	)
	// }

	// Apply cursor style (if cursor is within selection, selection style will be overridden)
	// if r.editor.cursorRow >= 0 && r.editor.cursorCol >= 0 {
	// 	r.grid.SetStyle(r.editor.cursorRow, r.editor.cursorCol, cursorStyle)
	// }

	// Force the grid to redraw
	r.grid.Refresh()
}

func (r *noteEditorRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.grid}
}

func (r *noteEditorRenderer) Destroy() {}
