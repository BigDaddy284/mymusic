package youtube

// A channel's page of YouTube's website: its name, handle, subscribers and
// avatar, whether the account subscribes, and its videos, latest first, or
// its playlists.

import (
	"context"
	"net/url"
	"strings"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// The params of a channel's tabs: its videos, and its playlists.
const (
	channelVideosParams    = "EgZ2aWRlb3PyBgQKAjoA"
	channelPlaylistsParams = "EglwbGF5bGlzdHPyBgQKAkIA"
)

// VideoChannel is a channel of YouTube's website, and a page of its
// videos, or its playlists.
type VideoChannel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Handle string `json:"handle,omitempty"`
	// Subscribers and Count are as YouTube writes them: "1.2M
	// subscribers", "340 videos".
	Subscribers string          `json:"subscribers,omitempty"`
	Count       string          `json:"count,omitempty"`
	Avatar      []media.Image   `json:"avatar"`
	Subscribed  bool            `json:"subscribed,omitempty"`
	Videos      VideoPage       `json:"videos"`
	Playlists   []VideoPlaylist `json:"playlists,omitempty"`
}

// ChannelVideos returns the channel id, a browse ID as "UC…", with its
// latest videos.
func (c *Client) ChannelVideos(ctx context.Context, id string) (VideoChannel, error) {
	return c.channelTab(ctx, id, channelVideosParams)
}

// ChannelPlaylists returns the channel id with its playlists.
func (c *Client) ChannelPlaylists(ctx context.Context, id string) (VideoChannel, error) {
	return c.channelTab(ctx, id, channelPlaylistsParams)
}

func (c *Client) channelTab(ctx context.Context, id, params string) (VideoChannel, error) {
	p, _ := url.QueryUnescape(params)
	n, err := c.callWith(ctx, websiteClient, "browse", map[string]any{"browseId": id, "params": p}, nil, callOptions{website: true})
	if err != nil {
		return VideoChannel{}, err
	}
	ch := parseVideoChannel(n)
	ch.ID = id
	return ch, nil
}

// parseVideoChannel reads a channel's page: its header, of the older
// c4TabbedHeaderRenderer or the newer pageHeaderRenderer, and its videos
// and playlists.
func parseVideoChannel(n node) VideoChannel {
	var ch VideoChannel
	if h := n.find("c4TabbedHeaderRenderer"); h.exists() {
		ch.Name = h.get("title").str()
		ch.Handle = h.get("channelHandleText").text()
		ch.Subscribers = h.get("subscriberCountText").text()
		ch.Count = h.get("videosCountText").text()
		ch.Avatar = imagesOf(h.get("avatar", "thumbnails"))
	}
	if h := n.find("pageHeaderViewModel"); h.exists() {
		if ch.Name == "" {
			ch.Name = h.get("title", "dynamicTextViewModel", "text", "content").str()
		}
		if len(ch.Avatar) == 0 {
			ch.Avatar = imagesOf(h.find("avatarViewModel").get("image", "sources"))
		}
		for _, row := range h.find("contentMetadataViewModel").get("metadataRows").list() {
			for _, part := range row.get("metadataParts").list() {
				text := part.get("text", "content").str()
				switch {
				case strings.HasPrefix(text, "@") && ch.Handle == "":
					ch.Handle = text
				case strings.Contains(text, "subscriber") && ch.Subscribers == "":
					ch.Subscribers = text
				case strings.Contains(text, "video") && ch.Count == "":
					ch.Count = text
				}
			}
		}
	}
	if ch.Name == "" {
		ch.Name = n.find("channelMetadataRenderer").get("title").str()
	}
	if len(ch.Avatar) == 0 {
		ch.Avatar = imagesOf(n.find("channelMetadataRenderer").get("avatar", "thumbnails"))
	}
	if s := n.find("subscriptionStateEntity"); s.exists() {
		ch.Subscribed = s.get("subscribed").boolean()
	} else if b := n.find("subscribeButtonRenderer"); b.exists() {
		ch.Subscribed = b.get("subscribed").boolean()
	}
	videos, lists := collectVideos(n)
	if videos == nil {
		videos = []Video{}
	}
	for i := range videos {
		if videos[i].Channel == "" {
			videos[i].Channel = ch.Name
		}
	}
	ch.Videos = VideoPage{Videos: videos, Continue: nextVideosToken(n)}
	ch.Playlists = lists
	return ch
}
