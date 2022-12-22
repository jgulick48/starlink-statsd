package main

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/jgulick48/starlink-statsd/internal/ecowitt"
	"github.com/jgulick48/starlink-statsd/internal/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	metrics.Metrics, err = statsd.New("192.168.3.92:8125")
	if err != nil {
		log.Printf("Error creating stats client %s", err.Error())
	} else {
		metrics.StatsEnabled = true
	}
	client := ecowitt.NewClient("192.168.3.134", 5*time.Second, http.DefaultClient)
	client.StartScan()
	fmt.Println("Running application")
	<-done
	client.StopScan()
	fmt.Println("exiting")
}
