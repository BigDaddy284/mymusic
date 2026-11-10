# Webviews

Two web pages still show, and only while they're needed:

- **Google's sign-in to YouTube Music**, whose cookies MyMusic keeps.
- **The picture of music videos** (`video.go`, `internal/app/video.go`),
  with "Show the video of music videos" on in Settings › Playback: a web
  view of the main window, on the now playing page, in place of the
  cover, of `public/video.html`, the app's one page. It is made
  when a picture first shows, and closed a minute after the last one. The
  page plays the picture, muted, through Media Source Extensions, fetching
  it from the proxy at the stream scheme, which serves that page alone,
  and follows the engine's time (`Player.Clock`), which the window sends
  it four times a second with the moment it read it. A box over the web
  view keeps the pointer and the keyboard in the native interface.
