# The player

- **`internal/player`** plays the queue: one goroutine owns the queue and
  what plays, and the network calls run beside it. It plays YouTube Music
  and the music on the computer through the engine. It shuffles, repeats, crossfades, starts radios when the
  queue runs out, or at once after a suggested song (`PlayWithRadio`),
  passes over the songs the user dislikes but those they
  pick, counts plays and keeps up to 500 tracks of the queue,
  around the one playing, in `queue.json`, as it changes and as MyMusic
  quits. A track that can't play for want of the network gives way at once
  to the next that plays without it, a download or a track the offline
  cache holds, as the app tells; once none is left, it waits for the
  network where it was, and plays on once it answers. A track that plays
  without the network, whose audio fails, doesn't wait for it: it is tried
  once more, at once, and then fails with why. The player pauses as the
  computer sleeps.
  As it moves on by itself it passes over, on one path
  (`skip.go`), the entries the settings skip, but those the user picks:
  songs the user dislikes, and songs YouTube labels as made with AI
  (`labels.go`). A track whose label the app doesn't know is checked
  before it's announced, beside the search for its stream, and the next
  track before it's preloaded or crossfaded into.
- **Failures** (`failures.go`): the player keeps the last 10 tracks that
  couldn't play, until MyMusic quits (`Failures`, `FailureOf`), and the
  notice of each names its own (`Notice.Failure`). Each holds the steps
  tried of the track since it loaded, in turn, with why each failed: the
  5.1 audio it gave up, the stream and its tries again, one step for each run of the same
  error, and the error that failed it, named after how the track was to
  play (`way`): the file, a download or the offline cache, the stream,
  or finding it. An error that tells its own steps, as the resolver's
  does when yt-dlp failed after YouTube's player API (`stream.bothFailed`),
  gives them in its place: `await` keeps the error of the answer the
  player takes in, for `tried` to read.
- **`internal/online`** tells the player whether the network answers: a
  check of YouTube's addresses over HTTPS, shared by its
  callers, which audio coming from YouTube answers too, and a wait for the
  network, with checks further apart the longer it lasts. The system's
  word on the network, through MyGo, only has it check at once.
