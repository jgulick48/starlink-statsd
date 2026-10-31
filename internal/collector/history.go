package collector

import (
	"math"
	"sort"

	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

// Window returns the sequence numbers (monotonic sample counters) of the
// samples that should be processed this poll, oldest first.
//
// current is the total number of samples the dish has ever written, last is
// the value seen on the previous poll (valid only if hasLast), ringLen is the
// ring buffer length and fallback is how many samples to take when there is no
// usable previous value (first poll or the counter went backwards).
func Window(current, last uint64, hasLast bool, ringLen, fallback int) []uint64 {
	if ringLen <= 0 || current == 0 {
		return nil
	}
	n := uint64(fallback)
	if hasLast && current >= last {
		n = current - last
	}
	if n > uint64(ringLen) {
		n = uint64(ringLen)
	}
	if n > current {
		n = current
	}
	out := make([]uint64, 0, n)
	for s := current - n; s < current; s++ {
		out = append(out, s)
	}
	return out
}

func at(arr []float32, seq uint64) (float64, bool) {
	if len(arr) == 0 {
		return 0, false
	}
	return float64(arr[seq%uint64(len(arr))]), true
}

// Sample is one second of history.
type Sample struct {
	Drop, Latency, Down, Up, Power float64
	HasPower                       bool
}

// Samples extracts the samples for the given sequence numbers.
func Samples(h *pb.DishGetHistoryResponse, seqs []uint64) []Sample {
	out := make([]Sample, 0, len(seqs))
	for _, s := range seqs {
		var sm Sample
		sm.Drop, _ = at(h.PopPingDropRate, s)
		sm.Latency, _ = at(h.PopPingLatencyMs, s)
		sm.Down, _ = at(h.DownlinkThroughputBps, s)
		sm.Up, _ = at(h.UplinkThroughputBps, s)
		sm.Power, sm.HasPower = at(h.PowerIn, s)
		out = append(out, sm)
	}
	return out
}

// Summary is the aggregate of a window of samples.
type Summary struct {
	Samples                         int
	SecondsDown                     int
	DropMean                        float64
	LatencyMean, LatencyMax, LatP95 float64
	HasLatency                      bool
	DownMean, DownMax               float64
	UpMean, UpMax                   float64
	PowerMean, PowerMax, PowerMin   float64
	EnergyWh                        float64
	HasPower                        bool
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// Summarize computes window statistics. Latency ignores samples whose drop
// rate is 1 (or non-finite latency), since latency is meaningless there.
func Summarize(ss []Sample) Summary {
	var s Summary
	s.Samples = len(ss)
	if len(ss) == 0 {
		return s
	}
	var dropSum, downSum, upSum, latSum, powSum float64
	var dropN, downN, upN, powN int
	var lats []float64
	for _, x := range ss {
		if finite(x.Drop) {
			dropSum += x.Drop
			dropN++
			if x.Drop >= 1 {
				s.SecondsDown++
			}
		}
		if finite(x.Down) {
			downSum += x.Down
			downN++
			s.DownMax = math.Max(s.DownMax, x.Down)
		}
		if finite(x.Up) {
			upSum += x.Up
			upN++
			s.UpMax = math.Max(s.UpMax, x.Up)
		}
		if x.Drop < 1 && finite(x.Latency) {
			lats = append(lats, x.Latency)
			latSum += x.Latency
		}
		if x.HasPower && finite(x.Power) {
			if powN == 0 {
				s.PowerMax, s.PowerMin = x.Power, x.Power
			}
			powSum += x.Power
			powN++
			s.PowerMax = math.Max(s.PowerMax, x.Power)
			s.PowerMin = math.Min(s.PowerMin, x.Power)
		}
	}
	if dropN > 0 {
		s.DropMean = dropSum / float64(dropN)
	}
	if downN > 0 {
		s.DownMean = downSum / float64(downN)
	}
	if upN > 0 {
		s.UpMean = upSum / float64(upN)
	}
	if len(lats) > 0 {
		s.HasLatency = true
		s.LatencyMean = latSum / float64(len(lats))
		sort.Float64s(lats)
		s.LatencyMax = lats[len(lats)-1]
		s.LatP95 = Percentile(lats, 0.95)
	}
	if powN > 0 {
		s.HasPower = true
		s.PowerMean = powSum / float64(powN)
		s.EnergyWh = powSum / 3600
	}
	return s
}

// Percentile returns the nearest-rank percentile of an ascending-sorted slice.
func Percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return math.NaN()
	}
	rank := int(math.Ceil(p * float64(len(sorted))))
	if rank < 1 {
		rank = 1
	}
	if rank > len(sorted) {
		rank = len(sorted)
	}
	return sorted[rank-1]
}
