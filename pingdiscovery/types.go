package pingdiscovery

import (
	"context"
	"sync"
)

type LoggerFunc func(format string, v ...any)

type ReadMessagesFunc func(args []string) ([]KafkaIpRetryTimeoutMessage, error)

type PostResultsFunc func(args []string, results map[int]bool) error

type IPTask struct {
	ID         int    `json:"id"`
	IP         string `json:"ip_address"`
	RetryCount int    `json:"retry_count"`
	Timeout    int    `json:"timeout"`
}

type PingResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type PingDiscoveryService struct {
	logger       LoggerFunc
	workerCount  int
	loopInterval int

	// Store original args for flexible parameter passing
	originalArgs []string

	// Data read and write functions
	readData   ReadMessagesFunc
	postResult PostResultsFunc

	// Graceful shutdown support
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type KafkaIpRetryTimeoutMessage struct {
	Id         int    `json:"id"`
	IP         string `json:"ip_address"`
	RetryCount int    `json:"retry_count"`
	Timeout    int    `json:"timeout"`
}

// Get all args for flexible usage
func (s *PingDiscoveryService) GetArgs() []string {
	return s.originalArgs
}
