package main

// Settings › Videos: SponsorBlock's skipping, and the channels, words and
// Shorts kept off the Video side's pages.

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/settings"
)

func init() {
	addSettings("videos", 10, (*shell).sponsorBlockSettings)
	addSettings("videos", 20, (*shell).blockSettings)
}

// segmentActions are the actions a kind of segment can take, as the
// settings name them, and those of a highlight, a moment.
var (
	segmentActions      = []string{settings.SegmentSkip, settings.SegmentMute, settings.SegmentShow, settings.SegmentOff}
	segmentActionLabels = []string{"skip", "mute", "show", "off"}
	pointActions        = []string{settings.SegmentShow, settings.SegmentOff}
	pointActionLabels   = []string{"show", "off"}
)

// sponsorBlockSettings choose what the player does with each kind of
// segment SponsorBlock's users marked.
func (s *shell) sponsorBlockSettings(c *ui.Context) {
	s.groupWith(c, "Skip the fluff", func() {
		s.linkOf(c, "SponsorBlock", "https://sponsor.ajay.app")
	}, func() {
		on := !s.set.SponsorBlockOff
		s.row(c, "Skip with SponsorBlock",
			"Skips sponsor reads, intros, reminders to subscribe and the like, which SponsorBlock's users marked. Show marks them on the seek bar, without skipping. The app asks by a part of a hash of the video, so SponsorBlock doesn't learn what you watch.",
			func() {
				if s.switchOf(c.Key("sponsorblock"), &on).Label("Skip with SponsorBlock").Changed() {
					s.set.SponsorBlockOff = !on
					s.save()
				}
			}, "sponsor", "skip", "segments", "intro", "outro")
		if !on {
			return
		}
		for _, kind := range settings.SegmentKinds {
			actions, labels := segmentActions, segmentActionLabels
			if kind.Point {
				actions, labels = pointActions, pointActionLabels
			}
			now := s.set.SegmentAction(kind.ID)
			sel := 0
			for i, a := range actions {
				if a == now {
					sel = i
				}
			}
			s.row(c, kind.Name, kind.About, func() {
				if s.segmented(c, "segment-"+kind.ID, kind.Name, &sel, labels...) {
					if s.set.Segments == nil {
						s.set.Segments = map[string]string{}
					}
					s.set.Segments[kind.ID] = actions[sel]
					s.save()
				}
			}, "sponsorblock")
		}
	})
}

// blockSettings keep channels, words and Shorts off the Video side.
func (s *shell) blockSettings(c *ui.Context) {
	s.group(c, "Keep off my pages", func() {
		s.row(c, "Blocked channels", "Their videos never show: on Home, in search, Random and Up next. Separate them with commas, or block one from a video's menu.", func() {
			s.blockListInput(c, "blocked-channels", "Blocked channels", &s.set.BlockedChannels)
		}, "block", "hide", "channel")
		s.row(c, "Blocked words", "Videos whose titles hold any of these, in any case, never show. Separate them with commas.", func() {
			s.blockListInput(c, "blocked-words", "Blocked words", &s.set.BlockedWords)
		}, "block", "hide", "keyword", "title")
		hide := !s.set.ShowShorts
		s.row(c, "Hide Shorts", "Leaves YouTube's Shorts out of every page of the Video side.", func() {
			if s.switchOf(c.Key("hide-shorts"), &hide).Label("Hide Shorts").Changed() {
				s.set.ShowShorts = !hide
				s.save()
			}
		}, "shorts")
	})
}

// blockListInput edits list as words separated by commas.
func (s *shell) blockListInput(c *ui.Context, key, label string, list *[]string) {
	in, _ := s.settingInput(c, key, strings.Join(*list, ", "), func(v string) {
		*list = splitList(v)
		s.save()
	})
	in.Width(280).Label(label).Placeholder("none")
}

// splitList splits text at its commas, each part trimmed, the blank ones
// left out.
func splitList(text string) []string {
	var out []string
	for _, p := range strings.Split(text, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
