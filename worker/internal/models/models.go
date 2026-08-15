package models

// IPSecLog represents the expected JSON structure from the firewall
type IPSecLog struct {
	Action   string `json:"action"`
	SourceIP string `json:"source_ip"`
}

// Config holds the application configuration
type Config struct {
	QueueURL    string
	RedisURL    string
	DatabaseURL string
}
