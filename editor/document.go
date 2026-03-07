package editor

import (
	"strings"
)

// Document manages a slice of styled segments.
type Document struct {
	segments []TextSegment
}

// NewDocument creates an empty document.
func NewDocument() *Document {
	return &Document{
		segments: []TextSegment{{Text: "", Style: TextStyle{}}},
	}
}

// Text returns the plain text (for debugging/export).
func (d *Document) Text() string {
	var b strings.Builder
	for _, seg := range d.segments {
		b.WriteString(seg.Text)
	}
	return b.String()
}

// Lines splits the document into lines (by '\n') and returns each line's plain text.
func (d *Document) Lines() []string {
	return strings.Split(d.Text(), "\n")
}

// InsertText inserts text at the given character position with the specified style.
// Returns the new end position.
func (d *Document) InsertText(pos int, text string, style TextStyle) int {
	if pos < 0 {
		pos = 0
	}
	// Find the segment and offset where insertion happens.
	segIdx, offset := d.findSegmentAndOffset(pos)
	seg := d.segments[segIdx]

	// Split the segment if insertion point is inside it.
	if offset < len(seg.Text) {
		// Split into before and after
		before := TextSegment{Text: seg.Text[:offset], Style: seg.Style}
		after := TextSegment{Text: seg.Text[offset:], Style: seg.Style}
		// Replace the original segment with before and after
		d.segments = append(d.segments[:segIdx], append([]TextSegment{before, after}, d.segments[segIdx+1:]...)...)
		segIdx++ // now points to 'after'? Actually after is at segIdx+1, we need to insert new text between.
		// Insert the new text as a separate segment between before and after.
		newSeg := TextSegment{Text: text, Style: style}
		d.segments = append(d.segments[:segIdx], append([]TextSegment{newSeg}, d.segments[segIdx:]...)...)
	} else {
		// Insertion at the end of the segment: just add a new segment after it.
		newSeg := TextSegment{Text: text, Style: style}
		d.segments = append(d.segments[:segIdx+1], append([]TextSegment{newSeg}, d.segments[segIdx+1:]...)...)
	}

	// Merge adjacent segments with identical style to keep the model tidy.
	d.mergeSegments()

	// Return the new cursor position (original pos + len(text)).
	return pos + len(text)
}

// DeleteRange removes characters from start to end (exclusive).
func (d *Document) DeleteRange(start, end int) {
	if start >= end || start < 0 || end > d.len() {
		return
	}
	// Remove characters by adjusting segments.
	// This is more complex; for a prototype we'll rebuild from plain text.
	// But that would lose styles. For now, we'll implement a simple but style‑aware deletion.

	// We'll create a new slice of segments by walking through the old ones and cutting out the range.
	var newSegs []TextSegment
	currentPos := 0
	for _, seg := range d.segments {
		segStart := currentPos
		segEnd := currentPos + len(seg.Text)
		if segEnd <= start || segStart >= end {
			// Segment completely outside deletion range
			newSegs = append(newSegs, seg)
		} else if segStart >= start && segEnd <= end {
			// Segment fully inside deletion range: skip it
			// (do nothing)
		} else if segStart < start && segEnd > end {
			// Deletion range is inside this segment: split it
			before := TextSegment{Text: seg.Text[:start-segStart], Style: seg.Style}
			after := TextSegment{Text: seg.Text[end-segStart:], Style: seg.Style}
			if before.Text != "" {
				newSegs = append(newSegs, before)
			}
			if after.Text != "" {
				newSegs = append(newSegs, after)
			}
		} else if segStart < start && segEnd <= end {
			// Overlap at the end of the segment: keep beginning
			before := TextSegment{Text: seg.Text[:start-segStart], Style: seg.Style}
			if before.Text != "" {
				newSegs = append(newSegs, before)
			}
		} else if segStart >= start && segStart < end && segEnd > end {
			// Overlap at the start of the segment: keep tail
			after := TextSegment{Text: seg.Text[end-segStart:], Style: seg.Style}
			if after.Text != "" {
				newSegs = append(newSegs, after)
			}
		}
		currentPos = segEnd
	}
	d.segments = newSegs
	d.mergeSegments()
}

// ApplyStyle modifies the style of a character range.
func (d *Document) ApplyStyle(start, end int, fn func(*TextStyle)) {
	// TODO: Implement similarly to Insert/Delete but adjusting style.
	// For prototype, we'll skip detailed implementation; you can add later.
}

// Helper: find segment index and offset within that segment for a given character position.
func (d *Document) findSegmentAndOffset(pos int) (int, int) {
	if pos < 0 {
		return 0, 0
	}
	current := 0
	for i, seg := range d.segments {
		segLen := len(seg.Text)
		if pos < current+segLen {
			return i, pos - current
		}
		current += segLen
	}
	// Position beyond the end: return last segment and its length.
	lastIdx := len(d.segments) - 1
	return lastIdx, len(d.segments[lastIdx].Text)
}

// Helper: total character count.
func (d *Document) len() int {
	total := 0
	for _, seg := range d.segments {
		total += len(seg.Text)
	}
	return total
}

// Helper: merge adjacent segments with identical style.
func (d *Document) mergeSegments() {
	if len(d.segments) == 0 {
		return
	}
	merged := []TextSegment{d.segments[0]}
	for _, seg := range d.segments[1:] {
		last := &merged[len(merged)-1]
		if last.Style == seg.Style {
			last.Text += seg.Text
		} else {
			merged = append(merged, seg)
		}
	}
	d.segments = merged
}
