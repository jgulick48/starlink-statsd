// Package metrics is a thin wrapper over a DogStatsD client that adds base
// tags and refuses to send non-finite values.
package metrics

import (
	"log"
	"math"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

// Sink sends metrics under a namespace with a fixed set of base tags.
type Sink struct {
	c        statsd.ClientInterface
	prefix   string
	baseTags []string
}

// New returns a Sink. prefix is prepended (with a dot) to every metric name.
func New(c statsd.ClientInterface, prefix string, baseTags ...string) *Sink {
	return &Sink{c: c, prefix: prefix, baseTags: baseTags}
}

func (s *Sink) name(n string) string {
	if s.prefix == "" {
		return n
	}
	return s.prefix + "." + n
}

func (s *Sink) tags(extra []string) []string {
	out := make([]string, 0, len(extra)+len(s.baseTags))
	out = append(out, extra...)
	return append(out, s.baseTags...)
}

func (s *Sink) report(name string, err error) {
	if err != nil {
		log.Printf("statsd send %s: %v", name, err)
	}
}

// Gauge sends a gauge, silently skipping NaN and Inf values.
func (s *Sink) Gauge(name string, v float64, tags ...string) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return
	}
	s.report(name, s.c.Gauge(s.name(name), v, s.tags(tags), 1))
}

// Count sends a counter increment.
func (s *Sink) Count(name string, v int64, tags ...string) {
	s.report(name, s.c.Count(s.name(name), v, s.tags(tags), 1))
}

// Timing sends a timing value.
func (s *Sink) Timing(name string, d time.Duration, tags ...string) {
	s.report(name, s.c.Timing(s.name(name), d, s.tags(tags), 1))
}

// Distribution sends a distribution value, skipping non-finite values.
func (s *Sink) Distribution(name string, v float64, tags ...string) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return
	}
	s.report(name, s.c.Distribution(s.name(name), v, s.tags(tags), 1))
}

// Event sends a DogStatsD event, adding the base tags.
func (s *Sink) Event(title, text string, alert statsd.EventAlertType, tags ...string) {
	e := statsd.NewEvent(title, text)
	e.AlertType = alert
	e.Tags = s.tags(tags)
	s.report("event", s.c.Event(e))
}

// Bool converts a bool to 0/1.
func Bool(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// Tag formats a key:value tag.
func Tag(k, v string) string { return k + ":" + v }
