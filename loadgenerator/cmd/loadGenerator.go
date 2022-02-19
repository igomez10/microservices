package cmd

import (
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type LoadGeneratorConfig struct {
	Entries []RequestConfig `json:"entries"`
}

type RequestConfig struct {
	URL     string `json:"url"`
	Method  string `json:"method"`
	Headers struct {
		Headername string `json:"headername"`
	} `json:"headers,omitempty"`
	Body string `json:"body,omitempty"`
}

var wg sync.WaitGroup

var defaultWait = 1 * time.Second

var httpClient = http.Client{Timeout: 60 * time.Second}
var maxRetry = 5

func Run(conf LoadGeneratorConfig) {
	wg.Add(len(conf.Entries))
	for i := range conf.Entries {
		log.Debug().Msgf("Processing: %+v", conf.Entries[i].URL)
		go IssueRequest(conf.Entries[i])
	}
	wg.Wait()
	log.Warn().Msgf("Exiting Run")
}

func IssueRequest(reqConfig RequestConfig) {
	for {

		req, err := http.NewRequest(reqConfig.Method, reqConfig.URL, nil)
		if err != nil {
			log.Err(err).Msg("Failed to create request")
			return
		}

		// var res http.Response
		for retryCount := 0; retryCount < maxRetry; retryCount++ {
			res, err := httpClient.Do(req)
			if err != nil || res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500 {
				log.Err(err).Msgf("Failed to issue request: Retry %d", retryCount)
				time.Sleep(100 * time.Millisecond)
			} else {
				log.Debug().Msgf("Processed %+v: status code: %d", reqConfig.URL, res.StatusCode)
				break // request was succesful
			}

		}

		log.Debug().Msgf("Processing %+v", reqConfig.URL)
		time.Sleep(defaultWait)
	}

	log.Warn().Msgf("Exiting %+v", reqConfig.URL)
}
