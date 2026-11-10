package app

// Whether a window of MyTube's is in front, which shows what a
// notification would tell: the notification would only repeat it over the
// window the user looks at (#289).

import "github.com/egoist/mygo"

// focused reports whether a window has the keyboard; tests, which open no
// windows, replace it.
var focused = (*mygo.Window).IsFocused

// inFront reports whether the main window or the mini player has the
// keyboard: both show what plays, and whether it is liked or disliked.
func (a *App) inFront() bool {
	a.mu.Lock()
	main, mini := a.main, a.mini
	a.mu.Unlock()
	return main != nil && focused(main) || mini != nil && focused(mini)
}
