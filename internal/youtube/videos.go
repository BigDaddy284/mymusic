package youtube

// The Video page reads YouTube's own website, www.youtube.com, with the
// session YouTube Music signs in with: Google signs both in with the same
// cookies. It reads search results, the subscriptions feed, the watch
// history, the account's playlists and their videos, as YouTube's website
// shows them, in the renderers of either of the two designs YouTube serves
// now: the older videoRenderer family, and the lockupViewModel.

import (
	"context"
	"net/url"
	"slices"
	"strings"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// Video is a video of YouTube's website.
type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Channel names the channel that uploaded it, and ChannelID is its
	// browse ID, when the result tells it.
	Channel   string `json:"channel,omitempty"`
	ChannelID string `json:"channelId,omitempty"`
	// Duration is as YouTube writes it, as "12:34", or "LIVE".
	Duration string `json:"duration,omitempty"`
	// Views and Published are as YouTube writes them: "1.2M views", "2
	// days ago".
	Views     string        `json:"views,omitempty"`
	Published string        `json:"published,omitempty"`
	Thumbnail []media.Image `json:"thumbnail"`
}

// VideoPlaylist is a playlist of YouTube's website.
type VideoPlaylist struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Count is as YouTube writes it, as "24 videos".
	Count     string        `json:"count,omitempty"`
	Owner     string        `json:"owner,omitempty"`
	Thumbnail []media.Image `json:"thumbnail"`
}

// VideoPage is a page of videos, and the token of the next, if any.
type VideoPage struct {
	Videos   []Video `json:"videos"`
	Continue string  `json:"continue,omitempty"`
}

// websiteClient is YouTube's website, which the Video page reads as.
var websiteClient = webClient

// searchVideosParams filters a search of YouTube's website to videos.
const searchVideosParams = "EgIQAQ%3D%3D"

// searchThisWeekParams filters a search of YouTube's website to videos
// uploaded this week.
const searchThisWeekParams = "EgQIAxAB"

// SearchVideos searches YouTube's website for videos.
func (c *Client) SearchVideos(ctx context.Context, query string) (VideoPage, error) {
	return c.searchVideos(ctx, query, searchVideosParams)
}

// SearchRecentVideos searches YouTube's website for videos uploaded this
// week.
func (c *Client) SearchRecentVideos(ctx context.Context, query string) (VideoPage, error) {
	return c.searchVideos(ctx, query, searchThisWeekParams)
}

// HomeVideos returns the videos YouTube's home page recommends: the
// account's own, signed in.
func (c *Client) HomeVideos(ctx context.Context) (VideoPage, error) {
	n, err := c.callWith(ctx, websiteClient, "browse", map[string]any{"browseId": "FEwhat_to_watch"}, nil, callOptions{website: true})
	if err != nil {
		return VideoPage{}, err
	}
	return parseVideoPage(n), nil
}

// searchVideos searches YouTube's website, filtered by params.
func (c *Client) searchVideos(ctx context.Context, query, filter string) (VideoPage, error) {
	params, _ := url.QueryUnescape(filter)
	n, err := c.callWith(ctx, websiteClient, "search", map[string]any{"query": query, "params": params}, nil, callOptions{website: true})
	if err != nil {
		return VideoPage{}, err
	}
	return parseVideoPage(n), nil
}

// MoreVideos returns the page of videos after one of search's or browse's,
// by its continuation token.
func (c *Client) MoreVideos(ctx context.Context, endpoint, token string) (VideoPage, error) {
	if endpoint != "search" {
		endpoint = "browse"
	}
	n, err := c.callWith(ctx, websiteClient, endpoint, map[string]any{"continuation": token}, nil, callOptions{website: true})
	if err != nil {
		return VideoPage{}, err
	}
	return parseVideoPage(n), nil
}

// Subscriptions returns the latest videos of the channels the account
// subscribes to.
func (c *Client) Subscriptions(ctx context.Context) (VideoPage, error) {
	return c.browseVideos(ctx, "FEsubscriptions")
}

// WatchHistory returns the videos the account watched, the latest first.
func (c *Client) WatchHistory(ctx context.Context) (VideoPage, error) {
	return c.browseVideos(ctx, "FEhistory")
}

// PlaylistVideos returns the videos of a playlist of YouTube's website, as
// "WL" for Watch later and "LL" for liked videos.
func (c *Client) PlaylistVideos(ctx context.Context, id string) (VideoPage, error) {
	return c.browseVideos(ctx, "VL"+id)
}

// VideoPlaylists returns the playlists of the account on YouTube's
// website: those it made and those it saved.
func (c *Client) VideoPlaylists(ctx context.Context) ([]VideoPlaylist, error) {
	if err := c.session(); err != nil {
		return nil, err
	}
	var lists []VideoPlaylist
	for _, id := range []string{"FEplaylist_aggregation", "FElibrary"} {
		n, err := c.callWith(ctx, websiteClient, "browse", map[string]any{"browseId": id}, nil, callOptions{website: true})
		if err != nil {
			if len(lists) > 0 {
				break
			}
			return nil, err
		}
		_, lists = collectVideos(n)
		if len(lists) > 0 {
			break
		}
	}
	return lists, nil
}

