package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	*process.Process
}

func FindProcess(ctx context.Context, pid int) (*Process, error) {
	p, err := process.NewProcess(int32(pid))
	return &Process{Process: p}, err
}

// MonitorProcess starts monitoring a process with a given PID at a specified interval.
func MonitorProcess(ctx context.Context, proc *Process, interval time.Duration) {
	// Create a ticker that fires at an interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			{
				exists, err := process.PidExists(int32(proc.Pid))
				if err != nil {
					log.Printf("Error checking if process exists: %v", err)
					continue
				}

				if !exists {
					log.Printf("Process with PID %d has exited", proc.Pid)
					return
				}

				// CPU info
				cpuPercent, err := proc.CPUPercent()
				if err != nil {
					log.Printf("Could not get CPU percentage: %v", err)
					continue
				}

				// Memory info
				memInfo, err := proc.MemoryInfo()
				if err != nil {
					log.Printf("Could not get memory info: %v", err)
					continue
				}

				fmt.Printf("[%s] PID: %d, CPU Usage: %.2f%%, Memory RSS: %d MB\n",
					t.Format("2006-01-02 15:04:05"),
					proc.Pid,
					cpuPercent,
					memInfo.RSS/1024/1024)
			}
		case <-ctx.Done():
			{
				log.Printf("Received signal to terminate")
				return
			}
		}
	}
}
