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

type Info struct {
	Timestamp  time.Time
	PID        int32
	CPUPercent float64
	RSS        uint64
}

func (i *Info) String() string {
	return fmt.Sprintf("[%s] PID: %d, CPU Usage: %.2f%%, Memory RSS: %d MB",
		i.Timestamp.Format("2006-01-02 15:04:05"),
		i.PID,
		i.CPUPercent,
		i.RSS/1024/1024)
}

func FindProcess(ctx context.Context, pid int) (*Process, error) {
	p, err := process.NewProcess(int32(pid))
	return &Process{Process: p}, err
}

// MonitorProcess starts monitoring a process with a given PID at a specified interval.
func MonitorProcess(ctx context.Context, proc *Process, interval time.Duration, dataChan chan<- Info) {
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
					close(dataChan)
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

				info := Info{
					Timestamp:  t,
					PID:        proc.Pid,
					CPUPercent: cpuPercent,
					RSS:        memInfo.RSS,
				}
				dataChan <- info
			}
		case <-ctx.Done():
			{
				log.Printf("Received signal to terminate")
				close(dataChan)
				return
			}
		}
	}
}
