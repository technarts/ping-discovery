package pingdiscovery

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