// browseVideos browses a page of YouTube's website for its videos, signed
// in.
func (c *Client) browseVideos(ctx context.Context, browseID string) (VideoPage, error) {
	if err := c.session(); err != nil {
		return VideoPage{}, err
	}
	n, err := c.callWith(ctx, websiteClient, "browse", map[string]any{"browseId": browseID}, nil, callOptions{website: true})
	if err != nil {
		return VideoPage{}, err
	}
	return parseVideoPage(n), nil
}

// parseVideoPage reads the videos of a response of YouTube's website, and
// the token of the page after.
func parseVideoPage(n node) VideoPage {
	videos, _ := collectVideos(n)
	if videos == nil {
		videos = []Video{}
	}
	return VideoPage{Videos: videos, Continue: nextVideosToken(n)}
}

// collectVideos walks a response for the videos and playlists it shows,
// in order, each once. Shorts, ads and the like are left out.
func collectVideos(n node) ([]Video, []VideoPlaylist) {
	var videos []Video
	var lists []VideoPlaylist
	seenV, seenP := map[string]bool{}, map[string]bool{}
	addVideo := func(v Video, ok bool) {
		if ok && v.ID != "" && !seenV[v.ID] {
			seenV[v.ID] = true
			videos = append(videos, v)
		}
	}
	addList := func(p VideoPlaylist, ok bool) {
		if ok && p.ID != "" && !seenP[p.ID] {
			seenP[p.ID] = true
			lists = append(lists, p)
		}
	}
	var walk func(v any, depth int)
	walk = func(v any, depth int) {
		if depth > 64 {
			return
		}
		switch v := v.(type) {
		case []any:
			for _, x := range v {
				walk(x, depth+1)
			}
		case map[string]any:
			for _, key := range sortedKeys(v) {
				x := node{v[key]}
				switch key {
				case "videoRenderer", "gridVideoRenderer", "compactVideoRenderer", "playlistVideoRenderer":
					addVideo(parseVideoRenderer(x))
				case "playlistRenderer", "gridPlaylistRenderer":
					addList(parsePlaylistRenderer(x))
				case "lockupViewModel":
					switch x.get("contentType").str() {
					case "LOCKUP_CONTENT_TYPE_VIDEO":
						addVideo(parseVideoLockup(x))
					case "LOCKUP_CONTENT_TYPE_PLAYLIST", "LOCKUP_CONTENT_TYPE_PODCAST":
						addList(parsePlaylistLockup(x))
					}
				case "reelShelfRenderer", "richShelfRenderer", "shortsLockupViewModel", "reelItemRenderer", "adSlotRenderer":
					// Shorts and ads: left out.
				default:
					walk(v[key], depth+1)
				}
			}
		}
	}
	walk(n.v, 0)
	return videos, lists
}

// sortedKeys returns the keys of m in an order that keeps the renderers of
// a response in the order they come: JSON objects of InnerTube responses
// hold one renderer each, or a renderer's fields, so any fixed order will
// do, and a fixed one keeps the results the same each time.
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// parseVideoRenderer reads a video of the videoRenderer family.
func parseVideoRenderer(r node) (Video, bool) {
	id := r.get("videoId").str()
	if id == "" || r.get("upcomingEventData").exists() {
		return Video{}, false
	}
	v := Video{ID: id, Title: r.get("title").text()}
	if v.Title == "" {
		v.Title = r.get("title", "accessibility", "accessibilityData", "label").str()
	}
	for _, path := range []string{"ownerText", "longBylineText", "shortBylineText"} {
		if runs := r.get(path).runs(); len(runs) > 0 {
			v.Channel = runs[0].get("text").str()
			v.ChannelID = runs[0].get("navigationEndpoint", "browseEndpoint", "browseId").str()
			break
		}
	}
	v.Duration = r.get("lengthText").text()
	if v.Duration == "" {
		if o := r.find("thumbnailOverlayTimeStatusRenderer"); o.exists() {
			v.Duration = o.get("text").text()
		}
	}
	if v.Duration == "" && hasLiveBadge(r) {
		v.Duration = "LIVE"
	}
	v.Views = r.get("shortViewCountText").text()
	if v.Views == "" {
		v.Views = r.get("viewCountText").text()
	}
	v.Published = r.get("publishedTimeText").text()
	if v.Views == "" && v.Published == "" {
		// A playlist's video tells them in its info: "1.2M views • 2 years ago".
		parts := strings.Split(r.get("videoInfo").text(), " • ")
		if len(parts) > 0 {
			v.Views = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			v.Published = strings.TrimSpace(parts[1])
		}
	}
	v.Thumbnail = videoThumbnail(id, r.get("thumbnail", "thumbnails"))
	return v, v.Title != ""
}

// hasLiveBadge reports whether a video's badges say it is live now.
func hasLiveBadge(r node) bool {
	for _, b := range r.get("badges").list() {
		if b.get("metadataBadgeRenderer", "style").str() == "BADGE_STYLE_TYPE_LIVE_NOW" {
			return true
		}
	}
	return false
}

