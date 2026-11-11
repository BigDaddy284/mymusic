# Capture mode

`MYTUBE_CAPTURE=<dir>` screenshots pages of the main window, and
`MYTUBE_CAPTURE_QUIT=1` quits once they are saved:

```sh
MYTUBE_CAPTURE=shots MYTUBE_CAPTURE_QUIT=1 MYTUBE_CAPTURE_ROUTES="/|/search?q=jazz|/settings" MYTUBE_CAPTURE_SIZE=1440x900 ./MyTube
```

`MYTUBE_CAPTURE_WAIT` sets how long each step gets (`5s`). Steps that
aren't pages act: `act:shot:<name>` saves the window as it is,
`act:mute` mutes the player, `act:play:<search>` plays what a search
finds, `act:local:<search>` the music on this computer it finds, all of
it without one, `act:resume` plays the queue kept in the data folder,
`act:wait:<duration>` waits, `act:palette` opens the command palette,
`act:search:<text>` types text in the search box, whose suggestions
show, `act:suggested:<song|video|album|artist|playlist>` chooses the
first suggestion of that kind, as a click would,
`act:refresh` refreshes the page, as F5 does,
`act:panel:<queue|lyrics|none>` shows a side panel, `act:scroll:<y|end>`
scrolls the page, `act:theme:<id>` wears a theme,
`act:mini` opens the mini player, which saves
`mini.png`, `act:library:<add|remove>:<search>` saves the first YouTube
Music song a search finds to the library there, or takes it out, as its
menu does, `act:title` logs the window's title, which the screenshots
leave out, `act:nowplaying` the cover MyTube hands the system's media controls,
and `act:memory` and `act:heap` log the Go heap and save a
profile of it. Nothing plays unless the player is muted,
or `MYTUBE_SILENT=1` has the audio engine play into no device, which
keeps the visualizer and the lyrics moving without a sound.
The
site's screenshots are made this way, from a throwaway profile, at
1536x900 (`MYTUBE_CAPTURE_SIZE=1552x909`: the size includes the
window's invisible borders). `go test . -run Render` with
`MYTUBE_RENDER=<dir>` draws the pages and themes without a window.
