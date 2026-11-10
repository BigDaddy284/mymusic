package app

// The watching of the Video side's videos, reported to the account's
// history on YouTube as its own player reports it: the embedded player
// plays signed out, so without this YouTube wouldn't know of it, nor show
// the videos in the history, in Continue watching, or with their red bars
// on other devices. It goes as the episodes' listening does
// (youtube_watchtime.go): on a schedule while the video plays, as it
// pauses and plays again, and a last time as another video plays or none
// does; the history setting governs it, and signed out nothing goes.

import "time"

// ReportWatching tells where the video id is, at seconds, playing or not,
// at rate, as the window's player tells it: "" once no video plays.
// loading is set while the player tells nothing yet.
func (v *Videos) ReportWatching(id string, at float64, playing bool, rate float64, loading bool) {
	a := v.a
	on := a.settings.Get().YouTubeHistory
	pos := time.Duration(at * float64(time.Second))
	if rate <= 0 {
		rate = 1
	}
	a.followWatchtime(&a.videoWatch, id, id, pos, playing, rate, loading, on, true)
}
