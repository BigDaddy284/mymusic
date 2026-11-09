package output

// The 5.1 layouts the outputs open: the channels of a Stream of six
// channels are, in order, left, right, center, LFE, left surround and
// right surround, as WAVE files and ac3's decoder order them.

// Surround is the channels of a Stream of 5.1 surround.
const Surround = 6

// The speaker positions of WAVEFORMATEXTENSIBLE's dwChannelMask.
const (
	speakerFrontLeft    = 0x1
	speakerFrontRight   = 0x2
	speakerFrontCenter  = 0x4
	speakerLowFrequency = 0x8
	speakerBackLeft     = 0x10
	speakerBackRight    = 0x20
	speakerSideLeft     = 0x200
	speakerSideRight    = 0x400

	// speakerStereo is KSAUDIO_SPEAKER_STEREO.
	speakerStereo = speakerFrontLeft | speakerFrontRight
	// speaker51 is KSAUDIO_SPEAKER_5POINT1, of back speakers, and
	// speaker51Side KSAUDIO_SPEAKER_5POINT1_SURROUND, of side ones.
	speaker51     = speakerStereo | speakerFrontCenter | speakerLowFrequency | speakerBackLeft | speakerBackRight
	speaker51Side = speakerStereo | speakerFrontCenter | speakerLowFrequency | speakerSideLeft | speakerSideRight
)

// surroundMask returns the channel mask of a WASAPI stream of 5.1 on a
// device of channels whose mix format has mask, 0 when unknown: the 5.1
// of the device's surround speakers, side ones where it has them, as 5.1
// and 7.1 systems have, else back ones. It reports false for a device of
// fewer than six channels, or without a center, an LFE and a pair of
// surrounds, which a Stream plays stereo on.
func surroundMask(channels int, mask uint32) (uint32, bool) {
	if channels < Surround {
		return 0, false
	}
	if mask == 0 {
		// A device that names no speakers: six channels or more are taken
		// in the order of the masks, the side speakers' for 5.1.
		return speaker51Side, true
	}
	if mask&(speakerStereo|speakerFrontCenter|speakerLowFrequency) != speakerStereo|speakerFrontCenter|speakerLowFrequency {
		return 0, false
	}
	switch {
	case mask&(speakerSideLeft|speakerSideRight) == speakerSideLeft|speakerSideRight:
		return speaker51Side, true
	case mask&(speakerBackLeft|speakerBackRight) == speakerBackLeft|speakerBackRight:
		return speaker51, true
	}
	return 0, false
}
