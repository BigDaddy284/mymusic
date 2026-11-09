// Package brand holds MyMusic's brand files: the app's icon, which
// mygo.config.ts names, the tray's, which the app embeds, and the mark of
// its logo, the music note, for the code that draws it. make_icon.py draws
// them all.
// The app builds them in, so they stay out of resources/, whose every file
// `mygo build` ships beside the app (#303).
//
// discord/ holds the badges of the Discord status, which tools/discordbadges
// draws: Discord reads them from GitHub, and the app doesn't build them in.
package brand

import _ "embed"

// Mark is MyMusic's mark, the note in its gradient, without the dark tile
// of the app's icon.
//
//go:embed mark.svg
var Mark []byte
