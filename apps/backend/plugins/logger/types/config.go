package types

type LoggerPluginConfig struct {
	Enabled bool `json:"enabled" toml:"enabled"`
	// MaxLogCount is the maximum number of logs to keep before stopping
	MaxLogCount int `json:"max_log_count" toml:"max_log_count"`
}

func (c *LoggerPluginConfig) Validate() error {
	if c.MaxLogCount <= 0 {
		c.MaxLogCount = 1000
	}
	return nil
}
