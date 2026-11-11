# The services

- **`internal/media`** defines what every service shares: tracks, albums,
  artists, playlists, shelves. IDs carry their service: `yt:<video id>`,
  `sp:<track id>`, `mx:<mix id>`. The services give a playlist with its
  first page of tracks, and `More`, the token of the next;
  `Playlist.ReadOn` reads on a page at a time, asks again after 1, 2 and
  4 seconds for a page that fails, and stops at `MaxTracks`, 20,000, past
  which the playlist is `Cut`. `Catalog.Playlist` reads one whole, for
  what takes all of it, as a download.
- **`internal/youtube`** reads YouTube Music through the API its website uses:
  search, suggestions (searches, then songs, albums, artists and
  playlists), the home and explore pages, charts, moods and genres, albums,
  artists, playlists, related tracks, radios and plain lyrics. Below a
  playlist the user owns, YouTube Music suggests songs to add, five at a
  time: the playlist's `Suggest` is the token of the first, and each
  answer has the token of a refresh, for the next (`PlaylistSuggestions`).
  Signed in with the cookies of the user's session, which it signs its
  requests with as the website does, it reads and changes their library:
  liked songs, playlists, albums, subscriptions, likes, dislikes and
  saves, songs saved to the library apart from likes, by the feedback
  tokens of their menus in the queue of their "next", and it deletes the playlists the user made, which it tells from
  those they saved by the menus YouTube sends with them. An add of a
  video the playlist holds already YouTube refuses, offering in its
  answer the same add with its check skipped (`DEDUPE_OPTION_SKIP`):
  that, not the answer's words, which are in the account's language,
  tells it from other refusals (`ErrAlreadyInPlaylist`), and
  `AddToPlaylistAnyway` sends it. It lists the
  Google accounts of the session and their brand channels, from YouTube's
  website, and uses a brand channel by its page ID: one whose page ID the
  list doesn't tell is listed as `Unswitchable`, as without it the account
  itself would be used. Once the app switches to a brand channel, it goes
  back if YouTube Music answers as the Google account itself, as the list
  names it; any other answer keeps the switch. It also reads
  whether YouTube labels a video as made with AI, from the answer of
  YouTube's own website, asked signed out with a field mask
  (`MadeWithAI`): nothing of YouTube Music carries the label. Its parsers
  are tested against recorded responses in `testdata/`, and against the
  live service with `MYTUBE_LIVE=1 go test ./internal/youtube -run Live`
  (`MYTUBE_YT_COOKIES` adds a signed-in account, and
  `MYTUBE_YT_WRITE=1` lets it make and delete a playlist there, and
  dislike a song and take the dislike back).
- **`internal/sysproxy`** picks the proxy of every request MyTube makes,
  as the `Proxy` of `http.DefaultTransport`, which the clients use or
  clone: the environment's (`HTTPS_PROXY` and the like) when it names
  one; on Windows, failing that, the one of Windows' settings, which
  browsers, WebView2 and yt-dlp use too. Automatic setup comes first,
  detection or a setup script, which WinHTTP answers for each site and
  which is remembered for 5 minutes, failures too; then the proxy
  server, unless the bypass list names the address. Loopback addresses
  always connect directly. It reads the settings again when their
  registry key changes, as VPN and proxy apps turn them on and off,
  logs what MyTube connects through, without user names or passwords,
  and gives yt-dlp a `--proxy` where yt-dlp wouldn't find that proxy
  itself. A package that clones `http.DefaultTransport` as it
  initializes imports it, as `internal/online` does.
- **`internal/match`** finds the song of a music video, by its title and artists, which
  the video's radio starts from: YouTube Music makes the radio of a video
  out of videos.
- **`lyrics`** finds synced lyrics on Unison and LRCLIB, reads TTML and LRC,
  and translates lyrics on Unison.
- **`internal/lastfm`** signs in to Last.fm with the user's API account and
  scrobbles by Last.fm's rules, keeping what fails offline for later.
