package pingdiscovery

import (
	"context"
	"sync"
)

type LoggerFunc func(format string, v ...any)

type ReadMessagesFunc func(server, topic string) ([]KafkaIpRetryTimeoutMessage, error)

type PostResultsFunc func(server, writeTopic, signalTopic string, results map[int]bool) error

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
	server       string
	readTopic    string
	writeTopic   string
	signalTopic  string
	logger       LoggerFunc
	workerCount  int
	loopInterval int

	// Datas Read and write functions
	readData   ReadMessagesFunc
	postResult PostResultsFunc

	// Add these for graceful shutdown
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
