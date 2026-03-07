package editor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (e *NoteEditor) showContextMenu(pos fyne.Position) {
	// Build menu items based on selection and clipboard state
	items := []*fyne.MenuItem{}

	if e.selStart >= 0 && e.selEnd > e.selStart {
		// Have selection -> enable Copy and Cut
		items = append(items, fyne.NewMenuItem("Copy", func() {
			e.copyToClipboard()
		}))
		items = append(items, fyne.NewMenuItem("Cut", func() {
			e.cutToClipboard()
		}))
	}

	// Paste is always available if clipboard has text
	clipboard := fyne.CurrentApp().Clipboard()
	if clipboard.Content() != "" {
		items = append(items, fyne.NewMenuItem("Paste", func() {
			e.pasteFromClipboard()
		}))
	}

	items = append(items, fyne.NewMenuItem("Select All", func() {
		e.selectAll()
	}))

	if len(items) == 0 {
		return // nothing to show
	}

	menu := fyne.NewMenu("", items...)
	popup := widget.NewPopUpMenu(menu, fyne.CurrentApp().Driver().CanvasForObject(e))
	popup.ShowAtPosition(pos)
}

func (e *NoteEditor) copyToClipboard() {
	if e.selStart < 0 || e.selEnd <= e.selStart {
		return
	}
	text := e.doc.Text()[e.selStart:e.selEnd]
	clipboard := fyne.CurrentApp().Clipboard()
	clipboard.SetContent(text)
}

func (e *NoteEditor) cutToClipboard() {
	e.copyToClipboard()
	e.doc.DeleteRange(e.selStart, e.selEnd)
	e.cursor = e.selStart
	e.selStart, e.selEnd = -1, -1
	e.Refresh()
}

func (e *NoteEditor) pasteFromClipboard() {
	clipboard := fyne.CurrentApp().Clipboard()
	text := clipboard.Content()
	if text == "" {
		return
	}

	// If there's a selection, delete it first
	if e.selStart >= 0 && e.selEnd > e.selStart {
		e.doc.DeleteRange(e.selStart, e.selEnd)
		e.cursor = e.selStart
		e.selStart, e.selEnd = -1, -1
	}

	// Insert at cursor with default style
	style := TextStyle{} // You'll need to track current style (e.g., from toolbar)
	newCursor := e.doc.InsertText(e.cursor, text, style)
	e.cursor = newCursor
	e.Refresh()
}

func (e *NoteEditor) selectAll() {
	totalLen := e.doc.len()
	if totalLen > 0 {
		e.selStart, e.selEnd = 0, totalLen
		e.cursor = totalLen
		e.Refresh()
	}
}
