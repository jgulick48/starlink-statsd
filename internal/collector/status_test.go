package collector

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"github.com/jgulick48/starlink-statsd/internal/metrics"
	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

type call struct {
	kind  string
	name  string
	value float64
	tags  []string
}

type fakeClient struct {
	statsd.NoOpClient
	calls  []call
	events []*statsd.Event
}

func (f *fakeClient) Gauge(n string, v float64, t []string, _ float64) error {
	f.calls = append(f.calls, call{"gauge", n, v, t})
	return nil
}
func (f *fakeClient) Count(n string, v int64, t []string, _ float64) error {
	f.calls = append(f.calls, call{"count", n, float64(v), t})
	return nil
}
func (f *fakeClient) Timing(n string, d time.Duration, t []string, _ float64) error {
	f.calls = append(f.calls, call{"timing", n, float64(d.Milliseconds()), t})
	return nil
}
func (f *fakeClient) Event(e *statsd.Event) error { f.events = append(f.events, e); return nil }

func (f *fakeClient) find(name string, tags ...string) (call, bool) {
	for _, c := range f.calls {
		if c.name != name {
			continue
		}
		ok := true
		for _, want := range tags {
			found := false
			for _, tg := range c.tags {
				found = found || tg == want
			}
			ok = ok && found
		}
		if ok {
			return c, true
		}
	}
	return call{}, false
}

func newSink() (*fakeClient, *metrics.Sink) {
	f := &fakeClient{}
	return f, metrics.New(f, "starlink", "deploymentId:T")
}

func TestEmitStatus(t *testing.T) {
	f, m := newSink()
	st := &pb.DishGetStatusResponse{
		HasActuators:                  pb.HasActuators_HAS_ACTUATORS_YES,
		MobilityClass:                 pb.UserMobilityClass_STATIONARY,
		PopPingLatencyMs:              float32(math.NaN()),
		Alerts:                        &pb.DishAlerts{DishWaterDetected: true, IsPowerSaveIdle: true},
		ObstructionStats:              &pb.DishObstructionStats{AvgProlongedObstructionDurationS: float32(math.NaN())},
		AlignmentStats:                &pb.AlignmentStats{TiltAngleDeg: 3},
		InitializationDurationSeconds: &pb.InitializationDurationSeconds{GpsValid: 12, RfReady: -1},
	}
	EmitStatus(m, st)

	if c, ok := f.find("starlink.hasactuators", "has_actuators:HAS_ACTUATORS_YES"); !ok || c.value != float64(pb.HasActuators_HAS_ACTUATORS_YES) {
		t.Errorf("hasactuators wrong: %+v", c)
	}
	if _, ok := f.find("starlink.mobilityclass", "mobility_class:STATIONARY", "deploymentId:T"); !ok {
		t.Error("mobility tag not key:value")
	}
	for _, key := range []string{"class_of_service:", "software_update_state:", "disablement_code:"} {
		found := false
		for _, c := range f.calls {
			for _, tg := range c.tags {
				found = found || strings.HasPrefix(tg, key)
			}
		}
		if !found {
			t.Errorf("missing tag %s", key)
		}
	}
	if _, ok := f.find("starlink.poppinglatencyms"); ok {
		t.Error("NaN gauge was sent")
	}
	for _, c := range f.calls {
		if c.name == "starlink." || c.name == "starlink" {
			t.Errorf("junk metric %q", c.name)
		}
		if math.IsNaN(c.value) || math.IsInf(c.value, 0) {
			t.Errorf("non-finite sent: %+v", c)
		}
	}
	if _, ok := f.find("starlink.obstructionstats.avgprolongedobstructiondurrations"); ok {
		t.Error("prolonged duration sent while invalid")
	}
	if _, ok := f.find("starlink.obstructionstats.avgprolongedobstructionintervals"); ok {
		t.Error("prolonged interval sent while invalid")
	}
	if _, ok := f.find("starlink.obstructionstats.avgprolongedobstructionvalid"); !ok {
		t.Error("valid flag missing")
	}
	if c, ok := f.find("starlink.alerts", "alertType:DishWaterDetected"); !ok || c.value != 1 {
		t.Error("new alert missing")
	}
	if _, ok := f.find("starlink.alerts", "alertType:IsPowerSaveIdle"); !ok {
		t.Error("IsPowerSaveIdle missing")
	}
	if _, ok := f.find("starlink.alerts", "alertType:Roaming"); ok {
		t.Error("Roaming should be gone")
	}
	if _, ok := f.find("starlink.alignment.tiltangledeg"); !ok {
		t.Error("alignment missing")
	}
	if _, ok := f.find("starlink.initduration.gpsvalid"); !ok {
		t.Error("initduration gpsvalid missing")
	}
	if _, ok := f.find("starlink.initduration.rfready"); ok {
		t.Error("negative initduration sent")
	}
	if c, ok := f.find("starlink.outage.active", "cause:NONE"); !ok || c.value != 0 {
		t.Error("outage.active none missing")
	}
	if _, ok := f.find("starlink.upsu.dishpowerw"); ok {
		t.Error("upsu sent without UpsuStats")
	}
}

