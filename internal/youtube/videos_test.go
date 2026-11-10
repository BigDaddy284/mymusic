package youtube

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) node {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatal(err)
	}
	return node{v}
}

// TestParseVideoPage reads videos of both designs YouTube's website
// serves, marks Shorts, and finds the token of the next page.
func TestParseVideoPage(t *testing.T) {
	n := decode(t, `{"contents":{"sectionListRenderer":{"contents":[
	 {"itemSectionRenderer":{"contents":[
	  {"videoRenderer":{"videoId":"vid00000001","title":{"runs":[{"text":"Old design"}]},
	   "ownerText":{"runs":[{"text":"A Channel","navigationEndpoint":{"browseEndpoint":{"browseId":"UC1"}}}]},
	   "lengthText":{"simpleText":"3:12"},"shortViewCountText":{"simpleText":"12K views"},
	   "publishedTimeText":{"simpleText":"2 years ago"},
	   "thumbnail":{"thumbnails":[{"url":"https://i.ytimg.com/vi/vid00000001/hq720.jpg","width":720,"height":404}]},
	   "thumbnailOverlays":[{"thumbnailOverlayResumePlaybackRenderer":{"percentDurationWatched":37}}]}},
	  {"reelShelfRenderer":{"items":[{"videoRenderer":{"videoId":"short000001","title":{"runs":[{"text":"A Short"}]}}}]}},
	  {"lockupViewModel":{"contentId":"vid00000002","contentType":"LOCKUP_CONTENT_TYPE_VIDEO",
	   "contentImage":{"thumbnailViewModel":{"image":{"sources":[{"url":"https://i.ytimg.com/vi/vid00000002/hq.jpg","width":480,"height":270}]},
	    "overlays":[{"thumbnailOverlayBadgeViewModel":{"thumbnailBadges":[{"thumbnailBadgeViewModel":{"text":"14:30"}}]}},
	     {"thumbnailBottomOverlayViewModel":{"progressBar":{"thumbnailOverlayProgressBarViewModel":{"startPercent":100}}}}]}},
	   "metadata":{"lockupMetadataViewModel":{"title":{"content":"New design"},
	    "metadata":{"contentMetadataViewModel":{"metadataRows":[
	     {"metadataParts":[{"text":{"content":"Another Channel"}}]},
	     {"metadataParts":[{"text":{"content":"48K views"}},{"text":{"content":"4 years ago"}}]}]}}}}}},
	  {"continuationItemRenderer":{"continuationEndpoint":{"continuationCommand":{"token":"NEXT"}}}}
	 ]}}]}}}`)
	page := parseVideoPage(n)
	if len(page.Videos) != 3 {
		t.Fatalf("read %d videos: %+v", len(page.Videos), page.Videos)
	}
	old, short, lockup := page.Videos[0], page.Videos[1], page.Videos[2]
	if old.Short || !short.Short || short.ID != "short000001" || lockup.Short {
		t.Errorf("Shorts marked: %v, %v, %v", old.Short, short.Short, lockup.Short)
	}
	if old.ID != "vid00000001" || old.Title != "Old design" || old.Channel != "A Channel" || old.ChannelID != "UC1" ||
		old.Duration != "3:12" || old.Views != "12K views" || old.Published != "2 years ago" || len(old.Thumbnail) != 1 {
		t.Errorf("the old design reads %+v", old)
	}
	if lockup.ID != "vid00000002" || lockup.Title != "New design" || lockup.Channel != "Another Channel" ||
		lockup.Duration != "14:30" || lockup.Views != "48K views" || lockup.Published != "4 years ago" {
		t.Errorf("the new design reads %+v", lockup)
	}
	if old.Watched != 37 || lockup.Watched != 100 {
		t.Errorf("watched %d%% and %d%%", old.Watched, lockup.Watched)
	}
	if page.Continue != "NEXT" {
		t.Errorf("the next page's token is %q", page.Continue)
	}
}

