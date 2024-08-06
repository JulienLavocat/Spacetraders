package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const RPS_LIMIT = 2

var burstDepleted = false

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

func FatalIfNotRateLimitError(response *http.Response, err error, logger zerolog.Logger, msg string, args ...interface{}) bool {
	if response.StatusCode == 429 {
		resetAt, err := time.Parse(time.RFC3339, response.Header.Get("x-ratelimit-reset"))
		if err != nil {
			logger.Fatal().Msgf("unable to parse rate limit date: %s", response.Header.Get("x-ratelimit-reset"))
		}

		sleepFor := resetAt.Add(time.Second).Sub(time.Now().UTC())
		// We've hit the burst limit and are now limited by the regular pool (i.e RPS_LIMIT)
		// isBurstLimit := sleepFor.Seconds() <= RPS_LIMIT
		// if isBurstLimit {
		// 	sleepFor += time.Minute
		// }

		logger.Info().Msgf("hitting rate limit, resets at %s (sleeping for %.2f)", resetAt, sleepFor.Seconds())
		time.Sleep(sleepFor)
		return true
	}

	FatalIfHttpError(response, err, logger, msg, args...)
	return false
}
