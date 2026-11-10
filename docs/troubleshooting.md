# Troubleshooting
<!-- docs: group=start order=90 -->

MyTube writes what goes wrong to `mytube.log`, without your sign-ins:
[Settings › About](settings.md#about) › open logs folder shows it, beside
the reports of any crash, in
[the logs folder](data-and-privacy.md#the-logs-folder). Settings › About ›
report a problem opens a new issue on GitHub with your version, system,
tools, theme and plugins filled in; attach the log after reading it. After
a crash, MyTube offers to report it the next time it starts.

![Settings › About: the version of MyTube and whether it's the newest, with the Report a problem and Open logs folder buttons under it, and the Updates group below](shots/about.webp)

## Tracks that fail to play

When a track can't play, a toast says why, and MyTube moves on to the
next one. Its **details** shows what MyTube tried, in turn, and why each
try failed, with the time:

- YouTube's player API, then yt-dlp, which find a track's stream;
- the stream itself, as many times as MyTube tried it again;
- a download or the offline cache, or a file on this computer.

**copy** puts it all on the clipboard, as text ready for an issue, with
MyTube's version and your system. It leaves out the names of your
accounts, cookies and tokens, your computer's addresses, and the name of
your home folder, which a path shows as `~`.

The last 10 stay until MyTube closes, in
[Settings › About](settings.md#about) › Recent playback errors, each with
its **details**. The keyboard reaches **details** in the toast with
<kbd>Tab</kbd>, and **copy** has the focus as the details open.

## Common problems

- **A YouTube Music track doesn't play.** It says why: some are not
  available in your country, or only to signed-in adults on YouTube's own
  website. The tracks MyTube can't find the audio of itself, such as those
  made for children, [need yt-dlp](streams.md#when-yt-dlp-is-needed): check
  Settings › Playback that it is installed, and a JavaScript runtime found. MyTube updates the yt-dlp it
  installed every day, as an old one stops finding streams; one of your
  own, chosen there, you keep up to date yourself (`yt-dlp -U`, or your
  package manager).
- **The JavaScript runtime is not found.** Opened from the Start menu or
  a shortcut, MyTube does not get the `PATH` of your terminal. It also
  looks where Deno, Bun, Volta, Scoop and WinGet install;
  elsewhere, choose the runtime in Settings › Playback.
- **YouTube Music signed you out.** Settings › Accounts says so when YouTube
  [ends the session](accounts.md#when-youtube-ends-the-session), as it does
  when you sign out of Google everywhere or change your password. Sign out
  there, then sign in again.
- **Google refuses the sign-in window** with "This browser or app may not
  be secure". Google decides this; try again later, or from another
  network. Or sign in with a browser instead, in Settings › Accounts:
  [import the session](accounts.md#import-from-a-browser) of a browser you
  signed in to YouTube Music with, or
  [paste the Cookie header](accounts.md#paste-a-cookie-header) of a request
  of a signed-in music.youtube.com tab.
- **Importing a browser's session fails.** Close the browser first: while
  it runs, it may hold its cookies, or not have saved a new sign-in yet.
  Recent Chrome on Windows seals its cookies against other programs:
  import from Firefox there, or paste a cookie header. Zen and Helium are
  read from their profile folders: give the folder when MyTube doesn't
  find it.
- **settings.json has a mistake.** MyTube tells where, as "line 4,
  column 1: a comma before } has nothing after it", and goes on with the
  settings it had, without saving to the file, until you fix it. If it
  quits first, it keeps your text as `settings.json.broken` beside it, and
  saves its settings in its place. A file that doesn't read as MyTube
  starts is kept the same way, and MyTube starts with the settings of a
  new install: see [Settings](settings.md#when-it-has-a-mistake).
- **MyTube takes too much memory.** Turn on Settings › About › Record
  memory use, and use MyTube as when it grew: it writes how it uses its
  memory to `memory.log` each minute, and saves a heap profile once it
  passes 1 GB, or when you press save a heap profile. Report a problem
  then tells of them: attach them to the issue. They hold only how
  MyTube uses its memory ([the logs folder](data-and-privacy.md#the-logs-folder)).
  It's off again at the next start.

## Known issues

- MyTube reads YouTube Music through the interfaces of its own website,
  and plays it through YouTube's player API and yt-dlp. When YouTube
  changes them, searching, playing or signing in can stop working until
  MyTube or yt-dlp catches up.
- YouTube Music has no sign-in for other apps: MyTube signs in with a
  browser session, which can do with your account what YouTube Music can.
  Treat it like a password; sign out in Settings › Accounts to delete it.
- Tracks YouTube restricts to some countries, to signed-in adults or to a
  channel's members do not play, as streams are found signed out.
- The Windows installer is not signed yet, so SmartScreen may ask you to
  allow the first launch.
- The portable zip does not register
  [`mytube://` links](desktop.md#mytube-links), and Windows on Arm is
  less tested than the rest.
