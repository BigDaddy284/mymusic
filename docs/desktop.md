# On the desktop
<!-- docs: group=using order=60 -->

MyTube on your desktop: the mini player, the tray,
media keys and media controls, global shortcuts, notifications and Discord.

Most of what this page covers is in [Settings › Desktop](settings.md#desktop),
and the global shortcuts are in [Settings › Shortcuts](settings.md#shortcuts).

## The mini player

The mini player is a small window that stays on top of your other
windows, with what plays and the buttons to control it. Open it with:

- the mini player button at the end of the player bar;
- <kbd>Ctrl</kbd> <kbd>M</kbd>, or
  **Mini player** in the command palette;
- **Mini player** in the tray icon's menu.

It shows the cover, the title, the artists with the mark of their service,
and the time and the length ("live" for a live stream), with a line of the
accent along its foot for the progress. Its buttons:

- at the top right, one button that opens MyTube's main window and one
  that closes the mini player;
- beside the time, the buttons of its [panel](#its-panel): the queue, the
  lyrics and the output device, and, after a track couldn't play, a red
  warning sign;
- [shuffle and repeat](playing-music.md#shuffle-and-repeat), in the accent
  while they're on, around previous, play or pause, and next;
- the heart, which [likes](library.md#liked-songs) what plays, and the
  thumbs-down, which [dislikes](playing-music.md#dislikes-and-skipping) it;
- the volume: a click mutes and unmutes, the mouse wheel over it turns the
  volume up and down, and a small slider shows above it while the pointer
  is over it. <kbd>Tab</kbd> reaches the button and then the slider, whose
  arrow keys set the volume; <kbd>Escape</kbd> closes it.

![The mini player playing アイドル by YOASOBI: the cover on the left, the title, the artist with the mark of YouTube Music, 0:17 of 3:33, previous, pause and next, the heart, the thumbs-down and the volume, and at the top right the buttons that open the main window and close the mini player](shots/mini-player.webp)

The line along its foot seeks: click it or drag it, or reach it with
<kbd>Tab</kbd> and use the arrow keys. Under the pointer it thickens, and
the time there shows above it.

Drag the mini player anywhere but its buttons to move it, and drag its
edges to resize it: the cover takes its height, and grows up to 180 px at
100% zoom. It opens where you left it. A right click opens a menu of
[speeds](playing-music.md#speed) for what plays. The mini player takes no
keyboard shortcuts but the zoom's, and <kbd>Escape</kbd>, which closes its
panel. It [zooms](look-and-feel.md#density-text-size-and-zoom) with the
main window, and its window grows and shrinks with it.

### Its panel

The buttons beside the time open a panel under the mini player, or above
it when it's at the foot of the screen, one at a time:

- **Queue:** the track playing and what plays after it, with your place in
  the queue, as "3 of 24". Click a track to play it, and right-click it
  for its menu, as in [Up next](playing-music.md#the-queue-and-up-next).
  A page its menu goes to, as Go to album, opens in the main window.
- **Lyrics:** the line sung now, which lights up word by word where the
  [lyrics](lyrics.md) have the words' times, and the line after it, in
  Latin letters where the lyrics panel shows them so. Lyrics without times
  show whole. For a podcast's episode, the chapter that plays, and the
  next.
- **Output device:** the [output devices](playing-music.md#the-output-device)
  of your computer, and the [cast devices](casting.md) on your network and
  your Spotify Connect devices, while the player bar shows their buttons.
  Click one to play there.
- **The red warning sign,** after a track couldn't play since the mini
  player opened: why, as its toast's
  [details](troubleshooting.md#tracks-that-fail-to-play) tell, with
  **copy**. Once you close it, the sign goes until another track fails.

The window grows to show the panel, and shrinks back as it closes, the
mini player staying where it is. Its button again, the × at the panel's
top, or <kbd>Escape</kbd> closes it. The mini player opens with the panel
it showed as it closed, but the warning sign's.

## The tray

MyTube puts an icon in the tray. Its menu has **Play / Pause**, **Next**,
**Previous**, **Show MyTube**, **Mini player**, and the item that quits
MyTube. A click on the icon shows the main window and a right click opens
the menu.

![The tray icon's menu on Windows: Play / Pause, Next and Previous, then Show MyTube and Mini player, then Exit](shots/tray-menu.webp)

### Close to the tray

With **Close to the tray** on, in Settings › Desktop, closing the main
window keeps the music playing: MyTube hides in the tray until you show
it again from the tray icon, by starting MyTube again, or with a global
shortcut of your own for **Show MyTube**. It is off on a new install, so
closing the window quits MyTube.

- To quit with Close to the tray on, choose the quit item of the tray
  icon's menu.

## Media keys and the system's media controls

MyTube shows what plays in Windows' media flyout, which shows with the
volume keys and in Quick Settings, with the cover, the title and the
artists, play or pause, previous, next and a seek bar, and takes its
buttons.

The covers there come from the
[covers MyTube keeps](library.md#the-offline-cache) on disk, so a cover
you saw in the window shows there at once, and offline too. A track [on
this computer](local-music.md) shows the cover its file holds, or the one
beside it, from MyTube's copy in its cache.

**Media keys**, in Settings › Desktop and on by default, has the keyboard's
play or pause, next and previous keys control MyTube, and its stop key
pause it. While a track loads, they and the system's controls do nothing:
a track can't be paused or moved until you hear it.

### The taskbar thumbnail on Windows

On Windows, the thumbnail that shows while the pointer rests on MyTube's
taskbar button has three buttons: previous, play or pause, and next. They
are off while there's no track to play.

### The window title

The main window's title names the track that plays, for the taskbar, the
window switcher and screen readers, and starts with "Paused:" while it's
paused. Turn off **Show what plays in the window title** in Settings ›
Desktop to keep it to MyTube's name.

## Keep the window on top

**Keep on top**, in Settings › Desktop and off at first, keeps MyTube's
main window above your other windows, for when you work beside it.
MyTube remembers it between starts. Turn it on or off with:

- the **Keep on top** switch in Settings › Desktop;
- **Keep the window on top, or not** in the command palette, which has no
  keys at first: Settings › Shortcuts can give it some. A toast says
  which it did;
- `"keepOnTop": true` in `settings.json`.

The [mini player](#the-mini-player) stays on top either way.

## Global shortcuts

Global shortcuts are keys of your own that control MyTube whichever
window is in front, even with its window closed to the tray. Record them in
the **From any app** group at the top of
[Settings › Shortcuts](settings.md#change-a-shortcut), for:

- Play or pause
- Next track and Previous track
- Volume up, Volume down and Mute
- Like what plays and Dislike what plays
- Show MyTube

![Settings › Shortcuts: the From any app group at the top, "These keys control MyTube whichever window is in front", with Play or pause, Next track, Previous track and Volume up, each with none and a Set button](shots/shortcuts.webp)

A new install has none. To set one, choose **set** on its row and press
the keys; <kbd>Escape</kbd> cancels and <kbd>Backspace</kbd> removes it.

- A global shortcut needs <kbd>Ctrl</kbd>, <kbd>Alt</kbd> or the Windows
  key. <kbd>F1</kbd> to <kbd>F11</kbd> may stand alone; Windows keeps
  <kbd>F12</kbd> for itself.
- While MyTube runs, those keys do nothing else, in MyTube's own window
  too: a command of the window that had them gives them up.
- Keys the system refuses, as when another app has them, stay yours: the
  row says why, with **try again**, and MyTube tries again at each
  start.
- Liking or disliking from another window shows a notification, while
  notifications are on and neither MyTube's window nor the mini player
  is in front.

The media keys work as they do without any global shortcut.

## Notifications

With **Notifications** on, in Settings › Desktop and on by default,
MyTube shows a notification when a track starts while its window isn't
in front: the title, the artists and the album, with the cover. A click
on it shows MyTube. None shows while the mini player is in front: it
shows the track already. The tray icon's tip names what plays too.

- The cover is its centre square, as the player bar shows it: the middle
  of a music video's thumbnail. It takes the place of the tray icon.
- It is the cover the system's media controls show, from the covers
  MyTube keeps: a track on this computer shows the cover its file holds,
  or the one beside it.
- MyTube waits up to 2 seconds for a cover it hasn't downloaded yet:
  without one by then, the notification shows without it. When the next
  track starts meanwhile, only the next one's shows.

## Discord

While the Discord app runs on your computer, your Discord profile shows
what you listen to: the title, the artists, the cover and the time left.
Pausing clears it, and a live stream shows how long it has played. You
can [choose what it shows](#what-your-discord-status-shows). Discord's
member list says "Listening to MyTube". There's nothing to set up: turn
it off with **Show what plays on Discord** in Settings › Desktop.

- **Discord application ID**, under the switch, shows the name of a Discord
  application of your own in place of MyTube's: create one in Discord's
  [developer portal](https://discord.com/developers/applications) and paste
  its application ID.
- Under them, a **preview** shows your status as Discord's card does, of
  what plays, or of a song while nothing does, and follows what you choose
  below it, the lines as you type them.
- Discord's own Spotify status comes from your Discord account's link to
  Spotify. To show only MyTube, turn off "Display Spotify as your status"
  in Discord's Connections.
- MyTube finds Discord's Canary build too.

### What your Discord status shows

The rows under the preview choose what your profile says. At first it
says what it always did: the title, the artists, the cover and the time
left, and a pause clears it.

| Row | What it does |
|---|---|
| First line, Second line | What the two lines say. Type text and placeholders: `{title}`, `{artist}`, `{album}` and `{service}` (YouTube Music, Spotify or my computer), as `{artist} · {album}`. A placeholder at a line's start or end with nothing to show, as the album of a single, takes the separator beside it along. Anything else in braces shows as typed, and the row says so in red. A line holds 128 characters, as many as Discord shows; an empty one shows nothing. **reset** puts back `{title}` or `{artist}`. |
| Cover | The track's cover, where its service has one online. A file on your computer has none Discord can read. On at first. |
| Service badge | A small round mark at the cover's corner: a green disc for Spotify, a red square for YouTube Music, a record for a file on your computer. Without a cover, the badge takes its place. Off at first. |
| Time left | How long the track has played and how long it has left, as a bar; a live stream, how long it has played. On at first. |
| Listen button | A **Listen on Spotify** or **Listen on YouTube Music** button that opens the track there. Discord shows it to the others who see your profile, not always to you. A file on your computer has none. Off at first. |
| Keep while paused | A paused track stays, marked with a pause badge, whose tip says Paused, and without its time. Off, a pause clears your status, as before. |

Discord reads the badges from MyTube's GitHub repository, and an
application of your own shows them as MyTube's does.

## Start at login

**Start at login**, in Settings › Desktop, starts MyTube when you sign in
to your computer, with its window hidden in the tray, from where the tray
icon shows it.

## Sleep

MyTube pauses as the computer goes to sleep, and stays paused when it
wakes, ready where it was. A Spotify Connect device elsewhere plays on.

## mytube:// links

Links that start with `mytube://` open a page of MyTube, or play a
track, from a browser or another app:

| Link | What it does |
|---|---|
| `mytube://search?q=lofi+piano` | searches for "lofi piano" |
| `mytube://album/<id>`, `mytube://artist/<id>`, `mytube://playlist/<id>` | opens an album, an artist or a playlist |
| `mytube://library`, `mytube://explore`, `mytube://stats`, `mytube://marketplace`, `mytube://settings` | opens that page |
| `mytube://play/<track id>` | plays a track |

An ID is MyTube's own: the service's ID after `yt:` for YouTube Music or
`sp:` for Spotify. Opening MyTube a second time, or a link while it runs,
opens the page in the window that's there.
