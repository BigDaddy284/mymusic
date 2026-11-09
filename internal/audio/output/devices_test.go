package output

import (
	"slices"
	"testing"
	"time"
)

// TestEndpointDevices lists WASAPI's endpoints by their friendly names,
// sorted as people read them, the default one marked: an endpoint without
// a name shows its ID, and one without an ID isn't listed.
func TestEndpointDevices(t *testing.T) {
	got := endpointDevices([]Device{
		{ID: "{0.0.0.00000000}.{b}", Name: "Speakers (Realtek(R) Audio)"},
		{ID: "{0.0.0.00000000}.{a}", Name: "headphones (USB DAC)"},
		{ID: "{0.0.0.00000000}.{c}"},
		{Name: "No ID"},
	}, "{0.0.0.00000000}.{b}")
	want := []Device{
		{ID: "{0.0.0.00000000}.{a}", Name: "headphones (USB DAC)"},
		{ID: "{0.0.0.00000000}.{b}", Name: "Speakers (Realtek(R) Audio)", Default: true},
		{ID: "{0.0.0.00000000}.{c}", Name: "{0.0.0.00000000}.{c}"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("devices %+v, want %+v", got, want)
	}
}

// TestMover tells Moved where a stream that asks for a device plays, as it
// starts and as that changes, once for each change, and never for a
// stream of the default device.
func TestMover(t *testing.T) {
	var told []bool
	m := newMover(Options{Device: "headphones", Moved: func(fallback bool) { told = append(told, fallback) }})
	for _, fallback := range []bool{false, false, true, true, false} {
		m.tell(fallback)
	}
	if want := []bool{false, true, false}; !slices.Equal(told, want) {
		t.Errorf("told %v, want %v", told, want)
	}
	told = nil
	m = newMover(Options{Moved: func(fallback bool) { told = append(told, fallback) }})
	m.tell(false)
	m.tell(true)
	if told != nil {
		t.Errorf("a stream of the default device told %v", told)
	}
	newMover(Options{Device: "headphones"}).tell(true) // without Moved
}

// TestAvoid keeps a stream off the device it asks for for a while after
// the device failed to open.
func TestAvoid(t *testing.T) {
	var a avoid
	now := time.Now()
	if !a.may(now) {
		t.Fatal("a device that never failed is avoided")
	}
	a.failed(now)
	if a.may(now.Add(avoidFor - time.Millisecond)) {
		t.Error("a device that failed is tried again at once")
	}
	if !a.may(now.Add(avoidFor)) {
		t.Errorf("a device that failed isn't tried again after %v", avoidFor)
	}
}
