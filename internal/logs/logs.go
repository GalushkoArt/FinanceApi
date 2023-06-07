package logs

import (
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func Init(level string, logsPath string) {
	file, err := os.OpenFile(logsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	utils.PanicOnError(err)
	logLevel, err := zerolog.ParseLevel(level)
	utils.PanicOnError(err)
	logger := zerolog.New(file).Level(logLevel).With().Timestamp().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = logger
}
