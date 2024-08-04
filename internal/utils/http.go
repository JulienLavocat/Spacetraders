package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func FatalIfHttpError(res *http.Response, err error, logger zerolog.Logger, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	body, readBodyErr := io.ReadAll(res.Body)
	if readBodyErr != nil {
		logger.Fatal().Err(err).Msg("unable to read error body")
	}

	event := logger.Fatal().Err(err)

	if strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		var result map[string]any
		err := json.Unmarshal(body, &result)
		if err != nil {
			logger.Warn().Msg("server sent an invalid json")
			event.Str("body", string(body))
			return
		}

		event.Interface("body", result)
	} else {
		event.Str("body", string(body))
	}

	event.Msgf(msg, args...)
}