func TestEmitStatusValidProlongedAndOutage(t *testing.T) {
	f, m := newSink()
	EmitStatus(m, &pb.DishGetStatusResponse{
		ObstructionStats: &pb.DishObstructionStats{AvgProlongedObstructionValid: true, AvgProlongedObstructionDurationS: 2},
		Outage:           &pb.DishOutage{Cause: pb.DishOutage_OBSTRUCTED},
	})
	if c, ok := f.find("starlink.obstructionstats.avgprolongedobstructiondurrations"); !ok || c.value != 2 {
		t.Error("gated field missing when valid")
	}
	if c, ok := f.find("starlink.outage.active", "cause:OBSTRUCTED"); !ok || c.value != 1 {
		t.Error("active outage missing")
	}
}

type fakeSource struct {
	st  *pb.DishGetStatusResponse
	h   *pb.DishGetHistoryResponse
	err error
}

func (s *fakeSource) Status(context.Context) (*pb.DishGetStatusResponse, error) { return s.st, s.err }
func (s *fakeSource) History(context.Context) (*pb.DishGetHistoryResponse, error) {
	return s.h, s.err
}

func TestPollOutageEventsAndErrors(t *testing.T) {
	f, m := newSink()
	src := &fakeSource{st: &pb.DishGetStatusResponse{}, h: &pb.DishGetHistoryResponse{
		Current: 3, PopPingDropRate: []float32{0, 0, 0}, PopPingLatencyMs: []float32{1, 2, 3},
		Outages: []*pb.DishOutage{{Cause: pb.DishOutage_NO_SATS, StartTimestampNs: 1, DurationNs: uint64(5 * time.Second)}},
	}}
	c := New(src, m, 10*time.Second, true, "T")
	c.poll(context.Background())
	if len(f.events) != 0 {
		t.Fatal("first poll must not emit outage events")
	}
	src.h.Outages = append(src.h.Outages, &pb.DishOutage{Cause: pb.DishOutage_OBSTRUCTED, StartTimestampNs: 2, DurationNs: uint64(45 * time.Second), DidSwitch: true})
	c.poll(context.Background())
	c.poll(context.Background())
	if len(f.events) != 1 {
		t.Fatalf("want 1 event, got %d", len(f.events))
	}
	e := f.events[0]
	if e.Title != "Starlink outage: OBSTRUCTED" || e.AlertType != statsd.Warning {
		t.Errorf("event: %+v", e)
	}
	if _, ok := f.find("starlink.outage.count", "cause:OBSTRUCTED"); !ok {
		t.Error("outage.count missing")
	}
	if _, ok := f.find("starlink.collector.poll", "result:success"); !ok {
		t.Error("poll success counter missing")
	}
	src.err = context.DeadlineExceeded
	c.poll(context.Background())
	if _, ok := f.find("starlink.collector.poll", "result:error"); !ok {
		t.Error("poll error counter missing")
	}
	if _, ok := f.find("starlink.collector.poll_duration_ms"); !ok {
		t.Error("poll duration missing")
	}
}
