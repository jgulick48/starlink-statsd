package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"github.com/jgulick48/starlink-statsd/internal/collector"
	"github.com/jgulick48/starlink-statsd/internal/metrics"
	"github.com/jgulick48/starlink-statsd/internal/starlink"
)

// Set via -ldflags "-X main.version=... -X main.commit=...".
var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	var (
		addr         = flag.String("addr", "192.168.100.1:9200", "grpc addr (dishy is at 192.168.100.1:9200, wifi is at 192.168.1.1:9000)")
		statsDHost   = flag.String("statsDHost", "127.0.0.1:8125", "statsD address")
		deploymentID = flag.String("deploymentId", "Dishy", "unique Id for this deployment, defaults to Dishy")
		scanInterval = flag.Int64("scanInterval", 30, "interval in which to scan data from the dish in seconds")
		history      = flag.Bool("history", true, "also fetch get_history each poll and emit history/outage metrics")
	)
	flag.Parse()

	log.Printf("starlink-statsd version=%s commit=%s", version, commit)
	log.Printf("config: addr=%s statsDHost=%s deploymentId=%s scanInterval=%ds history=%t",
		*addr, *statsDHost, *deploymentID, *scanInterval, *history)
	if *scanInterval < 1 {
		log.Fatalf("scanInterval must be >= 1")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sd, err := statsd.New(*statsDHost)
	if err != nil {
		log.Printf("fatal: create statsd client: %v", err)
		os.Exit(1)
	}
	dish, err := starlink.New(*addr)
	if err != nil {
		sd.Close()
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}

	sink := metrics.New(sd, "starlink", metrics.Tag("deploymentId", *deploymentID))
	c := collector.New(dish, sink, time.Duration(*scanInterval)*time.Second, *history, *deploymentID)
	log.Printf("running")
	c.Run(ctx)

	log.Printf("shutting down")
	dish.Close()
	sd.Close()
}
