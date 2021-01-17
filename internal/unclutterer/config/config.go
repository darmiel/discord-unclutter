package config

import (
	"github.com/BurntSushi/toml"
	"strings"
	"time"
)

type Config struct {
	DiscordConfig   `toml:"Discord"`
	DatabaseConfig  `toml:"Database"`
	CooldownConfig  `toml:"Cooldown"`
	MessageConfig   `toml:"Messages"`
	LoggingConfig   `toml:"Logging"`
	GhostPingConfig `toml:"GhostPing"`
}

func LoadConfig() (conf *Config, err error) {
	conf = &Config{}
	_, err = toml.DecodeFile("config.toml", conf)
	return
}

/// Duration
type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

///

/// Placeholder
type PlaceholderMessage string

func createPlaceholder(key string) string {
	return "{{ " + key + " }}"
}

func (m PlaceholderMessage) Repl(repl ...string) (msg string) {
	// convert to string
	msg = string(m)

	// placeholders
	var key string
	for _, data := range repl {
		// make replacement
		if key != "" {
			msg = strings.ReplaceAll(msg, createPlaceholder(key), data)
			// reset key
			key = ""
		} else {
			key = data
		}
	}

	// return as string
	return
}

///
