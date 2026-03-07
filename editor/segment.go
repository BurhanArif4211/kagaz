package editor

// TextStyle holds formatting attributes for a span of text.
type TextStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
	// Add more later (Heading, Bullet, etc.)
}

// TextSegment represents a contiguous run of text with a single style.
type TextSegment struct {
	Text  string
	Style TextStyle
}
