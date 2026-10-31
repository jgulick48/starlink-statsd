package collector

import (
	"testing"
	"time"

	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

func outs(ts ...int64) []*pb.DishOutage {
	var o []*pb.DishOutage
	for _, t := range ts {
		o = append(o, &pb.DishOutage{StartTimestampNs: t})
	}
	return o
}

func TestOutageTracker(t *testing.T) {
	tr := NewOutageTracker(3)
	if got := tr.New(outs(1, 2)); len(got) != 0 {
		t.Fatalf("first poll must be suppressed, got %d", len(got))
	}
	if got := tr.New(outs(1, 2)); len(got) != 0 {
		t.Fatalf("known outages re-emitted: %d", len(got))
	}
	if got := tr.New(outs(1, 2, 3)); len(got) != 1 || got[0].StartTimestampNs != 3 {
		t.Fatalf("new outage: %v", got)
	}
	if got := tr.New(outs(1, 2, 3)); len(got) != 0 {
		t.Fatalf("duplicate emitted")
	}
	// Exceed bound: 1 is evicted, so it counts as new if it reappears.
	tr.New(outs(4))
	if len(tr.seen) != 3 || len(tr.order) != 3 {
		t.Fatalf("set unbounded: %d", len(tr.seen))
	}
	if got := tr.New(outs(1)); len(got) != 1 {
		t.Fatalf("evicted entry should be new again")
	}
}

func TestOutageTrackerFirstPollEmpty(t *testing.T) {
	tr := NewOutageTracker(10)
	tr.New(nil)
	if got := tr.New(outs(7)); len(got) != 1 {
		t.Fatalf("outage after empty first poll should emit, got %d", len(got))
	}
}

func TestOutageStart(t *testing.T) {
	unix := time.Date(2026, 10, 31, 3, 0, 0, 0, time.UTC)
	if got := outageStart(unix.UnixNano(), unix.Add(time.Hour)); !got.Equal(unix) {
		t.Errorf("unix timestamp: got %s, want %s", got, unix)
	}
	gps := unix.Add(-gpsEpochOffset).UnixNano()
	if got := outageStart(gps, unix.Add(time.Hour)); !got.Equal(unix) {
		t.Errorf("gps timestamp: got %s, want %s", got, unix)
	}
}
