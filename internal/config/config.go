package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"

	"cbeimers113/ilo/internal/locale"
)

type Config struct {
	Locale  string `yaml:"locale"`
}

// New creates a new Config from the given data
func New(data []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(data, c)
	return c, err
}

// Message returns the localized message for a given key
func (c *Config) Message(key int) string {
	messages := locale.LocalizedStrings[c.Locale]

	if key >= len(messages) {
		return fmt.Sprintf("no message at index %d for locale \"%s\"", key, c.Locale)
	}

	return messages[key]
}
