# Go Process Monitor

[![Go CI](https://github.com/juong-jko/monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/juong-jko/monitor/actions/workflows/ci.yml)

This is a simple web-based tool written in Go to monitor the CPU and memory usage of a specific process ID (PID) in real-time.

## Features

-   **Real-time Monitoring**: CPU and memory usage is streamed to the web interface using Server-Sent Events (SSE).
-   **Web-Based Interface**: A simple HTML interface displays the monitoring data.
-   **Configurable**: The process ID and monitoring interval can be configured via command-line flags.

## Dependencies

This project uses the following external library:

-   [gopsutil](https://github.com/shirou/gopsutil): A cross-platform library for retrieving process and system utilization data.

## Usage

To use the monitor, you need to provide the PID of the process you want to track. You can also provide an optional interval.

```bash
go run main.go -pid <pid> -interval <interval>
```

-   `<pid>`: The process ID of the application you want to monitor.
-   `<interval>` (optional): The interval at which to poll for CPU and memory usage (e.g., 5s, 1m). Defaults to 5 seconds.

### Example

```bash
go run main.go -pid 12345 -interval 10s
```

This will start a web server on `http://localhost:8080` and monitor the process with PID `12345` every `10` seconds.

## Web Interface

Once the server is running, you can view the monitoring data by opening a web browser and navigating to `http://localhost:8080`.

## Build

You can build the executable using the following command:

```bash
go build -o main.exe .
```

Then you can run the executable directly:

```bash
./main.exe -pid <pid> -interval <interval>
```