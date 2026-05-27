package logging

import (
	"log"
	"log/slog"
	"organization_structure/internal/config"
)

func InitLogger(config *config.Config) {
	var ll slog.Level
	err := ll.UnmarshalText([]byte(config.Log.Level))

	if err != nil {
		log.Fatal(err)
	}

	slog.SetLogLoggerLevel(ll)
}
