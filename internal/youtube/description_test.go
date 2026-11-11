package youtube

import "testing"

// TestParseVideoAbout reads what the watch page shows under a video: its
// views, date and podcast link, and its description, whose links, by
// UTF-16 offsets, open web pages, other videos and channels, and whose
// times seek the video itself.
func TestParseVideoAbout(t *testing.T) {
	n := decode(t, `{"contents":{"twoColumnWatchNextResults":{"results":{"results":{"contents":[
	 {"videoPrimaryInfoRenderer":{"title":{"runs":[{"text":"A talk"}]},
	  "viewCount":{"videoViewCountRenderer":{"viewCount":{"simpleText":"541 views"},"shortViewCount":{"simpleText":"541 views"}}},
	  "dateText":{"simpleText":"11 Nov 2026"},
	  "superTitleLink":{"runs":[{"text":"The Podcast","navigationEndpoint":{"browseEndpoint":{"browseId":"VLPL1"}}},{"text":" "}]}}},
	 {"videoSecondaryInfoRenderer":{"owner":{"videoOwnerRenderer":{"title":{"runs":[{"text":"Chan","navigationEndpoint":{"browseEndpoint":{"browseId":"UCchan"}}}]}}},
	  "attributedDescription":{"content":"Site: https://ex.com/a\n😀 Clip — • Other video\n1:05 intro @someone #tag",
	   "commandRuns":[
	    {"startIndex":6,"length":16,"onTap":{"innertubeCommand":{"urlEndpoint":{"url":"https://www.youtube.com/redirect?event=video_description&q=https%3A%2F%2Fex.com%2Fa"}}}},
	    {"startIndex":33,"length":13,"onTap":{"innertubeCommand":{"watchEndpoint":{"videoId":"other","startTimeSeconds":30}}}},
	    {"startIndex":47,"length":4,"onTap":{"innertubeCommand":{"watchEndpoint":{"videoId":"me","startTimeSeconds":65}}}},
	    {"startIndex":58,"length":8,"onTap":{"innertubeCommand":{"browseEndpoint":{"browseId":"UCsome"}}}},
	    {"startIndex":67,"length":4,"onTap":{"innertubeCommand":{"commandMetadata":{"webCommandMetadata":{"url":"/hashtag/tag"}},"browseEndpoint":{"browseId":"FEhashtag"}}}}]}}}]}}}}}`)
	a := parseVideoAbout(n, "me")
	if a.Title != "A talk" || a.Views != "541 views" || a.Date != "11 Nov 2026" || a.Channel != "Chan" || a.ChannelID != "UCchan" {
		t.Errorf("the facts read %+v", a)
	}
	if len(a.Tags) != 1 || a.Tags[0].Link != "https://www.youtube.com/playlist?list=PL1" {
		t.Errorf("the tags read %+v", a.Tags)
	}
	if len(a.Description) != 3 {
		t.Fatalf("the description read %d lines: %+v", len(a.Description), a.Description)
	}
	if l := a.Description[0]; len(l) != 2 || l[1].Text != "https://ex.com/a" || l[1].Link != "https://ex.com/a" {
		t.Errorf("the first line %+v", l)
	}
	if l := a.Description[1]; len(l) != 2 || l[1].Text != "• Other video" || l[1].Link != "https://www.youtube.com/watch?t=30s&v=other" {
		t.Errorf("the second line %+v", l)
	}
	l := a.Description[2]
	if len(l) != 5 || !l[0].Seek || l[0].AtMs != 65000 || l[2].Link != "https://www.youtube.com/channel/UCsome" || l[4].Link != "https://www.youtube.com/hashtag/tag" {
		t.Errorf("the third line %+v", l)
	}
}

// TestParseVideoAboutRuns reads a description of the older design, runs
// whose endpoints link.
func TestParseVideoAboutRuns(t *testing.T) {
	n := decode(t, `{"videoSecondaryInfoRenderer":{"description":{"runs":[{"text":"See "},
	 {"text":"this","navigationEndpoint":{"watchEndpoint":{"videoId":"x"}}}]}}}`)
	a := parseVideoAbout(n, "me")
	if len(a.Description) != 1 || len(a.Description[0]) != 2 || a.Description[0][1].Link != "https://www.youtube.com/watch?v=x" {
		t.Errorf("the description read %+v", a.Description)
	}
}
