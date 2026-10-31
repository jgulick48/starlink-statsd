// Package collector polls the dish and turns responses into metrics.
package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"github.com/jgulick48/starlink-statsd/internal/metrics"
	pb "github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device"
)

// Source is the dish API used by the collector.
type Source interface {
	Status(ctx context.Context) (*pb.DishGetStatusResponse, error)
	History(ctx context.Context) (*pb.DishGetHistoryResponse, error)
}

const (
	pollTimeout    = 8 * time.Second
	errLogInterval = time.Minute
	maxSeenOutages = 500
	longOutage     = 30 * time.Second
)

// Collector polls a Source on an interval.
type Collector struct {
	src          Source
	m            *metrics.Sink
	interval     time.Duration
	history      bool
	deploymentID string

	lastCurrent uint64
	hasLast     bool
	outages     *OutageTracker

	lastErr    string
	lastErrLog time.Time
}

// New builds a Collector. m must already carry the deploymentId base tag.
func New(src Source, m *metrics.Sink, interval time.Duration, history bool, deploymentID string) *Collector {
	return &Collector{
		src: src, m: m, interval: interval, history: history, deploymentID: deploymentID,
		outages: NewOutageTracker(maxSeenOutages),
	}
}

// Run polls immediately and then every interval until ctx is cancelled.
func (c *Collector) Run(ctx context.Context) {
	t := time.NewTicker(c.interval)
	defer t.Stop()
	for {
		c.poll(ctx)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (c *Collector) poll(parent context.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(parent, pollTimeout)
	err := c.pollOnce(ctx)
	cancel()

	result := "success"
	if err != nil {
		result = "error"
		if parent.Err() == nil {
			c.logErr(err)
		}
	} else if c.lastErr != "" {
		log.Printf("poll recovered after error: %s", c.lastErr)
		c.lastErr = ""
	}
	c.m.Count("collector.poll", 1, metrics.Tag("result", result))
	c.m.Timing("collector.poll_duration_ms", time.Since(start))
}

func (c *Collector) logErr(err error) {
	msg := err.Error()
	if msg == c.lastErr && time.Since(c.lastErrLog) < errLogInterval {
		return
	}
	c.lastErr, c.lastErrLog = msg, time.Now()
	log.Printf("poll failed: %v (repeats suppressed for %s)", err, errLogInterval)
}

func (c *Collector) pollOnce(ctx context.Context) error {
	st, err := c.src.Status(ctx)
	if err != nil {
		return fmt.Errorf("status: %w", err)
	}
	EmitStatus(c.m, st)
	if !c.history {
		return nil
	}
	h, err := c.src.History(ctx)
	if err != nil {
		return fmt.Errorf("history: %w", err)
	}
	c.processHistory(h)
	return nil
}

func (c *Collector) processHistory(h *pb.DishGetHistoryResponse) {
	ringLen := len(h.PopPingDropRate)
	fallback := int(c.interval / time.Second)
	seqs := Window(h.Current, c.lastCurrent, c.hasLast, ringLen, fallback)
	c.lastCurrent, c.hasLast = h.Current, true
	EmitHistory(c.m, Summarize(Samples(h, seqs)))
	for _, o := range c.outages.New(h.Outages) {
		c.emitOutage(o)
	}
}

// gpsEpochOffset is the Unix time of the GPS epoch (1980-01-06T00:00:00Z).
const gpsEpochOffset = 315964800 * time.Second

// outageStart converts an outage start timestamp to UTC. The dish appears to
// count from the GPS epoch rather than Unix. Outages in the history are always
// recent, so a value that reads as more than five years old as Unix time is
// treated as GPS time.
func outageStart(ns int64, now time.Time) time.Time {
	t := time.Unix(0, ns).UTC()
	if now.Sub(t) > 5*365*24*time.Hour {
		t = t.Add(gpsEpochOffset)
	}
	return t
}

func (c *Collector) emitOutage(o *pb.DishOutage) {
	cause := o.Cause.String()
	dur := time.Duration(o.DurationNs)
	start := outageStart(o.StartTimestampNs, time.Now()).Format(time.RFC3339)
	alert := statsd.Info
	if dur > longOutage {
		alert = statsd.Warning
	}
	text := fmt.Sprintf("start: %s\nduration: %s\ndid_switch: %t", start, dur.Round(time.Millisecond), o.DidSwitch)
	tag := metrics.Tag("cause", cause)
	c.m.Event("Starlink outage: "+cause, text, alert, tag)
	c.m.Count("outage.count", 1, tag)
	c.m.Gauge("outage.duration_s", dur.Seconds(), tag)
}
