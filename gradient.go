package main

import "github.com/egoist/mygo/ui"

// The MyTube themes' gradient, after Stills & Reels' Theme.cs and
// UiControls.cs: its primary buttons are filled with the accent gradient,
// blue to purple from left to right, under white ink that hovering lightens
// with 28 of 255 of white and pressing with 50; its checked buttons with
// the gradient at 70 of 255 in dark mode and 40 in light mode, under the
// text's ink.

// washAlpha is the opacity of the gradient on what is chosen.
func washAlpha(dark bool) float32 {
	if dark {
		return 70.0 / 255
	}
	return 40.0 / 255
}

// fillState is how far a filled control is pressed.
type fillState int

const (
	fillRest fillState = iota
	fillHover
	fillPress
)

// lift is the white Stills & Reels lays over a primary button in state.
func (f fillState) lift() float32 {
	switch f {
	case fillHover:
		return 28.0 / 255
	case fillPress:
		return 50.0 / 255
	}
	return 0
}

// stateOf is the state of the button b.
func stateOf(b ui.Element) fillState {
	switch {
	case b.Pressed():
		return fillPress
	case b.Hovered():
		return fillHover
	}
	return fillRest
}

// fillAccent fills e as the accent fills the main buttons, in state: with
// the theme's gradient where it has one, or with the accent, which hovering
// mixes towards the text and pressing towards the ground. It returns the
// ink that goes on it.
func (k *look) fillAccent(e ui.Element, state fillState) ui.Color {
	if k.gradient {
		w := ui.RGB(255, 255, 255)
		from, to := k.gradFrom.Mix(w, state.lift()), k.gradTo.Mix(w, state.lift())
		e.LinearGradient(ui.LinearGradient{From: from, To: to, Angle: 90})
		return w
	}
	switch state {
	case fillHover:
		e.Background(k.accent.Mix(k.text, 0.14))
	case fillPress:
		e.Background(k.accent.Mix(k.bg, 0.14))
	default:
		e.Background(k.accent)
	}
	return k.accentText
}

// fillSelected fills e as what is chosen: with a wash of the theme's
// gradient where it has one, or with the selected colour.
func (k *look) fillSelected(e ui.Element) {
	if k.gradient {
		a := washAlpha(k.dark)
		e.LinearGradient(ui.LinearGradient{From: k.gradFrom.Alpha(a), To: k.gradTo.Alpha(a), Angle: 90})
		return
	}
	e.Background(k.selected)
}
