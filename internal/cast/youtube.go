package cast

// Casting a video of YouTube to a device's own YouTube app: MyTube starts
// the app, asks it for its screen ID, and then controls it as YouTube's
// phone app does, through YouTube's lounge (lounge.go). The device plays
// the video itself, at its own quality.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// YouTubeApp is the ID of YouTube's app on cast devices.
const YouTubeApp = "233637DE"

// nsYouTube is the namespace YouTube's app speaks to senders in.
const nsYouTube = "urn:x-cast:com.google.youtube.mdx"

// youTubeWait bounds how long the device takes to start YouTube and tell
// its screen.
var youTubeWait = 20 * time.Second

// YouTubeScreen starts YouTube's app on d, unless it runs, and returns the
// screen ID it tells, by which the lounge reaches it.
func YouTubeScreen(ctx context.Context, d Device) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, youTubeWait)
	defer cancel()
	screens := make(chan string, 1)
	c, err := dial(ctx, d.Addr(), func(m message, _ header) {
		if m.Namespace != nsYouTube {
			return
		}
		var msg struct {
			Type string `json:"type"`
			Data struct {
				ScreenID string `json:"screenId"`
			} `json:"data"`
		}
		if json.Unmarshal(m.Payload, &msg) == nil && msg.Type == "mdxSessionStatus" && msg.Data.ScreenID != "" {
			select {
			case screens <- msg.Data.ScreenID:
			default:
			}
		}
	})
	if err != nil {
		return "", err
	}
	defer c.close(ErrClosed)
	if err := c.connect(receiverID); err != nil {
		return "", err
	}
	app, err := launchApp(ctx, c, YouTubeApp)
	if err != nil {
		return "", err
	}
	if err := c.connect(app.TransportID); err != nil {
		return "", err
	}
	ask := func() error {
		return c.send(nsYouTube, app.TransportID, map[string]any{"type": "getMdxSessionStatus"})
	}
	if err := ask(); err != nil {
		return "", err
	}
	// The app tells its screen once it's ready: it's asked again until it
	// does.
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case id := <-screens:
			return id, nil
		case <-ctx.Done():
			return "", errors.New("YouTube on the device didn't answer in time")
		case <-c.closed:
			return "", c.reason()
		case <-t.C:
			if err := ask(); err != nil {
				return "", err
			}
		}
	}
}

// launchApp launches the app appID on the device of c, unless it runs,
// and returns it once it does.
func launchApp(ctx context.Context, c *conn, appID string) (*application, error) {
	status := func() (receiverStatus, error) {
		r, err := c.request(ctx, nsReceiver, receiverID, map[string]any{"type": "GET_STATUS"})
		if err != nil {
			return receiverStatus{}, err
		}
		return parseReceiverStatus(r.msg.Payload)
	}
	st, err := status()
	if err != nil {
		return nil, err
	}
	if a := st.app(appID); a != nil {
		return a, nil
	}
	r, err := c.request(ctx, nsReceiver, receiverID, map[string]any{"type": "LAUNCH", "appId": appID})
	if err != nil {
		return nil, fmt.Errorf("the device didn't start YouTube: %w", err)
	}
	if st, err = parseReceiverStatus(r.msg.Payload); err == nil {
		if a := st.app(appID); a != nil {
			return a, nil
		}
	}
	t := time.NewTicker(launchPoll)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("the device didn't start YouTube in time")
		case <-t.C:
		}
		if st, err = status(); err == nil {
			if a := st.app(appID); a != nil {
				return a, nil
			}
		}
	}
}
