package cmd

import (
	"net"
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
		emfLogger := emf.New(emf.WithWriter(CWAGENT_CONNECTION), emf.WithLogGroup("ecsloadgenerator")).
			Namespace("loadgenerator").
			Property("URL", reqConfig.URL).
			Property("Method", reqConfig.Method)

		req, err := http.NewRequest(reqConfig.Method, reqConfig.URL, nil)
		if err != nil {
			log.Err(err).Msg("Failed to create request")
		}

		startTime := time.Now()
		res, err := httpClient.Do(req)
		finishTime := time.Now()
		if err != nil ||
			res.StatusCode == http.StatusTooManyRequests ||
			res.StatusCode >= http.StatusInternalServerError {

			log.Err(err).
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
		emfLogger.
			Property("StatusCode", strconv.Itoa(res.StatusCode)).
			MetricAs("Latency", latency, emf.Milliseconds).
			Log()

		time.Sleep(DEFAULT_WAIT)
	}

}
