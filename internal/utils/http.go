package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	RPS_LIMIT = 2
	MAX_RETRY = 5
)

var burstDepleted = false

func RetryRequest[T any](execute func() (*T, *http.Response, error), logger zerolog.Logger, msg string, args ...interface{}) *T {
	for i := 0; i < MAX_RETRY; i++ {
		res, http, err := execute()

		if fatalIfNotRateLimitError(http, err, logger, msg, args...) {
			continue
		}

		return res
	}

	logger.Fatal().Err(errors.New("max retry exceeded")).Msgf(msg, args...)
	return nil
}

func fatalIfHttpError(res *http.Response, err error, logger zerolog.Logger, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	event := logger.Fatal().Err(err)

	isJson, json, body, err := readJsonFromBody(res)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to read response body")
	}

	if isJson {
		event.Interface("body", json)
	} else {
		event.Str("body", string(body))
	}

	event.Msgf(msg, args...)
}

func fatalIfNotRateLimitError(response *http.Response, err error, logger zerolog.Logger, msg string, args ...interface{}) bool {
	if response.StatusCode == 429 {
		resetAt, err := time.Parse(time.RFC3339, response.Header.Get("x-ratelimit-reset"))
		if err != nil {
			logger.Fatal().Msgf("unable to parse rate limit date: %s", response.Header.Get("x-ratelimit-reset"))
		}

		sleepFor := resetAt.Add(time.Second).Sub(time.Now().UTC())

		logger.Info().Msgf("hitting rate limit, resets at %s (sleeping for %.2f)", resetAt, sleepFor.Seconds())
		time.Sleep(sleepFor)
		return true
	}

	fatalIfHttpError(response, err, logger, msg, args...)
	return false
}

func readJsonFromBody(response *http.Response) (bool, map[string]any, string, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, nil, "", err
	}

	isJson := strings.HasPrefix(response.Header.Get("Content-Type"), "application/json")
	var jsonBody map[string]any
	if isJson {
		err := json.Unmarshal(body, &jsonBody)
		if err != nil {
			return false, nil, string(body), err
		}
	}

	return isJson, jsonBody, string(body), err
}
