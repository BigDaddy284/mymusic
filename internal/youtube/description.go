package youtube

// A video's description, as YouTube's watch page shows it under the
// video: its full views, its date, the links beside them, and its text,
// whose links open web pages, other videos, channels and playlists, and
// whose times seek the video itself.

import (
	"net/url"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// VideoAbout is what the watch page shows under a video.
type VideoAbout struct {
	// Title, Channel and ChannelID are the video's, as the page tells
	// them: for a video opened by its ID alone.
	Title     string `json:"title,omitempty"`
	Channel   string `json:"channel,omitempty"`
	ChannelID string `json:"channelId,omitempty"`
	// Views is the whole count, "12,345 views", and Date the day the
	// video came out, as YouTube writes it: "11 Nov 2026", "Premiered
	// 3 Jan 2025".
	Views string `json:"views,omitempty"`
	Date  string `json:"date,omitempty"`
	// Tags are the links YouTube shows beside the date: the video's
	// hashtags, or the podcast it is of.
	Tags []media.TextRun `json:"tags,omitempty"`
	// Description is the video's description, line by line.
	Description [][]media.TextRun `json:"description,omitempty"`
}

// parseVideoAbout reads what the watch page n of the video id shows
// under it.
func parseVideoAbout(n node, id string) VideoAbout {
	if id == "" {
		id = n.get("currentVideoEndpoint", "watchEndpoint", "videoId").str()
	}
	prim := n.find("videoPrimaryInfoRenderer")
	sec := n.find("videoSecondaryInfoRenderer")
	head := n.find("videoDescriptionHeaderRenderer")
	a := VideoAbout{
		Title: prim.get("title").text(),
		Views: prim.get("viewCount", "videoViewCountRenderer", "viewCount").text(),
		Date:  prim.get("dateText").text(),
	}
	if a.Title == "" {
		a.Title = head.get("title").text()
	}
	if a.Views == "" {
		a.Views = head.get("views").text()
	}
	if a.Date == "" {
		a.Date = head.get("publishDate").text()
	}
	owner := sec.get("owner", "videoOwnerRenderer")
	a.Channel = owner.get("title").text()
	a.ChannelID = owner.get("navigationEndpoint", "browseEndpoint", "browseId").str()
	if a.ChannelID == "" {
		a.ChannelID = owner.get("title", "runs", 0, "navigationEndpoint", "browseEndpoint", "browseId").str()
	}
	if a.Channel == "" {
		a.Channel = head.get("channel").text()
	}
	for _, r := range prim.get("superTitleLink").runs() {
		run := endpointRun(r.get("text").str(), r.get("navigationEndpoint"), id)
		if strings.TrimSpace(run.Text) != "" {
			a.Tags = append(a.Tags, run)
		}
	}
	var runs []media.TextRun
	switch {
	case sec.get("attributedDescription", "content").exists():
		runs = attributedRuns(sec.get("attributedDescription"), id)
	case n.find("attributedDescriptionBodyText").get("content").exists():
		runs = attributedRuns(n.find("attributedDescriptionBodyText"), id)
	default:
		for _, r := range sec.get("description").runs() {
			runs = append(runs, endpointRun(r.get("text").str(), r.get("navigationEndpoint"), id))
		}
	}
	a.Description = media.Lines(runs)
	return a
}

// attributedRuns reads an attributed string of a description: its
// content, whose command runs, by their start and length in UTF-16 code
// units, are its links.
func attributedRuns(ad node, id string) []media.TextRun {
	text := utf16.Encode([]rune(ad.get("content").str()))
	type span struct {
		start, end int
		cmd        node
	}
	var spans []span
	for _, r := range ad.get("commandRuns").list() {
		start, length := r.get("startIndex").number(), r.get("length").number()
		cmd := r.get("onTap", "innertubeCommand")
		if length <= 0 || start < 0 || start+length > len(text) || !cmd.exists() {
			continue
		}
		spans = append(spans, span{start, start + length, cmd})
	}
	slices.SortStableFunc(spans, func(a, b span) int { return a.start - b.start })
	var out []media.TextRun
	plain := func(from, to int) {
		if from < to {
			out = append(out, media.TextRun{Text: string(utf16.Decode(text[from:to]))})
		}
	}
	at := 0
	for _, sp := range spans {
		if sp.start < at {
			continue
		}
		plain(at, sp.start)
		out = append(out, endpointRun(string(utf16.Decode(text[sp.start:sp.end])), sp.cmd, id))
		at = sp.end
	}
	plain(at, len(text))
	return out
}

// endpointRun is a run of text whose command, cmd, opens what it links:
// a web page, another video, a channel or a playlist on YouTube, which
// its Link gives as YouTube's address, or a time of the video id itself,
// which it seeks. A run without a command is plain text.
func endpointRun(text string, cmd node, id string) media.TextRun {
	run := media.TextRun{Text: strings.ReplaceAll(text, " ", " ")}
	if !cmd.exists() {
		return run
	}
	if u := cmd.get("urlEndpoint", "url").str(); u != "" {
		run.Link = descriptionLink(u)
		return run
	}
	if w := cmd.get("watchEndpoint"); w.exists() {
		vid, at := w.get("videoId").str(), w.get("startTimeSeconds")
		if (vid == "" || vid == id) && at.exists() {
			run.Seek, run.AtMs = true, at.number()*1000
			return run
		}
		if vid != "" {
			q := url.Values{"v": {vid}}
			if list := w.get("playlistId").str(); list != "" {
				q.Set("list", list)
			}
			if s := at.number(); s > 0 {
				q.Set("t", strconv.Itoa(s)+"s")
			}
			run.Link = "https://www.youtube.com/watch?" + q.Encode()
			return run
		}
	}
	if vid := cmd.get("reelWatchEndpoint", "videoId").str(); vid != "" {
		run.Link = "https://www.youtube.com/shorts/" + url.PathEscape(vid)
		return run
	}
	if b := cmd.get("browseEndpoint", "browseId").str(); b != "" {
		switch {
		case strings.HasPrefix(b, "UC"):
			run.Link = "https://www.youtube.com/channel/" + url.PathEscape(b)
			return run
		case strings.HasPrefix(b, "VL"):
			run.Link = "https://www.youtube.com/playlist?list=" + url.QueryEscape(b[2:])
			return run
		}
	}
	// Else the page of YouTube's website the command opens, as a
	// hashtag's.
	if web := cmd.get("commandMetadata", "webCommandMetadata", "url").str(); strings.HasPrefix(web, "/") {
		run.Link = media.WebLink("https://www.youtube.com" + web)
	} else if web != "" {
		run.Link = descriptionLink(web)
	}
	return run
}
