package settings

// Desktop is how MyTube fits in the desktop: Settings › Desktop.
type Desktop struct {
	// CloseToTray keeps the music playing as the window closes, from the
	// tray icon. It is off at first.
	CloseToTray bool `json:"closeToTray"`
	// Notifications tell what starts playing while MyTube is in the
	// background.
	Notifications bool `json:"notifications"`
	// ShowInTitle names what plays in the main window's title.
	ShowInTitle bool `json:"showInTitle"`
	// KeepOnTop keeps the main window above the other windows (#350). It
	// is off at first; the mini player stays on top either way.
	KeepOnTop bool `json:"keepOnTop"`
	// MiniPanel is the panel the mini player shows beside it as it opens
	// (#371): MiniQueue, MiniLyrics or MiniOutput, or "" for none.
	MiniPanel string `json:"miniPanel"`
	// MediaKeys has the keyboard's play, next and previous keys control
	// MyTube.
	MediaKeys bool `json:"mediaKeys"`
}

// defaultDesktop is how a new install fits in the desktop.
func defaultDesktop() Desktop {
	var d Desktop
	d.Notifications = true
	d.ShowInTitle = true
	d.MediaKeys = true
	return d
}

// clean is d with a panel of the mini player it has.
func (d Desktop) clean() Desktop {
	d.MiniPanel = cleanMiniPanel(d.MiniPanel)
	return d
}
