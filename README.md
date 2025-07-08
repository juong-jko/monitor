# Go Process Monitor

[![Go CI](https://github.com/juong-jko/monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/juong-jko/monitor/actions/workflows/ci.yml)


This is a simple command-line tool written in Go to monitor the CPU and memory usage of a specific process ID (PID) at regular intervals.

## Dependencies

This project uses the following external library:

- [gopsutil](https://github.com/shirou/gopsutil): A cross-platform library for retrieving process and system utilization data.

## Usage

To use the monitor, you need to provide the PID of the process you want to track. You can also provide an optional interval in seconds.

```bash
go run main.go <pid> [interval_in_seconds]
```

- `<pid>`: The process ID of the application you want to monitor.
- `[interval_in_seconds]` (optional): The interval in seconds at which to poll for CPU and memory usage. Defaults to 5 seconds.

### Example

```bash
go run main.go 12345 10
```

This will monitor the process with PID `12345` every `10` seconds.

## Build

You can build the executable using the following command:

```bash
go build -o main.exe .
```

Then you can run the executable directly:

```bash
./main.exe <pid> [interval_in_seconds]
```
