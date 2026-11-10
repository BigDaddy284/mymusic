# Getting started
<!-- docs: group=start order=10 -->

MyTube is a free desktop player for Spotify and YouTube Music on Windows. This page installs it and shows you around.

MyTube puts both services in one window: one search looks through YouTube Music and Spotify at once, one queue plays tracks of either, and your liked songs and mixes hold tracks of both. YouTube Music plays as soon as you open it, with no account. Spotify, and your own YouTube Music library, are there once you sign in. MyTube has no account, server, analytics or telemetry of its own, and it is open source, under the MIT license or the GPL 3.0.

## Install

Download the file for Windows 10 or 11, on x64 or Arm, from the
[latest release](https://github.com/bigdaddy284/mymusic/releases/latest): the installer, `mytube-<version>-windows-<arch>-setup.exe`, or the portable zip, `mytube-<version>-windows-<arch>-portable.zip`, which you unpack anywhere and start with `MyTube.exe`.

`<arch>` is `amd64` for x64 and `arm64` for Arm. Each release lists `SHA256SUMS` to check the downloads with, and files you don't need, which installed copies update from. The [README](../README.md#install) has the full table of downloads.

The installer and the portable zip update themselves: MyTube looks for a new version at start and every few hours, and asks before it installs one, unless you let it install updates by itself. [Settings › About](settings.md#about) shows your version and the newest one.

## The first launch

- **Windows.** The installer isn't signed yet, so SmartScreen may warn about it: choose More info, then Run anyway. Signing in to YouTube Music and playing Spotify in MyTube itself need the WebView2 Runtime, which Windows 11 includes. On Windows 10, install it from [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/).
- **yt-dlp.** MyTube finds the audio of most YouTube Music tracks itself. For the rest, such as tracks made for children, it runs [yt-dlp](https://github.com/yt-dlp/yt-dlp). Until yt-dlp is there, Home says "Some tracks need yt-dlp", with a button that installs it; Settings › Playback › Tools installs it too, and MyTube keeps the copy it installed up to date. yt-dlp needs a JavaScript runtime: [Deno](https://deno.com), Node.js, Bun or QuickJS, which MyTube finds once you install one. [Where the audio streams come from](streams.md#when-yt-dlp-is-needed) tells which tracks need it.

MyTube opens on Home in the [Blend layout](look-and-feel.md#layouts). Your settings and library stay in a [data folder](data-and-privacy.md#the-data-folder) of your own: `%AppData%\MyTube`.

## The window

![MyTube's window on Home, in its default look: the sidebar on the left, with its pages and Your mixes under them, the top bar with the back and forward arrows, Refresh and the search box, Home's Jump back in, Add Spotify, Liked songs and Your mixes tiles over its Recently played shelf, and the player bar along the foot](shots/home.webp)

### The sidebar

The sidebar on the left goes to MyTube's pages: Home, Search, Explore, Library, Stats and [Settings](settings.md), then the pages plugins add. Under them, Your mixes lists your liked songs, your downloads once you have some, and each [mix](library.md#mixes) you made, with how many tracks each holds; its + button makes a new mix. At its foot, a card shows what plays, and opens the now playing page. The button beside MyTube's name folds the sidebar to its icons, as Ctrl+B does; in a narrow window it folds by itself.

- **Home** has your own tiles first: Jump back in plays the track you played last, then Liked songs and Your mixes, and Add Spotify while you aren't signed in to it. Then come your shelves, such as Recently played, [New from artists you follow](library.md#new-releases), and those of YouTube Music and Spotify. The pivots at the top keep Everything, only Yours, only Spotify's or only YouTube Music's, and [Edit home](home.md#arrange-home) moves, hides and adds its shelves.
- **Explore** has YouTube Music's new albums, singles and trending songs, its charts, and its moods and genres, each with playlists of its own.
- **Library** has your liked songs, mixes, podcasts, downloads and history, and your library on each service you signed in to: see [Your library](library.md).
- **Stats** shows what you listened to, day by day: see [Stats](library.md#stats).

### The top bar

The top bar holds the back and forward arrows, and [Refresh this page](settings.md#refresh-a-page), which asks the services again for what the page shows (F5). [Pulling a page down](settings.md#pull-to-refresh) past its top refreshes it too. Then comes the search box, "Search Spotify and YouTube Music" (Ctrl+F), which suggests searches, songs, albums, artists and playlists as you type: see [Playing music](playing-music.md#search). Once you scroll down a playlist, an album, an artist or your liked songs, the top bar shows its cover and name, with Play, Shuffle and Download. The button at the right of the top bar, which shows its key, opens the [command palette](settings.md#the-command-palette).

### The player bar

The player bar runs along the foot of the window. In the Expressive layout it floats over it.

- **On the left,** what plays: its cover, which opens the now playing page (and its menu, on a right click), its title, which opens its album, and its artists. The line under the title says "via YouTube Music" for a Spotify track that plays from YouTube Music, "on Spotify" for one on Spotify Connect (see [How Spotify tracks play](playing-music.md#how-spotify-tracks-play)), the [speed](playing-music.md#speed) when it isn't 1×, and "waiting for the network" when the [network dropped](playing-music.md#without-the-network). In a wide window, the track that plays next shows under it. The heart [likes](library.md#liked-songs) the track and the thumbs-down [dislikes](playing-music.md#dislikes-and-skipping) it.
- **In the middle,** Shuffle, Previous, Play or Pause, Next and Repeat, and the seek bar with the time and the length.
- **On the right,** the 5.1 button for a track that has a 5.1 surround mix, Lyrics, Queue, the [cast button](casting.md) while a Chromecast, a Google Home or Nest speaker, or a TV with Chromecast built in is on your network, the Spotify Connect devices when Spotify tracks play on Spotify Connect, Speed, the [output device](playing-music.md#the-output-device) where the window has room for it, the volume, and the [mini player](desktop.md#the-mini-player).

### Up next, lyrics and now playing

- **Up next,** the queue, opens at the right of the window from the Queue button, or with Ctrl+J: the track playing, then what plays after it. See [The queue and Up next](playing-music.md#the-queue-and-up-next).
- **The lyrics panel** opens in the same place from the Lyrics button, or with Ctrl+Shift+L: see [Lyrics](lyrics.md). For a podcast episode, the same button shows [its details](podcasts.md#episode-details-and-chapters).
- **The now playing page** gives the track the whole window: open it with its cover in the player bar, the card in the sidebar, or Alt+4. It shows the cover with a visualizer over its foot (or the [picture of a music video](playing-music.md#the-picture-of-music-videos), if you turn that on), the artists and the album, the seek bar, play and next, like and dislike, the lyrics, the [widgets of plugins](addons.md#what-a-plugin-can-do), and at its foot the songs, playlists and artists YouTube Music relates to the track.

![The now playing page of アイドル by YOASOBI: its cover with the visualizer over its foot, the seek bar, play, next, like and dislike, and its synced lyrics under the Synced, Plain text, Translate and Romaji switches, with Up next open at the right of the window](shots/now-playing.webp)

### The command palette and shortcuts

The command palette, Ctrl+K, finds every command by its name, plugins' commands too, with its key beside it. Type anything and it also offers to search for it, and to search the settings for it. These keys work from the start; [Settings › Shortcuts](settings.md#change-a-shortcut) changes any of them, and gives keys to the commands that have none:

| Keys | What they do |
|---|---|
| Space | Play or pause |
| ← and → | Back or forward 10 seconds |
| ↑ and ↓ | Volume up or down |
| M | Mute |
| Ctrl+← and Ctrl+→ | Previous and next track |
| Ctrl+S, Ctrl+R | Shuffle, and Repeat |
| Ctrl+L, Ctrl+D | Like or dislike what plays |
| Ctrl+F | Search |
| Ctrl+J, Ctrl+Shift+L | Show or hide the queue, or the lyrics |
| Alt+1 to Alt+5 | Home, Library, Liked songs, Now playing, Explore |
| Ctrl+M | Mini player |
| Ctrl+B | Fold or unfold the sidebar |
| Ctrl+, | Settings |
| Ctrl+= and Ctrl+-, Ctrl+0 | Zoom in and out, back to 100% |
| F5 | Refresh this page |

The [media keys](desktop.md) of your keyboard and headphones work too, and Settings › Shortcuts can set [keys that control the music from any app](desktop.md#global-shortcuts).

## What works without an account

With no account at all, MyTube:

- searches YouTube Music and plays it, with its home shelves, Explore, artists, albums, playlists, [podcasts](podcasts.md) and [radios](playing-music.md#autoplay-and-radios);
- keeps your liked songs, dislikes, mixes, history and stats on your computer, in [your library](library.md);
- shows [lyrics](lyrics.md), [downloads tracks](library.md#downloads) to play without the network, and keeps the tracks you played in its [offline cache](library.md#the-offline-cache);
- shows what plays on your [Discord profile](desktop.md#discord), from the Discord app on your computer, with nothing to set up (Settings › Desktop turns it off).

Every account is optional, and you add them in Settings › Accounts: [Accounts](accounts.md) tells each step.

### What each account adds

- **YouTube Music.** Sign in with [Google's own page](accounts.md#the-sign-in-window) in a window MyTube opens, or [import the session of a browser](accounts.md#import-from-a-browser) you signed in with. You get your home feed, your library (playlists, albums, artists you follow and liked songs), likes and dislikes that reach YouTube Music, saving albums and playlists, making and deleting playlists, and a [brand channel](accounts.md#several-google-accounts-and-brand-channels) of your account if you have one. What you play joins your YouTube Music history, so its home and radios learn from it, and podcast episodes [go on where you left them](podcasts.md#with-youtube-music) on your phone. Tracks play the same signed in or out.
- **Spotify.** Spotify lets an app serve up to five accounts, and only while the app's owner has Premium, so you sign in with a Spotify app of your own: [Spotify](spotify.md) walks you through it. You can then search Spotify too, see your Spotify library, and like, save and follow there. Spotify tracks play in MyTube itself (with Premium), on Spotify Connect (with Premium), or from their recording on YouTube Music, which needs no Premium.
- **Last.fm.** Paste the API key and shared secret of a Last.fm API account of yours, and approve MyTube in your browser: what you listen to joins your Last.fm profile once half of a track, or four minutes of it, played. Scrobbles made offline are sent once you're back online.

## Next

- [Playing music](playing-music.md): search, the queue, radios, the sound and speed.
- [Your library](library.md): liked songs, mixes, downloads, history and stats.
- [Lyrics](lyrics.md): synced lyrics, translations and Latin letters.
- [Podcasts](podcasts.md): shows, episodes and where you left them.
- [Troubleshooting](troubleshooting.md), when something doesn't work.
