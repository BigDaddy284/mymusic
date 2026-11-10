# How MyTube works
<!-- docs: group=inside order=10 -->

MyTube is a native desktop player for YouTube Music and the music on your computer, written in Go. Here is how its engine, streams, services and updates fit together.

This page is the big picture for people who use MyTube and want to know
what happens behind a click. The pages after it go deeper: the
[audio engine](audio-engine.md), [where the audio comes from](streams.md)
and [your data and privacy](data-and-privacy.md). For the code itself, read
[the architecture](architecture/README.md).

## One program, drawn by itself

MyTube is one Go program. [MyGo](https://github.com/egoist/mygo), the
desktop framework it is built with, draws every window in Go, on the GPU:
the pages, the player bar, the queue, the lyrics, the mini player and the
command palette. There is no browser engine inside MyTube drawing its interface,
as there is in an Electron app, and no web page behind the windows.

A web view of the system, WebView2, opens only for two things, and only while it's needed:

- **Google's sign-in page**, when you sign in to YouTube Music.
- **The picture of music videos**, when "Show the video of music videos"
  is on in Settings › Playback. It closes a minute after the last picture.

MyTube builds without a C compiler: it calls the system's libraries
while it runs, for the audio output, the media controls, the tray and the
credential store, partly through
[purego](https://github.com/ebitengine/purego).

## From a click to your speakers

When you play a YouTube Music track, this happens:

1. **The player takes the track.** One part of MyTube owns the queue and
   what plays; the network calls run beside it, so the window never waits
   on them.
2. **MyTube finds the track's audio.** It asks YouTube's player API
   itself, as one of YouTube's own apps does, signed out. For the few
   tracks YouTube won't give it that way, such as those made for
   children, it runs [yt-dlp](https://github.com/yt-dlp/yt-dlp).
   [Where the audio comes from](streams.md) tells the details.
3. **MyTube's stream proxy fetches the audio.** It serves a download or
   the offline cache first, and otherwise fetches the stream
   from YouTube in chunks of 4 MB, in the background, so that seeking
   back costs nothing.
4. **The audio engine plays it.** MyTube decodes the audio itself, in
   Go, converts it to 48 kHz, turns a loud track down to an even
   loudness, applies the speed, the equalizer, volume levelling and the
   volume, and plays it on your system's output.
   [The audio engine](audio-engine.md) tells the details.
5. **The next track gets ready.** A second into a track, MyTube finds
   the next one's stream; ten seconds in, it fetches its audio, so that
   it starts at once when its turn comes.

A file on this computer skips the first three steps: the engine reads it
from its folder. [Music on this computer](local-music.md) tells more.

## YouTube Music

MyTube reads YouTube Music through the API YouTube Music's own website
uses. It needs no account and no API key:

- **Signed out**, it reads search and its suggestions, the home and
  explore pages, the charts, moods and genres, albums, artists,
  playlists, radios, podcasts and plain lyrics.
- **Signed in**, with the cookies of a Google session, it signs its
  requests as the website does, and reads and changes your library:
  liked songs, playlists, albums, subscriptions, likes, dislikes and
  saved shows. It can act as one of your brand channels.

Audio is a separate matter: MyTube always finds streams signed out, even
while you're signed in, as YouTube gives an account's requests fewer
streams.

## The services it talks to

MyTube has no server of its own. Apart from YouTube, it
talks to:

- **Lyrics:** [Unison](https://unison.boidu.dev), an open, crowdsourced
  database of lyrics, first, then [LRCLIB](https://lrclib.net), then
  YouTube Music's own lyrics. Unison also translates lyrics, through
  Google, and romanizes Arabic, but only while you turn those on.
- **Latin letters for lyrics** are made on your computer for every other
  script. Japanese romaji reads kanji with a dictionary MyTube downloads
  from GitHub the first time you turn romaji on (11.7 MB); the readings of
  Chinese are built in.
- **Last.fm**, once you connect it: what you play, to scrobble.
- **Discord**, through the Discord app on your computer, over its local
  pipe or socket: what plays, for your profile.
- **Cast devices** on your network, such as Chromecasts and Nest
  speakers, which MyTube finds and speaks to itself, in Go.
- **GitHub**, for new versions of MyTube, for yt-dlp, and for the kanji
  dictionary.

To tell whether the network answers, MyTube asks YouTube's addresses
over HTTPS. [Your data and privacy](data-and-privacy.md) lists
what each of these gets.

## Themes and plugins

- **Themes** are a `theme.json` file of colours, fonts and corners.
  MyTube draws every theme in the layout you choose in Settings ›
  Appearance: Blend (the default), Expressive, Crate or Tiles. It can
  also take its colours from a file that pywal or matugen write, and
  follows the file as it changes. [Themes](themes.md) lists every token.
- **Plugins** are JavaScript modules, usually compiled from TypeScript.
  Each plugin you enable runs in a JavaScript engine of its own,
  [Sobek](https://github.com/grafana/sobek), written in Go, with
  `console`, timers, `fetch` and MyTube's plugin API. A plugin's pages
  and widgets are elements that MyTube draws in your theme, and its
  canvases are drawing calls that MyTube paints. The permissions a
  plugin lists are shown to you, not enforced, so MyTube asks you to
  trust a plugin before it enables it. [Plugins](plugins.md) documents the
  API. MyTube loads the add-ons it ships with and those you put in your
  add-ons folders; it doesn't download any.

## Updates

### MyTube itself

MyTube doesn't look for new versions of itself, or install them: get a
new one from the releases page when you want it.

### yt-dlp and the JavaScript runtime

MyTube needs yt-dlp for some tracks, and yt-dlp needs a JavaScript
runtime to solve the challenges YouTube sends. Settings › Playback ›
Tools shows both.

- **MyTube's own yt-dlp** is the one you install from Settings ›
  Playback or Home. MyTube downloads it from yt-dlp's releases on
  GitHub, checks it against the checksums the release publishes, and
  keeps it in the `bin` folder of its data folder. A minute after each
  start, and then once a day, MyTube replaces it when
  a newer release is out.
- **A yt-dlp of your own**, from a package manager, pip or winget, or the one you name under Own yt-dlp, is never changed.
  MyTube asks GitHub for the latest release at most every 6 hours, and
  when yours is older, Home and Settings › Playback say so, with the
  command that updates it when its folder tells how you installed it.
  They also offer MyTube's own instead.
- **The JavaScript runtime** is yours: MyTube installs none. It looks
  for Deno, Node.js, Bun and QuickJS, in that order, in its `bin` folder,
  on your `PATH`, and where Deno, Bun, Volta, Scoop and WinGet install
  them, or takes the one you choose under JavaScript runtime. You update
  it the way you installed it.