// parsePlaylistRenderer reads a playlist of the playlistRenderer family.
func parsePlaylistRenderer(r node) (VideoPlaylist, bool) {
	p := VideoPlaylist{ID: r.get("playlistId").str(), Title: r.get("title").text()}
	if c := r.get("videoCountText").text(); c != "" {
		p.Count = c
	} else if c := r.get("videoCount").str(); c != "" {
		p.Count = c + " videos"
	}
	if runs := r.get("shortBylineText").runs(); len(runs) > 0 {
		p.Owner = runs[0].get("text").str()
	}
	thumbs := r.get("thumbnails", 0, "thumbnails")
	if !thumbs.exists() {
		thumbs = r.get("thumbnail", "thumbnails")
	}
	p.Thumbnail = imagesOf(thumbs)
	return p, p.Title != ""
}

// parseVideoLockup reads a video of the lockupViewModel design.
func parseVideoLockup(r node) (Video, bool) {
	id := r.get("contentId").str()
	if id == "" {
		id = r.find("watchEndpoint").get("videoId").str()
	}
	meta := r.get("metadata", "lockupMetadataViewModel")
	v := Video{ID: id, Title: meta.get("title", "content").str()}
	rows := meta.get("metadata", "contentMetadataViewModel", "metadataRows").list()
	var parts []string
	for i, row := range rows {
		for _, part := range row.get("metadataParts").list() {
			text := part.get("text", "content").str()
			if text == "" {
				continue
			}
			if i == 0 && v.Channel == "" && len(rows) > 1 {
				v.Channel = text
				v.ChannelID = part.find("browseEndpoint").get("browseId").str()
				continue
			}
			parts = append(parts, text)
		}
	}
	for _, p := range parts {
		switch {
		case v.Views == "" && (strings.Contains(p, "view") || strings.Contains(p, "watching")):
			v.Views = p
		case v.Published == "" && (strings.Contains(p, "ago") || strings.Contains(p, "Streamed") || strings.Contains(p, "Premiered")):
			v.Published = p
		}
	}
	image := r.get("contentImage")
	v.Duration = badgeText(image)
	v.Thumbnail = videoThumbnail(id, image.find("image").get("sources"))
	return v, id != "" && v.Title != ""
}

// parsePlaylistLockup reads a playlist of the lockupViewModel design.
func parsePlaylistLockup(r node) (VideoPlaylist, bool) {
	id := r.get("contentId").str()
	meta := r.get("metadata", "lockupMetadataViewModel")
	p := VideoPlaylist{ID: id, Title: meta.get("title", "content").str()}
	image := r.get("contentImage")
	p.Count = badgeText(image)
	for _, row := range meta.get("metadata", "contentMetadataViewModel", "metadataRows").list() {
		for _, part := range row.get("metadataParts").list() {
			if text := part.get("text", "content").str(); text != "" && p.Owner == "" && !strings.Contains(text, "View full") {
				p.Owner = text
			}
		}
	}
	p.Thumbnail = imagesOf(image.find("image").get("sources"))
	return p, id != "" && p.Title != ""
}

// badgeText returns the text of the first badge over a lockup's image: a
// video's duration, or a playlist's count of videos.
func badgeText(image node) string {
	return image.find("thumbnailBadgeViewModel").get("text").str()
}

// imagesOf returns the images of a thumbnail list, or of an image's
// sources: each with its URL made absolute.
func imagesOf(list node) []media.Image {
	var out []media.Image
	for _, t := range list.list() {
		u := t.get("url").str()
		if u == "" {
			continue
		}
		if strings.HasPrefix(u, "//") {
			u = "https:" + u
		}
		out = append(out, media.Image{URL: u, Width: t.get("width").number(), Height: t.get("height").number()})
	}
	if out == nil {
		out = []media.Image{}
	}
	return out
}

// videoThumbnail returns the thumbnails of a video: those the response
// tells, and the ones of i.ytimg.com that every video has, so that one
// always shows.
func videoThumbnail(id string, list node) []media.Image {
	out := imagesOf(list)
	if len(out) == 0 {
		out = append(out,
			media.Image{URL: "https://i.ytimg.com/vi/" + id + "/mqdefault.jpg", Width: 320, Height: 180},
			media.Image{URL: "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg", Width: 480, Height: 360},
		)
	}
	return out
}

// nextVideosToken finds the token of the page after a response's: the
// last continuation item's, as the one of the list sits at its end.
func nextVideosToken(n node) string {
	token := ""
	var walk func(v any, depth int)
	walk = func(v any, depth int) {
		if depth > 64 {
			return
		}
		switch v := v.(type) {
		case []any:
			for _, x := range v {
				walk(x, depth+1)
			}
		case map[string]any:
			if r, ok := v["continuationItemRenderer"]; ok {
				if t := (node{r}).find("continuationCommand").get("token").str(); t != "" {
					token = t
				}
				return
			}
			for _, key := range sortedKeys(v) {
				walk(v[key], depth+1)
			}
		}
	}
	walk(n.v, 0)
	return token
}
