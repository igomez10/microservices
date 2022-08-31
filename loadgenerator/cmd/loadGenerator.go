package cmd

import (
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
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
var httpClient = http.Client{Timeout: 10 * time.Second}

const DEFAULT_WAIT = 500 * time.Millisecond
const MAX_RETRY = 5

func Run(conf LoadGeneratorConfig, metricOutput io.Writer) {

	wg.Add(len(conf.Entries))
	for i := range conf.Entries {
		log.Debug().Msgf("Processing: %+v", conf.Entries[i].URL)
		go IssueRequest(conf.Entries[i], metricOutput)
	}
	wg.Wait()
	log.Warn().Msgf("Exiting Run")
}

func IssueRequest(reqConfig RequestConfig, metricOutput io.Writer) {
	for {

		req, err := http.NewRequest(reqConfig.Method, reqConfig.URL, nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create request")
		}

		startTime := time.Now()
		res, err := httpClient.Do(req)
		finishTime := time.Now()
		if err != nil ||
			res.StatusCode == http.StatusTooManyRequests ||
			res.StatusCode >= http.StatusInternalServerError {

			log.Error().Err(err).
				Int("StatusCode", res.StatusCode).
				Str("URL", reqConfig.URL).
				Str("Method", reqConfig.Method).
				Msgf("Failed to issue request")
			//  retry with exponential backoff
		} else {
			log.Debug().
				Str("URL", reqConfig.URL).
				Str("Method", reqConfig.Method).
				Int("StatusCode", res.StatusCode).
				Int64("Latency", time.Since(startTime).Milliseconds()).
				Msgf("Processed")
		}

		latency := int(finishTime.Sub(startTime).Milliseconds())

		emf.New(emf.WithWriter(metricOutput), emf.WithLogGroup("ecsloadgenerator")).
			Namespace("loadgenerator").
			DimensionSet(
				emf.NewDimension("URL", reqConfig.URL),
				emf.NewDimension("Method", reqConfig.Method),
				emf.NewDimension("StatusCode", strconv.Itoa(res.StatusCode)),
			).
			MetricAs("Latency", latency, emf.Milliseconds).
			Log()

		time.Sleep(DEFAULT_WAIT)
	}

}