// TestParsePlaylists reads playlists of the lockup design.
func TestParsePlaylists(t *testing.T) {
	n := decode(t, `{"items":[{"lockupViewModel":{"contentId":"PL123","contentType":"LOCKUP_CONTENT_TYPE_PLAYLIST",
	 "contentImage":{"collectionThumbnailViewModel":{"primaryThumbnail":{"thumbnailViewModel":{"image":{"sources":[{"url":"https://i.ytimg.com/x.jpg","width":480,"height":270}]},
	  "overlays":[{"thumbnailOverlayBadgeViewModel":{"thumbnailBadges":[{"thumbnailBadgeViewModel":{"text":"24 videos"}}]}}]}}}},
	 "metadata":{"lockupMetadataViewModel":{"title":{"content":"Road trip"}}}}}]}`)
	_, lists := collectVideos(n)
	if len(lists) != 1 || lists[0].ID != "PL123" || lists[0].Title != "Road trip" || lists[0].Count != "24 videos" || len(lists[0].Thumbnail) != 1 {
		t.Errorf("the playlists read %+v", lists)
	}
}

// TestVideoThumbnailFallback gives every video a thumbnail of i.ytimg.com
// when the response tells none.
func TestVideoThumbnailFallback(t *testing.T) {
	th := videoThumbnail("abc", node{})
	if len(th) == 0 || th[0].URL != "https://i.ytimg.com/vi/abc/mqdefault.jpg" {
		t.Errorf("the fallback thumbnails are %+v", th)
	}
}

// TestParseVideoDetails reads the views and the account's rating of a
// watch page.
func TestParseVideoDetails(t *testing.T) {
	n := decode(t, `{"contents":{"twoColumnWatchNextResults":{"results":{"results":{"contents":[
	 {"videoPrimaryInfoRenderer":{"viewCount":{"videoViewCountRenderer":{"viewCount":{"simpleText":"12,345 views"},"shortViewCount":{"simpleText":"12K views"}}}}}]}}}},
	 "frameworkUpdates":{"entityBatchUpdate":{"mutations":[{"payload":{"likeStatusEntity":{"key":"k","likeStatus":"LIKE"}}}]}}}`)
	if d := parseVideoDetails(n); d.Views != "12K views" || d.Rating != "like" {
		t.Errorf("details %+v", d)
	}
}

// TestParseChapters reads a watch page's chapters from its seek bar.
func TestParseChapters(t *testing.T) {
	n := decode(t, `{"playerOverlays":{"playerOverlayRenderer":{"decoratedPlayerBarRenderer":{"decoratedPlayerBarRenderer":{"playerBar":{"multiMarkersPlayerBarRenderer":{"markersMap":[{"key":"DESCRIPTION_CHAPTERS","value":{"chapters":[
	 {"chapterRenderer":{"title":{"simpleText":"Intro"},"timeRangeStartMillis":0}},
	 {"chapterRenderer":{"title":{"simpleText":"The build"},"timeRangeStartMillis":95000}}]}}]}}}}}},
	 "engagementPanels":[{"macroMarkersListItemRenderer":{"title":{"simpleText":"Other"},"onTap":{"watchEndpoint":{"startTimeSeconds":3}}}}]}`)
	got := parseChapters(n)
	if len(got) != 2 || got[0].Title != "Intro" || got[1].Start != 95 {
		t.Errorf("chapters %+v", got)
	}
	panel := decode(t, `{"engagementPanels":[{"macroMarkersListItemRenderer":{"title":{"simpleText":"Start"},"onTap":{"watchEndpoint":{"startTimeSeconds":0}}}},
	 {"macroMarkersListItemRenderer":{"title":{"simpleText":"End"},"onTap":{"watchEndpoint":{"startTimeSeconds":60}}}}]}`)
	if got := parseChapters(panel); len(got) != 2 || got[1].Title != "End" {
		t.Errorf("panel chapters %+v", got)
	}
}

