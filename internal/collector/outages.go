package collector

import pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"

// OutageTracker reports each outage once, keyed by start timestamp.
type OutageTracker struct {
	max    int
	seen   map[int64]struct{}
	order  []int64
	primed bool
}

// NewOutageTracker keeps at most max remembered timestamps.
func NewOutageTracker(max int) *OutageTracker {
	return &OutageTracker{max: max, seen: map[int64]struct{}{}}
}

// New returns outages not seen before. The first call only records them and
// returns nothing, so a restart does not replay old outages.
func (t *OutageTracker) New(outages []*pb.DishOutage) []*pb.DishOutage {
	var fresh []*pb.DishOutage
	for _, o := range outages {
		if o == nil {
			continue
		}
		if _, ok := t.seen[o.StartTimestampNs]; ok {
			continue
		}
		t.seen[o.StartTimestampNs] = struct{}{}
		t.order = append(t.order, o.StartTimestampNs)
		if len(t.order) > t.max {
			delete(t.seen, t.order[0])
			t.order = t.order[1:]
		}
		fresh = append(fresh, o)
	}
	if !t.primed {
		t.primed = true
		return nil
	}
	return fresh
}
