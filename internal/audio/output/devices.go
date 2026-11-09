package output

// The output devices: Devices lists them, and a Stream plays on the one
// its Options name, or on the system's default device. What WASAPI tells
// of its devices is put together here, apart from the calls that ask it,
// so that tests give it what WASAPI would.

import (
	"cmp"
	"slices"
	"strings"
	"time"
)

// Device is an audio output device of the system, as Devices lists it.
type Device struct {
	// ID tells the device from the others, from one run to the next: its
	// WASAPI endpoint ID.
	ID string
	// Name is what the system calls it, as its sound settings show.
	Name string
	// Default is set on the device the system plays on by default, where
	// the system tells which it is.
	Default bool
}

// Options are those of a Stream: what Start takes, and the device it
// plays on.
type Options struct {
	// Channels is what the stream asks for: 2, or Surround for 5.1, which
	// it plays where the device mixes 5.1 or more, else stereo.
	Channels int
	// Render is called from the stream's goroutine for each period of the
	// device, about every 10 ms, to fill out with frames of the channels it
	// tells, interleaved: those of the device that plays, which change as
	// the device does. It must not block.
	Render func(out []float32, channels int)
	// Failed, when set, is called with why when the stream stops of
	// itself, as when no output device is left, from the stream's
	// goroutine once the stream has ended; Stop is not needed then.
	Failed func(error)
	// Device is the ID of the device to play on, as Devices tells it; ""
	// plays on the system's default device, and follows it as it changes.
	// While the device isn't there, or can't play, the stream plays on the
	// default device, and it goes back to the device once it can.
	Device string
	// Moved, when set, is told whether the stream plays on the default
	// device because Device isn't there or can't play: as the stream
	// starts, and each time that changes. It is never told while Device is
	// "". It is called from the stream's goroutines, and must not block.
	Moved func(fallback bool)
}

// Start starts a stream of channels on the system's default device, as
// Open does with Options of channels, render and failed.
func Start(channels int, render func(out []float32, channels int), failed func(error)) (*Stream, error) {
	return Open(Options{Channels: channels, Render: render, Failed: failed})
}

const (
	// checkInterval is how often a stream checks where it should play:
	// which device is the default, and whether the device it asks for is
	// there.
	checkInterval = time.Second
	// avoidFor is how long a stream plays on the default device after the
	// device it asks for was there but failed to open, before it tries
	// that device again.
	avoidFor = 10 * time.Second
)

// mover tells a stream's Moved where the stream plays, once for each
// change. One goroutine of the stream uses it at a time.
type mover struct {
	moved          func(bool)
	told, fallback bool
}

func newMover(o Options) *mover {
	if o.Device == "" {
		return &mover{}
	}
	return &mover{moved: o.Moved}
}

// tell tells Moved whether the stream plays on the default device in place
// of the one it asks for, unless it told so last.
func (m *mover) tell(fallback bool) {
	if m.moved == nil || m.told && m.fallback == fallback {
		return
	}
	m.told, m.fallback = true, fallback
	m.moved(fallback)
}

// avoid is when a stream may try the device it asks for again, after it
// failed to open: the zero time while it may.
type avoid struct{ until time.Time }

// failed has the stream play on the default device for avoidFor.
func (a *avoid) failed(now time.Time) { a.until = now.Add(avoidFor) }

// may reports whether the stream may try the device it asks for.
func (a *avoid) may(now time.Time) bool { return !now.Before(a.until) }

// sortDevices sorts devices by name, as people read them, then by ID.
func sortDevices(devices []Device) []Device {
	slices.SortStableFunc(devices, func(a, b Device) int {
		return cmp.Or(cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)), cmp.Compare(a.ID, b.ID))
	})
	return devices
}

// endpointDevices returns the devices of WASAPI's active render endpoints,
// whose names are their friendly names, the default one marked: the ID
// stands for a name the endpoint doesn't tell.
func endpointDevices(endpoints []Device, defaultID string) []Device {
	out := make([]Device, 0, len(endpoints))
	for _, d := range endpoints {
		if d.ID == "" {
			continue
		}
		if d.Name == "" {
			d.Name = d.ID
		}
		d.Default = d.ID == defaultID
		out = append(out, d)
	}
	return sortDevices(out)
}
