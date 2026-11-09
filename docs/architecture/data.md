# Data

Everything lives in the user's data folder (`%AppData%\MyMusic`): `settings.json`, `library.json`,
`mix-covers/` (the covers the user chose for mixes),
`queue.json`, `positions.json` (where the user left podcasts' episodes),
`local-files.json` (the index of the music on the computer),
`addons/`, `bin/` (yt-dlp), the YouTube session
(`youtube-session.bin`, encrypted with a key of the system's credential
store, which also keeps the Spotify token and the Last.fm connection), the
scrobbles still to send, and `updater.json`, the update preferences. A
secret the Credential Manager can't keep, as in a logon session without
credentials, is sealed for the user with DPAPI in `spotify-token.dpapi`,
`lastfm.dpapi` or `youtube-session-key.dpapi`. `secretStore` (`internal/app/secrets.go`) moves a
secret found in a file into the credential store once it can keep it.
`settings.json` may be edited while MyMusic runs: `App.watchSettings`
watches it (`mygo.WatchFiles`, or a look every two seconds where it can't),
and `App.ReloadSettings` applies an edit as Settings applies a change, then
tells the windows (`App.Watch`). A file with a mistake is told of after a
second, and kept as `settings.json.broken` if MyMusic quits before it is
fixed; a `settings.json` that is a link stays one. The
cache folder holds audio, the covers shown (`artwork/`), the covers of
the music on the computer (`local-art/`), matches, how loud
YouTube measured the tracks played (`loudness.json`), the videos YouTube
labels as made with AI, for 30 days (`ai-labels.json`), the releases of
the artists followed (`releases.json`), the dictionary of
the readings of kanji romaji downloads (`romaji/`), and yt-dlp's cache,
and the logs folder the log and the crash reports. The downloads are in a
folder of their own, which the user can move in Settings › Playback
(`downloadsFolder` in `settings.json`): by default
`%LocalAppData%\MyMusic\downloads`: never a folder the system treats as a
cache, whose files cleaning tools delete, nor the Windows folder that
roams. `MYMUSIC_DATA_DIR` moves all four into one folder,
for a portable copy or a test run.
`MYMUSIC_YTDLP_ONLY=1` has yt-dlp find every stream, as before MyMusic
asked YouTube's player API, to time one against the other or tell which
fails: the log tells which streams took over a second to find, and which
tracks over two to start.
