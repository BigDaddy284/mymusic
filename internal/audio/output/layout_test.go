package output

import (
	"testing"
)

// TestSurroundMask checks the channel masks of WASAPI streams of 5.1:
// those of the device's surround speakers, side ones before back ones, and
// none on a device without 5.1, which a stream plays stereo on.
func TestSurroundMask(t *testing.T) {
	const speaker71 = speaker51 | speakerSideLeft | speakerSideRight
	const speaker71Wide = speaker51 | 0x40 | 0x80 // front left and right of center
	for _, c := range []struct {
		name     string
		channels int
		mask     uint32
		want     uint32
		ok       bool
	}{
		{"stereo", 2, speakerStereo, 0, false},
		{"quad", 4, speakerStereo | speakerBackLeft | speakerBackRight, 0, false},
		{"5.1 of back speakers", 6, speaker51, speaker51, true},
		{"5.1 of side speakers", 6, speaker51Side, speaker51Side, true},
		{"7.1", 8, speaker71, speaker51Side, true},
		{"7.1 wide", 8, speaker71Wide, speaker51, true},
		{"six channels of no speakers", 6, 0, speaker51Side, true},
		{"six channels without an LFE", 6, speakerStereo | speakerFrontCenter | speakerBackLeft | speakerBackRight | 0x100, 0, false},
		{"five channels of 5.1's mask", 5, speaker51, 0, false},
	} {
		got, ok := surroundMask(c.channels, c.mask)
		if got != c.want || ok != c.ok {
			t.Errorf("%s: mask %#x, %v; want %#x, %v", c.name, got, ok, c.want, c.ok)
		}
	}
	// The masks a Stream opens, by the order of its channels.
	if speaker51 != 0x3F || speaker51Side != 0x60F {
		t.Errorf("KSAUDIO_SPEAKER_5POINT1 %#x, _SURROUND %#x", speaker51, speaker51Side)
	}
}
