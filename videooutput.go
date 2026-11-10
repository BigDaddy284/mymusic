package main

// The device the videos' sound plays on, apart from the music's: the
// Video side's output button chooses it, and Windows sends the sound of
// the video player's processes there.

import (
	"log"

	"github.com/egoist/mygo/ui"
)

// videoOutputButton lists the output devices, the one the videos play on
// checked.
func (s *shell) videoOutputButton(c *ui.Context, has bool) {
	vp := &s.videos.player
	chosen := s.set.VideoOutput
	label := "Videos play on the system's default device"
	if chosen.ID != "" {
		label = "Videos play on " + chosen.Name
	}
	if vp.casting != nil {
		label = "Playing on " + vp.casting.device.Name
	}
	b := s.iconButton(c, "v-output", "speaker", label, chosen.ID != "", 18).Disabled(vp.casting != nil)
	b.Menu(func(m *ui.Menu) {
		devices := s.outputDevices()
		for _, o := range outputChoices(devices.value, devices.done && devices.err == nil, chosen) {
			item := m.Item(o.Label).Radio(o.Chosen)
			if o.Detail != "" {
				item.Subtitle(o.Detail)
			}
			if item.Chosen() && o.Device.ID != s.set.VideoOutput.ID {
				s.set.VideoOutput = o.Device
				s.save()
				s.applyVideoOutput()
			}
		}
		switch {
		case !devices.done:
			m.Item("Looking for devices…").Disabled(true)
		case devices.err != nil:
			m.Item(devices.err.Error()).Disabled(true)
		}
	})
}

// applyVideoOutput sends the video player's sound to the device chosen.
func (s *shell) applyVideoOutput() {
	web := s.videos.player.web
	if web == nil {
		return
	}
	id := s.set.VideoOutput.ID
	go func() {
		if err := s.svc.Videos.RouteOutput(web.Page().ProcessID(), id); err != nil {
			log.Printf("mytube: the videos' output device: %v", err)
			if id != "" {
				s.update(func() {
					s.notes = append(s.notes, noticeOf(err))
				})
			}
		}
	}()
}
