# Ping Discovery Service Library

A flexible Go library for concurrent network connectivity testing and ping discovery.

## Features

- **Concurrent Ping Testing**: High-performance parallel ping execution with configurable worker pools
- **24/7 Service Mode**: Runs continuously with graceful shutdown support
- **One-Shot Execution**: Single execution mode for quick checks
- **Flexible Data Sources**: Works with Redis, Kafka, databases, files - anything you need
- **Production Ready**: Proper error handling, signal handling, and resource cleanup

## Quick Start
```go
package main

import (
    "log"
    "github.com/technarts/ping-discovery/pingdiscovery"
)

func main() {
    args := []string{"server1", "queue1", "topic1", "signal1"}
    
    // Start 24/7 service
    pingdiscovery.PingDiscoveryServiceAllTime(
        args,                    // Configuration arguments
        100,                     // Worker count
        5,                       // Loop interval (seconds)
        log.Printf,              // Logger function
        readFromYourSource,      // Your data reader
        writeToYourSink,         // Your data writer
    )
}

// Implement your data reader
func readFromYourSource(args []string) ([]pingdiscovery.KafkaIpRetryTimeoutMessage, error) {
    // Read from your data source (Redis, Kafka, database, file, etc.)
    return messages, nil
}

// Implement your data writer
func writeToYourSink(args []string, results map[int]bool) error {
    // Write results to your destination
    return nil
}
```

## One-Shot Mode
```go
// Run once and exit
err := pingdiscovery.RunPingDiscoveryOnce(
    args,
    50,              // Worker count
    log.Printf,      // Logger
    readerFunc,      // Data reader
    writerFunc,      // Data writer
)
```

## API Reference

### Core Types
```go
type KafkaIpRetryTimeoutMessage struct {
    Id         int    `json:"id"`
    IP         string `json:"ip_address"`
    RetryCount int    `json:"retry_count"`
    Timeout    int    `json:"timeout"`
}
```

### Main Functions

**PingDiscoveryServiceAllTime** - Starts a 24/7 service
```go
func PingDiscoveryServiceAllTime(
    args []string,           // Configuration arguments
    workerCount int,         // Number of concurrent workers
    loopInterval int,        // Seconds between cycles
    logger LoggerFunc,       // Logging function
    reader ReadMessagesFunc, // Data input function
    poster PostResultsFunc,  // Data output function
)
```

## Support

DMs: https://technarts.slack.com/team/U05E89JS324