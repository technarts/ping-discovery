package pingdiscovery

import (
	"sync"
	"time"

	"github.com/go-ping/ping"
)

const (
	StatusAvailable   = "AVAILABLE"
	StatusUnavailable = "UNAVAILABLE"
)

type LoggerFunc func(format string, v ...any)

func RunPingDiscovery(tasks []IPTask, workerCount int, logger LoggerFunc) []PingResult {
	if logger == nil {
		logger = func(string, ...any) {}
	}

	logger("Ping Discovery started")

	ipChan := make(chan IPTask, len(tasks))
	resultMap := make(map[int]bool)
	var wg sync.WaitGroup
	var mu sync.RWMutex

	for _, task := range tasks {
		resultMap[task.ID] = false
		ipChan <- task
	}
	close(ipChan)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ipChan {
				success := performPing(task, logger)
				mu.Lock()
				resultMap[task.ID] = success
				mu.Unlock()
			}
		}()
	}
	logger("Waiting for workers to finish")
	wg.Wait()

	var results []PingResult
	successCount := 0
	for id, success := range resultMap {
		status := StatusUnavailable
		if success {
			status = StatusAvailable
			successCount++
		}
		results = append(results, PingResult{ID: id, Status: status})
	}
	logger("%d successful response from %d devices", successCount, len(tasks))
	return results
}

func performPing(task IPTask, logger LoggerFunc) bool {
	pinger, err := ping.NewPinger(task.IP)
	if err != nil {
		logger("Failed to create pinger for IP %s: %v", task.IP, err)
		return false
	}

	pinger.Count = task.RetryCount
	pinger.Timeout = time.Duration(task.Timeout) * time.Second

	success := false
	pinger.OnFinish = func(stats *ping.Statistics) {
		success = stats.PacketsRecv > 0
	}

	if err := pinger.Run(); err != nil {
		logger("Ping run error for IP %s: %v", task.IP, err)
		return false
	}

	return success
}
