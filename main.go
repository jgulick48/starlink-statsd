package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/jgulick48/starlink-statsd/internal/metrics"
	"github.com/jgulick48/starlink-statsd/internal/starlink"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	var err error
	metrics.Metrics, err = statsd.New(statsDHost)
	if err != nil {
		log.Printf("Error creating stats client %s", err.Error())
	} else {
		metrics.StatsEnabled = true
	}
	client := starlink.NewClient(addr, deploymentId)
	fmt.Println("Running application")
	client.StartStatsLoop(done, time.Duration(scanInterval)*time.Second)
	fmt.Println("exiting")
}

var (
	statsDHost   = "127.0.0.1:8125"
	addr         = "192.168.100.1:9200"
	deploymentId = "Dishy"
	scanInterval = int64(30)
)

func init() {
	flag.StringVar(&addr, "addr", addr, "grpc addr (dishy is at 192.168.100.1:9200, wifi is at 192.168.1.1:9000")
	flag.StringVar(&statsDHost, "statsDHost", statsDHost, "statsD address")
	flag.StringVar(&deploymentId, "deploymentId", deploymentId, "unique Id for this deployment, defaults to Dishy")
	flag.Int64Var(&scanInterval, "scanInterval", scanInterval, "interval in which to scan data from the dish in seconds, defaults to 30")
	flag.Parse()
}