// TestParseVideoFeedback reads the tokens of Not interested and Don't
// recommend channel of a video's menu, of both designs.
func TestParseVideoFeedback(t *testing.T) {
	old := decode(t, `{"menuRenderer":{"items":[
	 {"menuServiceItemRenderer":{"text":{"runs":[{"text":"Add to queue"}]},"icon":{"iconType":"ADD_TO_QUEUE_TAIL"}}},
	 {"menuServiceItemRenderer":{"text":{"runs":[{"text":"Not interested"}]},"icon":{"iconType":"NOT_INTERESTED"},
	  "serviceEndpoint":{"feedbackEndpoint":{"feedbackToken":"NI"}}}},
	 {"menuServiceItemRenderer":{"text":{"runs":[{"text":"Don't recommend channel"}]},"icon":{"iconType":"REMOVE"},
	  "serviceEndpoint":{"feedbackEndpoint":{"feedbackToken":"NC"}}}}]}}`)
	if f := parseVideoFeedback(old); f.NotInterested != "NI" || f.NoChannel != "NC" {
		t.Errorf("the old design's feedback %+v", f)
	}
	lockup := decode(t, `{"buttonViewModel":{"onTap":{"innertubeCommand":{"showSheetCommand":{"panelLoadingStrategy":{"inlineContent":{"sheetViewModel":{"content":{"listViewModel":{"listItems":[
	 {"listItemViewModel":{"title":{"content":"Not interested"},"leadingImage":{"sources":[{"clientResource":{"imageName":"NOT_INTERESTED"}}]},
	  "rendererContext":{"commandContext":{"onTap":{"innertubeCommand":{"feedbackEndpoint":{"feedbackToken":"NI2"}}}}}}},
	 {"listItemViewModel":{"title":{"content":"Don't recommend channel"},"leadingImage":{"sources":[{"clientResource":{"imageName":"REMOVE"}}]},
	  "rendererContext":{"commandContext":{"onTap":{"innertubeCommand":{"feedbackEndpoint":{"feedbackToken":"NC2"}}}}}}}]}}}}}}}}}}`)
	if f := parseVideoFeedback(lockup); f.NotInterested != "NI2" || f.NoChannel != "NC2" {
		t.Errorf("the lockup's feedback %+v", f)
	}
}

// TestParseVideoChannel reads a channel's header, of the newer design,
// and its videos.
func TestParseVideoChannel(t *testing.T) {
	n := decode(t, `{"header":{"pageHeaderRenderer":{"content":{"pageHeaderViewModel":{
	  "title":{"dynamicTextViewModel":{"text":{"content":"Some Channel"}}},
	  "image":{"decoratedAvatarViewModel":{"avatar":{"avatarViewModel":{"image":{"sources":[{"url":"https://yt3/a.jpg","width":72,"height":72}]}}}}},
	  "metadata":{"contentMetadataViewModel":{"metadataRows":[{"metadataParts":[{"text":{"content":"@some"}}]},
	   {"metadataParts":[{"text":{"content":"1.2M subscribers"}},{"text":{"content":"340 videos"}}]}]}}}}}},
	 "contents":{"richGridRenderer":{"contents":[{"richItemRenderer":{"content":{"videoRenderer":{"videoId":"vid00000009","title":{"runs":[{"text":"Latest"}]}}}}}]}},
	 "frameworkUpdates":{"entityBatchUpdate":{"mutations":[{"payload":{"subscriptionStateEntity":{"subscribed":true}}}]}}}`)
	ch := parseVideoChannel(n)
	if ch.Name != "Some Channel" || ch.Handle != "@some" || ch.Subscribers != "1.2M subscribers" || ch.Count != "340 videos" ||
		len(ch.Avatar) != 1 || !ch.Subscribed || len(ch.Videos.Videos) != 1 || ch.Videos.Videos[0].Channel != "Some Channel" {
		t.Errorf("channel %+v", ch)
	}
}
