package slick

import (
	log "github.com/sirupsen/logrus"
)

// Logging contains the configuration for logrus
type Logging struct {
	Level string `json:"level" mapstructure:"level"`
	Type  string `json:"type" mapstructure:"type"`
}

// getLoggingConfig return the corresponding formatter and level for logging.
func getLoggingConfig(bot *Bot) (log.Formatter, log.Level) {
	var f log.Formatter

	switch bot.Logging.Type {
	case "json":
		f = &log.JSONFormatter{}
	default:
		f = &log.TextFormatter{}
	}

	l, err := log.ParseLevel(bot.Logging.Level)
	if err != nil {
		l = log.InfoLevel
	}
	return f, l
}

// setupLogging choose the config and setup the logging.
func (bot *Bot) setupLogging() error {
	formatter, level := getLoggingConfig(bot)
	log.SetFormatter(formatter)
	log.SetLevel(level)
	return nil
}
