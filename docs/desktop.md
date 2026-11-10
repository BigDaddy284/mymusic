# On the desktop
<!-- docs: group=using order=60 -->

MyTube on your desktop: the mini player, the tray,
media keys and media controls, global shortcuts and notifications.

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
  of your computer, and the [cast devices](casting.md) on your network,
  while the player bar shows their buttons.
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

## Start at login

**Start at login**, in Settings › Desktop, starts MyTube when you sign in
to your computer, with its window hidden in the tray, from where the tray
icon shows it.

## Sleep

MyTube pauses as the computer goes to sleep, and stays paused when it
wakes, ready where it was.

## mytube:// links

Links that start with `mytube://` open a page of MyTube, or play a
track, from a browser or another app:

| Link | What it does |
|---|---|
| `mytube://search?q=lofi+piano` | searches for "lofi piano" |
| `mytube://album/<id>`, `mytube://artist/<id>`, `mytube://playlist/<id>` | opens an album, an artist or a playlist |
| `mytube://library`, `mytube://explore`, `mytube://stats`, `mytube://settings` | opens that page |
| `mytube://play/<track id>` | plays a track |

An ID is MyTube's own: YouTube Music's ID after `yt:`. Opening MyTube a second time, or a link while it runs,
opens the page in the window that's there.
