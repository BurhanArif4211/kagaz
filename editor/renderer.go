package editor

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type noteEditorRenderer struct {
	editor      *NoteEditor
	grid        *widget.TextGrid
	handleStart *canvas.Circle
	handleEnd   *canvas.Circle
}
type cursorStyle struct{}

// selectionStyle implements widget.TextGridStyle for the selection highlight.
type selectionStyle struct{}

var (
	cursorStyleInstance    = &cursorStyle{}
	selectionStyleInstance = &selectionStyle{}
)

func (c *cursorStyle) Style() fyne.TextStyle {
	return fyne.TextStyle{} // no font style changes
}
func (c *cursorStyle) TextColor() color.Color {
	return theme.BackgroundColor()
}
func (c *cursorStyle) BackgroundColor() color.Color {
	return theme.ForegroundColor()
}

func (s *selectionStyle) Style() fyne.TextStyle {
	return fyne.TextStyle{}
}
func (s *selectionStyle) TextColor() color.Color {
	return theme.ForegroundColor()
}
func (s *selectionStyle) BackgroundColor() color.Color {
	return theme.PrimaryColor()
}
func (r *noteEditorRenderer) Layout(size fyne.Size) {
	r.grid.Resize(size)
	// Handles are positioned in Refresh, not here.
}

func (r *noteEditorRenderer) MinSize() fyne.Size {
	return r.grid.MinSize()
}

func (r *noteEditorRenderer) Refresh() {
	// Update grid text
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

	// Apply selection style
	if r.editor.selStartRow >= 0 && r.editor.selEndRow >= 0 {
		r.grid.SetStyleRange(
			r.editor.selStartRow, r.editor.selStartCol,
			r.editor.selEndRow, r.editor.selEndCol,
			selectionStyleInstance,
		)
	}

	// Apply cursor style
	if r.editor.cursorRow >= 0 && r.editor.cursorCol >= 0 {
		r.grid.SetStyle(r.editor.cursorRow, r.editor.cursorCol, cursorStyleInstance)
	}

	// Update handle positions and visibility
	r.updateHandles()

	// Refresh grid and handles
	r.grid.Refresh()
	if r.handleStart.Visible() {
		r.handleStart.Refresh()
	}
	if r.handleEnd.Visible() {
		r.handleEnd.Refresh()
	}
}

func (r *noteEditorRenderer) updateHandles() {
	// Determine if we should show handles (selection exists and not empty)
	show := r.editor.selStart >= 0 && r.editor.selEnd > r.editor.selStart
	r.editor.showHandles = show

	if !show {
		r.handleStart.Hide()
		r.handleEnd.Hide()
		return
	}

	// Get cell positions for start and end
	startPos := r.grid.PositionForCursorLocation(r.editor.selStartRow, r.editor.selStartCol)
	endPos := r.grid.PositionForCursorLocation(r.editor.selEndRow, r.editor.selEndCol)

	// Center of the cell (assuming cell size is roughly 20x20)
	// We don't have cell size, so we'll just use top-left plus a small offset.
	// For better accuracy, we could measure text size, but this is fine for hit detection.
	cellOffset := fyne.NewPos(10, 10) // half of assumed cell size
	r.editor.handleStartPos = startPos.Add(cellOffset)
	r.editor.handleEndPos = endPos.Add(cellOffset)

	// Position the circles (centered at the calculated points)
	r.handleStart.Move(r.editor.handleStartPos.Subtract(fyne.NewPos(handleRadius, handleRadius)))
	r.handleEnd.Move(r.editor.handleEndPos.Subtract(fyne.NewPos(handleRadius, handleRadius)))

	r.handleStart.Show()
	r.handleEnd.Show()
}

func (r *noteEditorRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.grid, r.handleStart, r.handleEnd}
}

func (r *noteEditorRenderer) Destroy() {}

// package editor

// import (
// 	"image/color"
// 	"strings"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/theme"
// 	"fyne.io/fyne/v2/widget"
// )

// type noteEditorRenderer struct {
// 	editor *NoteEditor
// 	grid   *widget.TextGrid
// }

// // cursorStyle implements widget.TextGridStyle for the cursor.
// type cursorStyle struct{}

// func (c *cursorStyle) Style() fyne.TextStyle {
// 	return fyne.TextStyle{} // no font style changes
// }
// func (c *cursorStyle) TextColor() color.Color {
// 	return theme.BackgroundColor()
// }
// func (c *cursorStyle) BackgroundColor() color.Color {
// 	return theme.ForegroundColor()
// }

// // selectionStyle implements widget.TextGridStyle for the selection highlight.
// type selectionStyle struct{}

// func (s *selectionStyle) Style() fyne.TextStyle {
// 	return fyne.TextStyle{}
// }
// func (s *selectionStyle) TextColor() color.Color {
// 	return theme.ForegroundColor()
// }
// func (s *selectionStyle) BackgroundColor() color.Color {
// 	return theme.PrimaryColor()
// }

// var (
// 	cursorStyleInstance    = &cursorStyle{}
// 	selectionStyleInstance = &selectionStyle{}
// )

// func (r *noteEditorRenderer) Layout(size fyne.Size) {
// 	r.grid.Resize(size)
// }

// func (r *noteEditorRenderer) MinSize() fyne.Size {
// 	return r.grid.MinSize()
// }

// func (r *noteEditorRenderer) Refresh() {
// 	// Update the grid's plain text from document lines.
// 	lines := r.editor.doc.Lines()
// 	r.grid.SetText(strings.Join(lines, "\n"))

// 	// Convert cursor and selection to row/col
// 	if r.editor.cursor >= 0 {
// 		r.editor.cursorRow, r.editor.cursorCol = r.editor.indexToRowColFor(r.editor.cursor)
// 	} else {
// 		r.editor.cursorRow, r.editor.cursorCol = -1, -1
// 	}
// 	if r.editor.selStart >= 0 && r.editor.selEnd > r.editor.selStart {
// 		r.editor.selStartRow, r.editor.selStartCol = r.editor.indexToRowColFor(r.editor.selStart)
// 		r.editor.selEndRow, r.editor.selEndCol = r.editor.indexToRowColFor(r.editor.selEnd)
// 	} else {
// 		r.editor.selStartRow, r.editor.selStartCol = -1, -1
// 		r.editor.selEndRow, r.editor.selEndCol = -1, -1
// 	}

// 	// Apply selection style first (cursor may override part of it)
// 	if r.editor.selStartRow >= 0 && r.editor.selEndRow >= 0 {
// 		r.grid.SetStyleRange(
// 			r.editor.selStartRow, r.editor.selStartCol,
// 			r.editor.selEndRow, r.editor.selEndCol,
// 			selectionStyleInstance,
// 		)
// 	}

// 	// Apply cursor style (if cursor is within selection, selection style will be overridden)
// 	if r.editor.cursorRow >= 0 && r.editor.cursorCol >= 0 {
// 		r.grid.SetStyle(r.editor.cursorRow, r.editor.cursorCol, cursorStyleInstance)
// 	}

// 	// Force the grid to redraw
// 	r.grid.Refresh()
// }

// func (r *noteEditorRenderer) Objects() []fyne.CanvasObject {
// 	return []fyne.CanvasObject{r.grid}
// }

// func (r *noteEditorRenderer) Destroy() {}
