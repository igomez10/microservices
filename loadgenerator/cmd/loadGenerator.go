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

const DEFAULT_WAIT = 500 * time.Millisecond
const MAX_RETRY = 5

var httpClient = http.Client{Timeout: 10 * time.Second}

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

				log.Debug().RawJSON("_aws", []byte(`{ "Version": "0", "Type": "Cluster", "ClusterName": "loadgeneratorcluster", "Timestamp": 1650587280000, "TaskCount": 1, "ContainerInstanceCount": 0, "ServiceCount": 1, "CloudWatchMetrics": [{ "Namespace": "microservices/loadgenerator", "Metrics": [{ "Name": "TaskCount", "Unit": "Count" }, { "Name": "ContainerInstanceCount", "Unit": "Count" }, { "Name": "ServiceCount", "Unit": "Count" } ], "Dimensions": [ [ "ClusterName" ] ] }] }`)).Msg("testing")

				break // request was succesful
			}

		}

		time.Sleep(DEFAULT_WAIT)
	}
}
