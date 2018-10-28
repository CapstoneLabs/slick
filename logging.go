package slick

import (
	log "github.com/sirupsen/logrus"
)

// Logging contains the configuration for logrus
type Logging struct {
	Level string `json:"level" mapstructure:"level"`
	Type  string `json:"type" mapstructure:"type"`
}

func (bot *Bot) setupLogging() error {
	// Set the formatter
	switch bot.Logging.Type {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	// Correlate a level from a string
	level, err := log.ParseLevel(bot.Logging.Level)
	if err != nil {
		level = log.InfoLevel
	}

	// Set the level
	log.SetLevel(level)

	return nil
}
