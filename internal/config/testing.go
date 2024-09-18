package config

// TestConfig returns an English config that can be used for testing
func TestConfig() (*Config, error) {
	return New([]byte("version: 0.0.0\nlocale: en"))
}
