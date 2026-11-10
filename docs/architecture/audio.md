# The audio engine

- **`internal/audio`** is MyTube's audio engine: two decks, each a
  goroutine that reads a track, Opus in WebM with `webm` and `opus`, AAC
  in MP4 with `mp4`, `aac` and `resample`, or the AAC of a live stream's
  packed audio, its ID3 tags passed over, which has no length and can't
  seek, or a file of the music on this computer (`pcm.go`): FLAC with
  `flac`, Opus and Vorbis in Ogg with `ogg`, `opus` and jfreymuth/vorbis,
  MP3 with go-mp3, the delay and padding of its LAME tag dropped, and
  WAV, each converted to 48 kHz and mixed down to stereo, as its first
  bytes tell, and decodes it into
  a ring the audio thread plays from, at 48 kHz, each turning a
  track louder than -14 LKFS down to it, from the loudness its `Measured`
  source tells (Even loudness), and the preamp, equalizer, levelling
  compressor, volume and analyser of the web player MyTube had, ported
  from Chromium's Web Audio, which they match sample for sample. A deck
  plays its track faster or slower at its own pitch by WSOLA
  (`stretch.go`), the waveform-similarity overlap-add of Verhelst and
  Roelands, between its ring and the mix; at 1× it plays the ring as it
  is, as a live stream always plays, with nothing ahead of its live edge
  to play faster. The six channels of 5.1 surround stretch alike, each
  window found on the mix of them all. A deck fails in the words of the
  web player's audio element, but a read of a `Local` source, a file on
  this computer or a download, that fails is the file's failure ("Can't
  read the file"), not the network's.
  `MYTUBE_STRETCH_WAV=<dir> go test ./internal/audio -run StretchListen`
  writes the test files at each speed, to listen to.
  `MYTUBE_LIVE=1 go test ./internal/audio -run Live` plays songs of both
  codecs, and a 24/7 live stream, as the app finds them on YouTube,
  against ffmpeg's decode. `Probe` tells whether a deck plays a file, and
  how long, reading little of it, for the index of the music on this
  computer: of an MP3, the frames its Xing, Info or VBRI header counts,
  less its LAME tag's delay and padding, else its first frame's bitrate
  and its size. With ffmpeg on the PATH, the tests play a file of each
  kind it makes against its decode.
  `internal/audio/output` plays it on the system's output, without cgo,
  through WASAPI. `Devices` lists the output devices (`devices.go` puts
  together what WASAPI tells, so that tests give it what WASAPI would), and a stream whose `Options`
  name one plays on it while it's there, checks every second that it
  still is, plays on the default device meanwhile, and tells `Moved`
  each time it moves. `internal/app/outputdevice.go` keeps the device
  Settings › Playback chose, opens the engine's outputs on it, has the
  engine open its output again (`Engine.ReopenOutput`) when the choice
  changes, and tells the windows when the device goes away and comes
  back. The player bar's device button and `outputChoices` and
  `outputMenu` in the root package's `outputdevice.go` show the choice.
  `MYTUBE_AUDIO=1 go test ./internal/audio/output -run TestStream`
  plays on the default device by its ID, and on one that isn't there.
- **`internal/volume`** has the volume's curves, which Settings › Playback
  › Volume chooses between (#383): Natural, each tenth of the slider 6 dB
  down to -40 dB, and Linear, whose gain is the slider's position. The
  player keeps the position, as settings.json does whatever the curve
  (`Saved`, `Of`). The app gives the engine the curve
  (`Engine.SetVolumeCurve`, in `applyEngine`).
