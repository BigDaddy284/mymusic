# The app and the desktop

- **`internal/app`** gives the windows the services (`App.Services`), tells
  them when the settings, the accounts, the add-ons or the updates change
  (`App.Watch`), makes the player, and puts MyMusic on the desktop: the
  tray, media keys, global shortcuts, deep links, notifications, updates.
- **`internal/nowplaying`** shows what the engine plays in the system's media
  controls, and passes on what they ask: the SystemMediaTransportControls
  of a hidden window, reached through WinRT without cgo. It downloads
  nothing: `internal/app` hands it each cover from the artwork store
  (`controls_cover.go`), as its bytes. The notification of a track that starts shows the
  same cover (`notecover.go`), its centre square at most 512 pixels
  across (`artwork.Square`), as the bytes MyGo's notifications take; it
  waits 2 seconds at most for it.
- **`internal/taskbar`** puts previous, play or pause, and next on the
  thumbnail of the main window's taskbar button on Windows, through the
  shell's ITaskbarList3, and takes their clicks in a window procedure in
  front of MyGo's.
- **`internal/keypad`** takes Ctrl and the keypad's +, which MyGo gives its
  windows as no key, so that they zoom in as Ctrl+= does (#251): in a
  keyboard hook of the window's thread.
- **`internal/diag`** writes the log, `mymusic.log` in the logs folder,
  rotated at a megabyte and without the cookies, hashes and tokens a line
  may hold, and has Go write a crash to `crash.log` there, which the next
  start keeps as `crash-<time>.log` and offers to report. Its
  `MemoryRecorder` records memory use while the user has Settings ›
  About › Record memory use on (`internal/app`'s `memory.go`, off at each
  start): a line of `ReadMemory` in `memory.log` each minute, the Go
  heap's figures and the process's resident size where the system tells
  it cheaply (`GetProcessMemoryInfo`, with its private bytes), and a heap
  profile (`WriteHeapProfile`) when asked, and
  once as the memory passes 1 GB, the newest five kept. Capture mode's
  `act:memory` and `act:heap` use the same. Report a problem names the
  files there are (`System.MemoryFiles`).
