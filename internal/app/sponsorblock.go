package app

// SponsorBlock: the parts of videos its users marked, as a sponsor's read
// or an intro, which the Video side's player skips, mutes or shows on the
// seek bar, as Settings › Videos asks. The app asks for them by the first
// characters of a hash of the video's ID, as SponsorBlock lets apps do,
// so SponsorBlock doesn't learn which video plays: it answers with the
// segments of every video whose hash starts the same, and the app keeps
// the video's own.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/egoist/mygo/netproxy"

	"github.com/bigdaddy284/mymusic/internal/settings"
)

// sponsorBlockAPI is SponsorBlock's API, replaced in tests.
var sponsorBlockAPI = "https://sponsor.ajay.app/api"

var sponsorBlockClient = &http.Client{Transport: netproxy.Transport(), Timeout: 10 * time.Second}

// Segment is a part of a video SponsorBlock's users marked.
type Segment struct {
	// Start and End are in seconds; a highlight's are the same.
	Start    float64 `json:"start"`
	End      float64 `json:"end"`
	Category string  `json:"category"`
	// Mute is set for a segment to be muted, not skipped, as SponsorBlock
	// tells of some.
	Mute bool   `json:"mute,omitempty"`
	UUID string `json:"uuid"`
}

// Segments asks SponsorBlock for the segments of the video id, by their
// start: none when SponsorBlock knows of none.
func (v *Videos) Segments(ctx context.Context, id string) ([]Segment, error) {
	return sponsorSegments(ctx, sponsorBlockAPI, id)
}

func sponsorSegments(ctx context.Context, api, id string) ([]Segment, error) {
	sum := sha256.Sum256([]byte(id))
	prefix := hex.EncodeToString(sum[:])[:4]
	var cats []string
	for _, k := range settings.SegmentKinds {
		cats = append(cats, k.ID)
	}
	c, _ := json.Marshal(cats)
	actions, _ := json.Marshal([]string{"skip", "mute", "poi"})
	q := url.Values{"categories": {string(c)}, "actionTypes": {string(actions)}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api+"/skipSegments/"+prefix+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MyTube")
	res, err := sponsorBlockClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return []Segment{}, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SponsorBlock answered %s", res.Status)
	}
	var videos []struct {
		VideoID  string `json:"videoID"`
		Segments []struct {
			Segment    []float64 `json:"segment"`
			Category   string    `json:"category"`
			ActionType string    `json:"actionType"`
			UUID       string    `json:"UUID"`
		} `json:"segments"`
	}
	if err := json.NewDecoder(res.Body).Decode(&videos); err != nil {
		return nil, err
	}
	out := []Segment{}
	for _, x := range videos {
		if x.VideoID != id {
			continue
		}
		for _, sg := range x.Segments {
			if len(sg.Segment) != 2 || sg.Segment[1] < sg.Segment[0] {
				continue
			}
			out = append(out, Segment{Start: sg.Segment[0], End: sg.Segment[1], Category: sg.Category, Mute: sg.ActionType == "mute", UUID: sg.UUID})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Start < out[j].Start })
	return out, nil
}
