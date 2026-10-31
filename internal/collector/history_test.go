package collector

import (
	"math"
	"reflect"
	"testing"

	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

func TestWindow(t *testing.T) {
	tests := []struct {
		name            string
		current, last   uint64
		hasLast         bool
		ringLen, fallbk int
		want            []uint64
	}{
		{"normal advance", 110, 107, true, 100, 30, []uint64{107, 108, 109}},
		{"no new samples", 110, 110, true, 100, 30, []uint64{}},
		{"first poll uses fallback", 500, 0, false, 100, 3, []uint64{497, 498, 499}},
		{"first poll capped at len", 500, 0, false, 4, 30, []uint64{496, 497, 498, 499}},
		{"first poll young dish", 2, 0, false, 100, 30, []uint64{0, 1}},
		{"current went backwards", 5, 1000, true, 100, 3, []uint64{2, 3, 4}},
		{"gap larger than ring", 1000, 100, true, 4, 30, []uint64{996, 997, 998, 999}},
		{"empty ring", 10, 5, true, 0, 3, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Window(tt.current, tt.last, tt.hasLast, tt.ringLen, tt.fallbk)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestSamplesWraparound(t *testing.T) {
	// ring of 4, Current=6 -> newest sample (seq 5) is at index 1.
	h := &pb.DishGetHistoryResponse{
		Current:          6,
		PopPingDropRate:  []float32{0.4, 0.5, 0.2, 0.3}, // idx0=seq4, idx1=seq5
		PopPingLatencyMs: []float32{40, 50, 20, 30},
	}
	seqs := Window(h.Current, 3, true, 4, 30) // 3,4,5
	ss := Samples(h, seqs)
	var lat []float64
	for _, s := range ss {
		lat = append(lat, s.Latency)
	}
	if !reflect.DeepEqual(lat, []float64{30, 40, 50}) {
		t.Fatalf("got %v", lat)
	}
	if ss[0].HasPower {
		t.Fatal("power should be absent")
	}
}

func TestSummarize(t *testing.T) {
	ss := []Sample{
		{Drop: 0, Latency: 30, Down: 100, Up: 10, Power: 40, HasPower: true},
		{Drop: 0.5, Latency: 50, Down: 300, Up: 30, Power: 60, HasPower: true},
		{Drop: 1, Latency: 999, Down: 0, Up: 0, Power: 50, HasPower: true},
		{Drop: 0, Latency: 40, Down: 200, Up: 20, Power: 50, HasPower: true},
	}
	s := Summarize(ss)
	approx := func(name string, got, want float64) {
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("%s: got %v want %v", name, got, want)
		}
	}
	if s.Samples != 4 || s.SecondsDown != 1 {
		t.Errorf("samples=%d down=%d", s.Samples, s.SecondsDown)
	}
	approx("lat mean", s.LatencyMean, 40)
	approx("lat max", s.LatencyMax, 50)
	approx("lat p95", s.LatP95, 50)
	approx("drop mean", s.DropMean, 0.375)
	approx("down mean", s.DownMean, 150)
	approx("down max", s.DownMax, 300)
	approx("power mean", s.PowerMean, 50)
	approx("power min", s.PowerMin, 40)
	approx("power max", s.PowerMax, 60)
	approx("energy", s.EnergyWh, 200.0/3600)
}

func TestSummarizeAllDropped(t *testing.T) {
	s := Summarize([]Sample{{Drop: 1, Latency: 0}, {Drop: 1}})
	if s.HasLatency || s.HasPower || s.SecondsDown != 2 {
		t.Fatalf("%+v", s)
	}
}

func TestPercentile(t *testing.T) {
	v := make([]float64, 100)
	for i := range v {
		v[i] = float64(i + 1)
	}
	if got := Percentile(v, 0.95); got != 95 {
		t.Fatalf("got %v", got)
	}
	if !math.IsNaN(Percentile(nil, 0.95)) {
		t.Fatal("want NaN")
	}
}
