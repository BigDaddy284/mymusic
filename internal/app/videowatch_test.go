package app

import (
	"testing"
	"time"
)

// TestVideoWatchingReported reports the watching of a video to YouTube's
// history from its start, on the schedule while it plays, and a last time
// as no video plays.
func TestVideoWatchingReported(t *testing.T) {
	w := newWatching(t, true, true)
	vw := &w.a.videoWatch
	vw.now = func() time.Time { return w.clock }
	vw.start = w.a.watchtime.start
	v := &Videos{a: w.a}
	at := func(d time.Duration, id string, s float64, plays bool) {
		w.clock = time.Date(2026, 11, 4, 9, 0, 0, 0, time.UTC).Add(d)
		v.ReportWatching(id, s, plays, 1, false)
	}
	at(0, "vid00000001", 0, true)
	deadline := time.Now().Add(5 * time.Second)
	for {
		vw.mu.Lock()
		ready := vw.cur != nil && vw.cur.watch != nil
		vw.mu.Unlock()
		if ready {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the video's watch didn't start")
		}
		time.Sleep(time.Millisecond)
	}
	if got := w.started(); len(got) != 1 || got[0] != "vid00000001" {
		t.Fatalf("the watches started: %q", got)
	}
	at(5*time.Second, "vid00000001", 5, true)
	at(12*time.Second, "vid00000001", 12, true)
	at(15*time.Second, "", 0, false)
	vw.mu.Lock()
	last := vw.last
	vw.mu.Unlock()
	if last != nil {
		<-last
	}
	got := w.fake.got()
	if len(got) < 2 || !got[len(got)-1].Final || got[len(got)-1].At != 12*time.Second {
		t.Errorf("the reports: %+v", got)
	}
}
