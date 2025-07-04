package pingdiscovery

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewPingDiscoveryService(args []string, workerCount int, loopInterval int, logger LoggerFunc, reader ReadMessagesFunc, poster PostResultsFunc) *PingDiscoveryService {
	if len(args) < 4 {
		logger("Usage: ./discovery ping-service <bootstrap-server[:port]> <topic-to-read> <topic-to-write> <topic-to-signal>\n")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &PingDiscoveryService{
		server:       args[0],
		readTopic:    args[1],
		writeTopic:   args[2],
		signalTopic:  args[3],
		workerCount:  workerCount,
		loopInterval: loopInterval,
		readData:     reader,
		postResult:   poster,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (s *PingDiscoveryService) Start() {
	s.logger("PING DISCOVERY SERVICE STARTED - Running 24/7\n")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	s.wg.Add(1)
	go s.processLoop()

	select {
	case sig := <-sigChan:
		s.logger("Received signal %v, shutting down gracefully...\n", sig)
		s.cancel()
	case <-s.ctx.Done():
		s.logger("Context cancelled, shutting down...\n")
	}

	s.wg.Wait()
	s.logger("PING DISCOVERY SERVICE STOPPED\n")
}

func (s *PingDiscoveryService) processLoop() {
	defer s.wg.Done()

	retryInterval := time.Duration(s.loopInterval) * time.Second

	for {
		select {
		case <-s.ctx.Done():
			s.logger("Processing loop stopped\n")
			return
		default:
			if err := s.processPingDiscovery(); err != nil {
				s.logger("Error processing ping discovery: %v\n", err)
			}

			select {
			case <-s.ctx.Done():
				return
			case <-time.After(retryInterval):
			}
		}
	}
}

func (s *PingDiscoveryService) processPingDiscovery() error {
	ipRetryTimeoutList, err := s.readIPRetryTimeoutMessages()
	if err != nil {
		return err
	}

	if len(ipRetryTimeoutList) == 0 {
		s.logger("No messages to process, waiting...\n")
		return nil
	}

	s.logger("Processing %d IP tasks\n", len(ipRetryTimeoutList))

	tasks := make([]IPTask, len(ipRetryTimeoutList))
	for i, msg := range ipRetryTimeoutList {
		tasks[i] = IPTask{
			ID:         msg.Id,
			IP:         msg.IP,
			RetryCount: msg.RetryCount,
			Timeout:    msg.Timeout,
		}
	}

	results := RunPingDiscovery(tasks, s.workerCount, s.logger)

	ipSuccessMap := make(map[int]bool)
	for _, res := range results {
		ipSuccessMap[res.ID] = (res.Status == StatusAvailable)
	}

	if err := s.postPingResults(ipSuccessMap); err != nil {
		return err
	}

	s.logger("Successfully processed %d ping results\n", len(results))
	return nil
}

func (s *PingDiscoveryService) postPingResults(ipSuccessMap map[int]bool) error {
	return s.postResult(s.server, s.writeTopic, s.signalTopic, ipSuccessMap)
}

func (s *PingDiscoveryService) readIPRetryTimeoutMessages() ([]KafkaIpRetryTimeoutMessage, error) {
	return s.readData(s.server, s.readTopic)
}

func PingDiscoveryServiceAllTime(args []string, workerCount int, loopInterval int, logger LoggerFunc, reader ReadMessagesFunc, poster PostResultsFunc) {
	if logger == nil {
		logger = func(string, ...any) {}
	}

	service := NewPingDiscoveryService(args, workerCount, loopInterval, logger, reader, poster)
	service.Start()
}
