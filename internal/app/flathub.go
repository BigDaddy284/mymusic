package app

// Flathub builds MyMusic from source with the manifest of
// packaging/flathub, whose build.sh leaves flathubMarker in the app: the
// release's own .flatpak, packed from the Linux archive, has none. Flathub
// updates its copies, through the desktop's software center, and they run
// the yt-dlp that comes with them, which Flathub updates too.

// flathubMarker is the file only Flathub's build of MyMusic has.
const flathubMarker = "/app/share/mymusic/flathub"

// flathubYTDLP is the yt-dlp of Flathub's build, a module of its manifest.
const flathubYTDLP = "/app/bin/yt-dlp"

// bundledYTDLP returns the yt-dlp that comes with this copy of MyMusic,
// which it runs instead of downloading one, or "" when it has none.
func bundledYTDLP() string {
	if installKind() == "flathub" {
		return flathubYTDLP
	}
	return ""
}
