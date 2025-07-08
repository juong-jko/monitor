package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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

	// Start process monitoring
	monitor.MonitorProcess(signalCtx, proc, *interval)
}
