package global

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func NewLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Logger = log.Logger
	Logger.Info().Msg("hello world")
}
