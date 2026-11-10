package app

// The stats of the Video side: how long the user watched each video, day
// by day, kept in video-stats.json beside the library.

import (
	"cmp"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/bigdaddy284/mymusic/internal/library"
	"github.com/bigdaddy284/mymusic/internal/store"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

const (
	// videoViewAfter is how long a video plays before it counts as
	// watched, as a track counts after thirty seconds.
	videoViewAfter = 30.0
	// videoSaveEvery is how many seconds of watching wait in memory
	// before they're saved.
	videoSaveEvery = 30.0
	// videoStatsDays is how many days the chart covers for all time.
	videoStatsDays = 90
)

// videoBook is video-stats.json.
type videoBook struct {
	// Days are the seconds each video was watched each day, by the day,
	// as 2026-11-10, and the video's ID.
	Days map[string]map[string]float64 `json:"days"`
	// Views are the videos counted as watched each day.
	Views map[string][]string `json:"views"`
	// Videos tell what each video is.
	Videos map[string]youtube.Video `json:"videos"`
}

// VideoWatch is a video, with how long it was watched and how often.
type VideoWatch struct {
	Video   youtube.Video `json:"video"`
	Seconds float64       `json:"seconds"`
	Views   int           `json:"views"`
}

// VideoChannel is a channel, with how long its videos were watched.
type VideoChannel struct {
	Name    string  `json:"name"`
	Seconds float64 `json:"seconds"`
	Views   int     `json:"views"`
}

// VideoStats sums up what the user watched.
type VideoStats struct {
	Seconds     float64            `json:"seconds"`
	Views       int                `json:"views"`
	Daily       []library.DayPlays `json:"daily"`
	TopVideos   []VideoWatch       `json:"topVideos"`
	TopChannels []VideoChannel     `json:"topChannels"`
}

// videoStats keeps the stats of the Video side, with what waits to be
// saved.
type videoStats struct {
	mu   sync.Mutex
	file *store.File[videoBook]
	now  func() time.Time
	// book is the file's value, with what waits to be saved.
	book    videoBook
	waiting float64
	// sitting is the video watched now, and for how long, which counts
	// as a view after videoViewAfter.
	sitting   string
	sat       float64
	satCounts bool
}

func openVideoStats(path string) (*videoStats, error) {
	f, err := store.Open(path, videoBook{})
	if err != nil {
		return nil, err
	}
	return &videoStats{file: f, now: time.Now, book: f.Get()}, nil
}

// watched adds seconds of watching v, now.
func (vs *videoStats) watched(v youtube.Video, seconds float64) {
	if v.ID == "" || seconds <= 0 {
		return
	}
	vs.mu.Lock()
	defer vs.mu.Unlock()
	b := &vs.book
	day := vs.now().Format(time.DateOnly)
	if b.Days == nil {
		b.Days = map[string]map[string]float64{}
	}
	if b.Days[day] == nil {
		b.Days[day] = map[string]float64{}
	}
	b.Days[day][v.ID] += seconds
	if b.Videos == nil {
		b.Videos = map[string]youtube.Video{}
	}
	known := b.Videos[v.ID]
	known.ID = v.ID
	known.Title = cmp.Or(v.Title, known.Title)
	known.Channel = cmp.Or(v.Channel, known.Channel)
	if len(v.Thumbnail) > 0 {
		known.Thumbnail = v.Thumbnail
	}
	// Views and Published change; the stats keep what doesn't.
	known.Views, known.Published = "", ""
	known.Duration = cmp.Or(v.Duration, known.Duration)
	b.Videos[v.ID] = known
	if vs.sitting != v.ID {
		vs.sitting, vs.sat, vs.satCounts = v.ID, 0, false
	}
	vs.sat += seconds
	if !vs.satCounts && vs.sat >= videoViewAfter {
		vs.satCounts = true
		if b.Views == nil {
			b.Views = map[string][]string{}
		}
		b.Views[day] = append(b.Views[day], v.ID)
	}
	vs.waiting += seconds
	if vs.waiting >= videoSaveEvery {
		vs.saveLocked()
	}
}

// started tells a video starts again from its start, which counts as a
// view of its own once it plays long enough.
func (vs *videoStats) started() {
	vs.mu.Lock()
	vs.sitting = ""
	vs.mu.Unlock()
}

func (vs *videoStats) saveLocked() {
	if vs.waiting == 0 {
		return
	}
	vs.waiting = 0
	if err := vs.file.Set(vs.book); err != nil {
		log.Printf("mytube: keeping the video stats: %v", err)
	}
}

// flush saves what waits.
func (vs *videoStats) flush() {
	vs.mu.Lock()
	vs.saveLocked()
	vs.mu.Unlock()
}

// stats sums up the last days, or all time when days is 0.
func (vs *videoStats) stats(days int) VideoStats {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	now := vs.now()
	span := cmp.Or(days, videoStatsDays)
	out := VideoStats{TopVideos: []VideoWatch{}, TopChannels: []VideoChannel{}}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1-span)
	index := map[string]int{}
	for i := range span {
		d := start.AddDate(0, 0, i).Format(time.DateOnly)
		index[d] = i
		out.Daily = append(out.Daily, library.DayPlays{Date: d})
	}
	since := ""
	if days > 0 {
		since = start.Format(time.DateOnly)
	}
	secs := map[string]float64{}
	views := map[string]int{}
	daySecs := make([]float64, span)
	for day, vids := range vs.book.Days {
		if day < since {
			continue
		}
		for id, s := range vids {
			secs[id] += s
			out.Seconds += s
			if i, ok := index[day]; ok {
				daySecs[i] += s
			}
		}
	}
	for day, ids := range vs.book.Views {
		if day < since {
			continue
		}
		for _, id := range ids {
			views[id]++
			out.Views++
			if i, ok := index[day]; ok {
				out.Daily[i].Plays++
			}
		}
	}
	for i := range out.Daily {
		out.Daily[i].Minutes = int(daySecs[i] / 60)
	}
	channels := map[string]*VideoChannel{}
	for id, s := range secs {
		v := vs.book.Videos[id]
		v.ID = id
		out.TopVideos = append(out.TopVideos, VideoWatch{Video: v, Seconds: s, Views: views[id]})
		if v.Channel == "" {
			continue
		}
		ch := channels[v.Channel]
		if ch == nil {
			ch = &VideoChannel{Name: v.Channel}
			channels[v.Channel] = ch
		}
		ch.Seconds += s
		ch.Views += views[id]
	}
	for _, ch := range channels {
		out.TopChannels = append(out.TopChannels, *ch)
	}
	slices.SortFunc(out.TopVideos, func(a, b VideoWatch) int {
		return cmp.Or(cmp.Compare(b.Seconds, a.Seconds), cmp.Compare(a.Video.ID, b.Video.ID))
	})
	slices.SortFunc(out.TopChannels, func(a, b VideoChannel) int {
		return cmp.Or(cmp.Compare(b.Seconds, a.Seconds), cmp.Compare(a.Name, b.Name))
	})
	out.TopVideos = out.TopVideos[:min(len(out.TopVideos), 20)]
	out.TopChannels = out.TopChannels[:min(len(out.TopChannels), 12)]
	return out
}

// Watched adds seconds of watching v.
func (v *Videos) Watched(video youtube.Video, seconds float64) {
	if vs := v.a.videoStats; vs != nil {
		vs.watched(video, seconds)
	}
}

// Started tells the video plays again from its start.
func (v *Videos) Started() {
	if vs := v.a.videoStats; vs != nil {
		vs.started()
	}
}

// Stats sums up what the user watched in the last days, or all time when
// days is 0.
func (v *Videos) Stats(days int) VideoStats {
	if vs := v.a.videoStats; vs != nil {
		return vs.stats(days)
	}
	return VideoStats{TopVideos: []VideoWatch{}, TopChannels: []VideoChannel{}}
}
