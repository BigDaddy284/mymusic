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
// serves, leaves out Shorts, and finds the token of the next page.
func TestParseVideoPage(t *testing.T) {
	n := decode(t, `{"contents":{"sectionListRenderer":{"contents":[
	 {"itemSectionRenderer":{"contents":[
	  {"videoRenderer":{"videoId":"vid00000001","title":{"runs":[{"text":"Old design"}]},
	   "ownerText":{"runs":[{"text":"A Channel","navigationEndpoint":{"browseEndpoint":{"browseId":"UC1"}}}]},
	   "lengthText":{"simpleText":"3:12"},"shortViewCountText":{"simpleText":"12K views"},
	   "publishedTimeText":{"simpleText":"2 years ago"},
	   "thumbnail":{"thumbnails":[{"url":"https://i.ytimg.com/vi/vid00000001/hq720.jpg","width":720,"height":404}]}}},
	  {"reelShelfRenderer":{"items":[{"videoRenderer":{"videoId":"short000001","title":{"runs":[{"text":"A Short"}]}}}]}},
	  {"lockupViewModel":{"contentId":"vid00000002","contentType":"LOCKUP_CONTENT_TYPE_VIDEO",
	   "contentImage":{"thumbnailViewModel":{"image":{"sources":[{"url":"https://i.ytimg.com/vi/vid00000002/hq.jpg","width":480,"height":270}]},
	    "overlays":[{"thumbnailOverlayBadgeViewModel":{"thumbnailBadges":[{"thumbnailBadgeViewModel":{"text":"14:30"}}]}}]}},
	   "metadata":{"lockupMetadataViewModel":{"title":{"content":"New design"},
	    "metadata":{"contentMetadataViewModel":{"metadataRows":[
	     {"metadataParts":[{"text":{"content":"Another Channel"}}]},
	     {"metadataParts":[{"text":{"content":"48K views"}},{"text":{"content":"4 years ago"}}]}]}}}}}},
	  {"continuationItemRenderer":{"continuationEndpoint":{"continuationCommand":{"token":"NEXT"}}}}
	 ]}}]}}}`)
	page := parseVideoPage(n)
	if len(page.Videos) != 2 {
		t.Fatalf("read %d videos: %+v", len(page.Videos), page.Videos)
	}
	old, lockup := page.Videos[0], page.Videos[1]
	if old.ID != "vid00000001" || old.Title != "Old design" || old.Channel != "A Channel" || old.ChannelID != "UC1" ||
		old.Duration != "3:12" || old.Views != "12K views" || old.Published != "2 years ago" || len(old.Thumbnail) != 1 {
		t.Errorf("the old design reads %+v", old)
	}
	if lockup.ID != "vid00000002" || lockup.Title != "New design" || lockup.Channel != "Another Channel" ||
		lockup.Duration != "14:30" || lockup.Views != "48K views" || lockup.Published != "4 years ago" {
		t.Errorf("the new design reads %+v", lockup)
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
