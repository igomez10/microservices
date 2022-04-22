package cmd

import (
	"fmt"
	"net"
	"net/http"
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

const DEFAULT_WAIT = 500 * time.Millisecond
const MAX_RETRY = 5

var httpClient = http.Client{Timeout: 10 * time.Second}
var CWAGENT_CONNECTION net.Conn

func Run(conf LoadGeneratorConfig) {
	log.Info().Msg("Sender: Dial OK.")

	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:25888", time.Millisecond*10000)
		if err != nil {
			log.Error().Err(err).Msgf("failed to connect to cloudwatch agent, attempt %d", i)
			time.Sleep(1 * time.Second)

			if i == 9 {
				log.Fatal().Msgf("Exhausted all attempts (%d) to connect to cwagent: %s", i, err)
			}

			continue
		}
		CWAGENT_CONNECTION = conn
		defer conn.Close()
		break
	}

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
		}

		for retryCount := 0; retryCount < MAX_RETRY; retryCount++ {
			startTime := time.Now()
			res, err := httpClient.Do(req)
			if err != nil ||
				res.StatusCode == http.StatusTooManyRequests ||
				res.StatusCode >= http.StatusInternalServerError {

				log.Err(err).
					Int("Retry", retryCount).
					Int("StatusCode", res.StatusCode).
					Str("URL", reqConfig.URL).
					Str("Method", reqConfig.Method).
					Msgf("Failed to issue request")

				//  retry with exponential backoff
				time.Sleep(100 * time.Millisecond * time.Duration(retryCount))
			} else {
				log.Debug().
					Str("URL", reqConfig.URL).
					Str("Method", reqConfig.Method).
					Int("StatusCode", res.StatusCode).
					Int64("Latency", time.Since(startTime).Milliseconds()).
					Msgf("Processed")

				fmt.Println(`{ "_aws": { "Timestamp": 1574109732004, "CloudWatchMetrics": [{ "Namespace": "lambda-function-metrics", "Dimensions": [ ["functionVersion"] ], "Metrics": [{ "Name": "time", "Unit": "Milliseconds" }] }] }, "functionVersion": "$LATEST", "time": 100, "requestId": "989ffbf8-9ace-4817-a57c-e4dd734019ee" }`)

				emf.New(emf.WithWriter(CWAGENT_CONNECTION)).Namespace("mtg").Metric("totalWins", 1500).Log()

				emf.New(emf.WithWriter(CWAGENT_CONNECTION)).Dimension("colour", "red").
					MetricAs("gameLength", 2, emf.Seconds).Log()

				emf.New(emf.WithWriter(CWAGENT_CONNECTION)).DimensionSet(
					emf.NewDimension("format", "edh"),
					emf.NewDimension("commander", "Muldrotha")).
					MetricAs("wins", 1499, emf.Count).Log()

				break // request was succesful
			}

		}

		time.Sleep(DEFAULT_WAIT)
	}
}
