package momolog

import (
	"time"

	"github.com/rs/zerolog"
)

type config struct {
	level      string
	formatTime string
}

func SetLevel(level string) func(*config) {
	return func(c *config) {
		c.level = level
	}
}

func SetFormatTime(formatTime string) func(*config) {
	return func(c *config) {
		c.formatTime = formatTime
	}
}

func defaultConfig() config {
	return config{
		level:      zerolog.InfoLevel.String(),
		formatTime: time.RFC3339,
	}
}

// Level: parse level from string into zerolog.Level
func (c *config) Level() zerolog.Level {
	level, err := zerolog.ParseLevel(c.level)
	if err != nil {
		level = zerolog.TraceLevel
	}

	return level
}
