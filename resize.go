package main

// The sidebar and the side panel are as wide as the user drags them: a
// thin line at their inner edge takes the pointer, and the width it ends at
// stays in the settings.

import "github.com/egoist/mygo/ui"

// The least and the most the sidebar and the side panel may be dragged to.
const (
	railLeast, railMost   = 180, 420
	panelLeast, panelMost = 260, 620
	// resizerWidth is the room the line that resizes takes the pointer in.
	resizerWidth = 7
)

// clampWidth keeps w within least and most.
func clampWidth(w, least, most float32) float32 { return min(max(w, least), most) }

// railResizer is the line at the sidebar's right edge, which drags its
// width. A folded sidebar keeps its icons' width.
func (s *shell) railResizer(c *ui.Context) {
	if s.railFolded() {
		return
	}
	s.resizer(c, "rail-resizer", "Resize the sidebar", s.railWidth(), 1, railLeast, railMost, func(w float32) {
		s.set.SidebarWidth = float64(w)
	})
}

// panelResizer is the line at the side panel's left edge, which drags its
// width.
func (s *shell) panelResizer(c *ui.Context) {
	s.resizer(c, "panel-resizer", "Resize the side panel", s.panelWidth(c), -1, panelLeast, panelMost, func(w float32) {
		s.set.PanelWidth = float64(w)
	})
}

// resizer is a vertical line that drags a width, now width, by the
// pointer's moves times sign, within least and most; set takes each new
// width, and the settings keep the last once the drag ends. A double
// click puts the layout's width back.
func (s *shell) resizer(c *ui.Context, key, label string, width, sign, least, most float32, set func(w float32)) {
	k := &s.lk
	h := ui.Box(c.Key(key)).Width(resizerWidth).Shrink(0).AlignSelf(ui.Stretch).Cursor(ui.CursorResizeColumn).Label(label).Tooltip(label)
	if dx, _, ok := h.Dragged(); ok && dx != 0 {
		set(clampWidth(width+sign*dx, least, most))
	}
	if _, _, ended := h.DragEnded(); ended {
		s.save()
	}
	if h.DoubleClicked() {
		set(0)
		s.save()
	}
	ink := k.border
	if h.Hovered() || h.Dragging() {
		ink = k.accent
	}
	h.Children(func() {
		ui.Box(c).Width(2).AlignSelf(ui.Stretch).MarginX((resizerWidth - 2) / 2).Background(ink)
	})
}
