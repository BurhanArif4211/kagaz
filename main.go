package main

import (
	"github.com/burhanarif4211/kagaz/editor"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Kagaz")

	ed := editor.NewNoteEditor()
	// Optionally set initial content
	ed.SetContent([]editor.TextSegment{{Text: "Hello, world!", Style: editor.TextStyle{}}})

	w.SetContent(container.NewBorder(nil, nil, nil, nil, ed))
	w.ShowAndRun()
}
