package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	monitor "juong.jko/monitor/internal"
)

func main() {
	ctx := context.Background()
	// Define and parse command-line flags
	pid := flag.Int("pid", 0, "Process ID to monitor")
	interval := flag.Duration("interval", 5*time.Second, "Monitoring interval")
	flag.Parse()

	if *pid == 0 {
		log.Fatal("PID must be provided with the -pid flag")
	}

	proc, err := monitor.FindProcess(ctx, *pid)
	if err != nil {
		log.Fatalf("Error finding process with PID %d: %v", *pid, err)
	}

	fmt.Printf("Starting monitoring...\n")
	fmt.Println(".........................................................")

	signalCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a channel to receive monitoring data
	dataChan := make(chan monitor.Info)

	// Start process monitoring in a separate goroutine
	go monitor.MonitorProcess(signalCtx, proc, *interval, dataChan)

	// Create a new broker
	broker := monitor.NewBroker()

	// Start a goroutine to broadcast monitoring data to clients
	go func() {
		for info := range dataChan {
			broker.Mu.Lock()
			for client := range broker.Clients {
				client <- info.String()
			}
			broker.Mu.Unlock()
		}
	}()

	// Set up HTTP server
	http.Handle("/events", broker)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Starting web server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
