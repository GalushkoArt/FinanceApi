package utils

import "github.com/rs/zerolog/log"

func PanicOnError(err error) {
	if err != nil {
		log.Panic().Err(err)
	}
}
