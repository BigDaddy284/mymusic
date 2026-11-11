<p align="center">
  <img src="brand/icon.png" width="112" height="112" alt="MyTube's icon" />
</p>

<h1 align="center">MyTube</h1>

<p align="center"><strong>YouTube Music, YouTube videos and your own music, in one Windows app.</strong></p>

<p align="center">
  <a href="https://github.com/BigDaddy284/mytube/actions/workflows/ci.yml"><img src="https://github.com/BigDaddy284/mytube/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
  <img src="https://img.shields.io/badge/platform-Windows%2010%20%7C%2011-lightgrey" alt="Windows 10 and 11" />
  <a href="#license"><img src="https://img.shields.io/badge/license-MIT%20or%20GPL--3.0-blue" alt="MIT or GPL-3.0" /></a>
</p>

MyTube has two sides, and the mark at the top of the sidebar flips
between them:

- **Music** plays YouTube Music and the music files on your computer from
  one search, one queue and one library.
- **Video** is YouTube: your home feed, subscriptions, history, playlists,
  Watch later and liked videos, with search, Trending and a Random button.

It's one Go program that draws its own windows. It has no account, server,
analytics or telemetry of its own, and YouTube works without signing in.

## Music

- **One search** across YouTube Music and your own files, each result
  marked with where it comes from.
- **One queue and mixed playlists.** A file on your computer can play right
  after a YouTube Music track, and a playlist can hold both. A radio
  carries on when the queue runs out.
- **Explore:** new releases, the charts, moods and genres.
- **Synced lyrics,** with a translation under each line if you want one.
- **Sound:** crossfade, loudness levelling, a ten-band equalizer with
  presets, a visualizer and an offline cache.
- **Podcasts,** a mini player that stays on top, casting, and listening
  stats.
- **Your YouTube Music account,** optional: likes, playlists, albums and
  followed artists stay in step, and what you play joins your history.
- **Last.fm** scrobbling, kept while you're offline.

## Video

- **Home** has YouTube's feed and, as chips beside it, your
  Subscriptions, History, Playlists, Watch later and Liked videos.
- **Keep watching while you browse.** The dock button on the player bar
  shrinks the video to the corner and takes you back to the page you were
  on. Full App fills the window, and pop-out puts the video in its own
  window over your other apps.
- **Video tabs** along the top bar keep videos to come back to, from where
  you left them.
- **SponsorBlock:** skip or mute sponsor reads and other marked parts.
- **Block lists:** keep channels and words off every page. Shorts are
  hidden unless you turn them back on.
- **Card menus** like YouTube's: Up next, Watch later (and remove from it),
  save to a playlist, not interested, block the channel.
- **Previews** play muted as you hover a card, with a red bar of how much
  you've watched, on any device.
- **Chapters, channel pages, likes, casting to a TV,** a separate output
  device for video sound, and **play as audio** to send a video to the
  music side with the EQ.
- **Stats:** how long you watched, and your top videos and channels.

## Make it yours

- **Five layouts:** Blend (the default), Studio, Expressive, Crate and
  Tiles.
- **Sixteen themes,** with an accent colour of your own or one taken from
  the artwork. You can add your own themes ([docs/themes.md](docs/themes.md)).
- **Keyboard first:** shortcuts for everything common, global media keys,
  your own key for any command, and a command palette
  (<kbd>Ctrl</kbd> <kbd>K</kbd>).
- **A settings file:** `settings.json` in `%AppData%\MyTube`. Save it
  while MyTube runs and it applies at once.

## Install

MyTube runs on Windows 10 and 11. There are no published releases yet, so
get a build from CI:

1. Open the [Actions tab](https://github.com/BigDaddy284/mytube/actions/workflows/ci.yml)
   and pick a green run (for a pull request, its run on the PR's branch),
   or press **Run workflow** to build any branch.
2. Scroll to **Artifacts** at the bottom of the run's page and download
   **windows-amd64**. You must be signed in to GitHub, and artifacts are
   kept for 7 days.
3. Unzip it. It holds the installer (`*-setup-*.exe`) and a portable zip;
   for the portable one, unpack it anywhere and open `MyTube.exe`.

The builds aren't signed yet, so:

- SmartScreen may warn: choose **More info**, then **Run anyway**.
- Windows Defender may quarantine `MyTube.exe`. Restore it from Windows
  Security › Protection history and choose **Allow on device**, or add the
  MyTube folder under Exclusions.

MyTube doesn't update itself. Signing in to YouTube needs the WebView2
Runtime, which Windows 11 includes; on Windows 10, get it from
[Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/).
For the few tracks MyTube can't play itself, it installs
[yt-dlp](https://github.com/yt-dlp/yt-dlp) for you from Settings ›
Playback.

Something not working? See [docs/troubleshooting.md](docs/troubleshooting.md).

## Build from source

You need [Go](https://go.dev/dl/) 1.27.3 or later, [Bun](https://bun.sh)
1.2 or later and Git.

```powershell
winget install GoLang.Go
winget install Oven-sh.Bun
winget install Git.Git
```

Then, in a new terminal:

```powershell
git clone https://github.com/BigDaddy284/mytube
cd mytube
bun install
bun run dev      # build, start, and rebuild as you edit
```

```sh
go test ./...    # the tests
bun run build    # the installer and the app, in build/
```

`bun run dev` keeps its data apart from an installed MyTube, as
`MyTube Dev`. With [Nix](https://nixos.org), `nix develop` gives you the
exact tools CI builds with. [docs/architecture](docs/architecture/README.md)
explains how the parts fit together, and [CONTRIBUTING.md](CONTRIBUTING.md)
lists the checks to run before a pull request.

## Privacy

MyTube talks only to the services it needs: YouTube and Google (the
catalog, streams, and your account if you sign in), Unison and LRCLIB for
lyrics, SponsorBlock (by part of a hash of the video, so it doesn't learn
what you watch), GitHub for yt-dlp, and Last.fm if you connect it. Your
settings, library, cache and log stay on your computer, and your sign-ins
are kept in Windows' Credential Manager.
[docs/data-and-privacy.md](docs/data-and-privacy.md) has the details.

## Credits

MyTube is built with [MyGo](https://github.com/egoist/mygo) and
[Go](https://go.dev), and started as a fork of
[MixMusic](https://github.com/fzn0x/mixmusic) by fzn0x and its
contributors, under the MIT license. [docs/credits.md](docs/credits.md)
lists every project it uses.

## License

MyTube is dual-licensed under the [MIT license](LICENSE-MIT) or the
[GNU GPL, version 3](LICENSE-GPL), whichever you prefer. The fonts in
`fonts/` are under the SIL Open Font License, and the Opus decoder, the
AAC decoder and the HCT colour code keep their own licenses, next to their
code.

MyTube is not affiliated with Google. YouTube and YouTube Music are
trademarks of Google.
