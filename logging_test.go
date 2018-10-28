package slick

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var mockedBotDefault = &Bot{}

var mockedBotJSONPanic = &Bot{
	Logging: Logging{
		Type:  "json",
		Level: "panic",
	},
}

var mockedBotTextDebug = &Bot{
	Logging: Logging{
		Type:  "text",
		Level: "debug",
	},
}

func TestShouldReturnTheCorrespondingConfig(t *testing.T) {

	tests := map[string]struct {
		bot       *Bot
		formatter log.Formatter
		level     log.Level
	}{
		"default": {
			bot:       mockedBotDefault,
			formatter: &log.TextFormatter{},
			level:     log.InfoLevel,
		},
		"JSON formatter and Panic Level": {
			bot:       mockedBotJSONPanic,
			formatter: &log.JSONFormatter{},
			level:     log.PanicLevel,
		},
		"Text formatter and Debug Level": {
			bot:       mockedBotJSONPanic,
			formatter: &log.JSONFormatter{},
			level:     log.PanicLevel,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		f, l := getLoggingConfig(test.bot)

		assert.Equal(t, f, test.formatter)
		assert.Equal(t, l, test.level)
	}
}
