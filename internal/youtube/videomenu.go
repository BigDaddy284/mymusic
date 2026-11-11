package youtube

// The menu of a video on YouTube's website: Not interested, Don't
// recommend channel, and saving the video to Watch Later or a playlist.

import (
	"context"
	"errors"
	"strings"
)

// VideoFeedback are the tokens a video of a feed carries in its menu, for
// telling YouTube what not to recommend. Either may be empty: search
// results carry none.
type VideoFeedback struct {
	NotInterested string `json:"notInterested,omitempty"`
	NoChannel     string `json:"noChannel,omitempty"`
}

// parseVideoFeedback reads the feedback tokens of the menu of a video of
// either design: menuServiceItemRenderer items of the old, listItemViewModel
// items of the lockup's sheet. It tells the items apart by their icons,
// and their words in English.
func parseVideoFeedback(r node) VideoFeedback {
	var f VideoFeedback
	visit(r.v, maxFindDepth, func(key string, item node) {
		if key != "menuServiceItemRenderer" && key != "listItemViewModel" {
			return
		}
		token := item.find("feedbackEndpoint").get("feedbackToken").str()
		if token == "" {
			return
		}
		icon := item.get("icon", "iconType").str()
		if icon == "" {
			icon = item.find("imageName").str()
		}
		label := strings.ToLower(item.get("text").text() + item.get("title", "content").str())
		switch {
		case icon == "NOT_INTERESTED" || strings.Contains(label, "not interested"):
			f.NotInterested = token
		case icon == "REMOVE" || strings.Contains(label, "recommend channel"):
			f.NoChannel = token
		}
	})
	return f
}

// visit calls fn with every key and value of the objects under v, to depth.
func visit(v any, depth int, fn func(key string, n node)) {
	if depth <= 0 {
		return
	}
	switch v := v.(type) {
	case map[string]any:
		for _, k := range sortedKeys(v) {
			fn(k, node{v[k]})
			visit(v[k], depth-1, fn)
		}
	case []any:
		for _, x := range v {
			visit(x, depth-1, fn)
		}
	}
}

// SendVideoFeedback sends YouTube the feedback of a token of a video's
// menu, as Not interested.
func (c *Client) SendVideoFeedback(ctx context.Context, token string) error {
	if err := c.session(); err != nil {
		return err
	}
	if token == "" {
		return errors.New("YouTube doesn't take that for this video")
	}
	n, err := c.callWith(ctx, websiteClient, "feedback", map[string]any{"feedbackTokens": []string{token}}, nil, callOptions{website: true})
	if err != nil {
		return err
	}
	if done := n.get("feedbackResponses", 0, "isProcessed"); done.exists() && !done.boolean() {
		return errors.New("YouTube didn't take it")
	}
	return nil
}

// WatchLater is the playlist of the account's Watch Later.
const WatchLater = "WL"

// SaveVideo adds the video id to the account's playlist, as WatchLater.
func (c *Client) SaveVideo(ctx context.Context, playlistID, id string) error {
	if err := c.session(); err != nil {
		return err
	}
	playlistID = strings.TrimPrefix(playlistID, "VL")
	action := map[string]any{"action": "ACTION_ADD_VIDEO", "addedVideoId": id}
	n, err := c.callWith(ctx, websiteClient, "browse/edit_playlist", map[string]any{"playlistId": playlistID, "actions": []any{action}}, nil, callOptions{website: true})
	if err != nil {
		return c.refusedOrEnded(ctx, err, errors.New("YouTube refused: the playlist may be another account's"))
	}
	if !strings.Contains(n.get("status").str(), "SUCCEEDED") {
		return refusedEdit(n, errors.New("YouTube didn't save the video"))
	}
	return nil
}

// UnsaveVideo takes the video id off the account's playlist, as
// WatchLater, by its ID, as YouTube's website does from a video's menu.
func (c *Client) UnsaveVideo(ctx context.Context, playlistID, id string) error {
	if err := c.session(); err != nil {
		return err
	}
	playlistID = strings.TrimPrefix(playlistID, "VL")
	action := map[string]any{"action": "ACTION_REMOVE_VIDEO_BY_VIDEO_ID", "removedVideoId": id}
	n, err := c.callWith(ctx, websiteClient, "browse/edit_playlist", map[string]any{"playlistId": playlistID, "actions": []any{action}}, nil, callOptions{website: true})
	if err != nil {
		return c.refusedOrEnded(ctx, err, errors.New("YouTube refused: the playlist may be another account's"))
	}
	if !strings.Contains(n.get("status").str(), "SUCCEEDED") {
		return refusedEdit(n, errors.New("YouTube didn't take the video off"))
	}
	return nil
}
