package global

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var Logger zerolog.Logger

func NewGlobalLogger(logfile io.Writer) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	Logger = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, logfile)).
		With().
		Timestamp().
		Logger()

	var logLevel zerolog.Level
	state := viper.GetString("APP_STATE")
	switch state {
	case "development":
		logLevel = zerolog.DebugLevel
	case "testing":
		logLevel = zerolog.WarnLevel
	case "production":
		logLevel = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	Logger.
		Info().
		Int("log_level", int(logLevel)).
		Str("app_state", state).
		Msgf("set log level to: %s", logLevel)
	return Logger
}
