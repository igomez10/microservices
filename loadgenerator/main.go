package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"loadgenerator/cmd"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

var FileName = "./exampleConfig.json"

func main() {
	jsonFile, err := os.Open(FileName)
	if err != nil {
		log.Fatal().Msgf("failed to open file: %s", err)
	}
	fmt.Println("Successfully Opened exampleConfig.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal().Msgf("Failed to read file: %s", err)
	}

	var config cmd.LoadGeneratorConfig
	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatal().Msgf("Failed to unmarshal config: %s", err)
	}

	// fallback to stdout if cwagent is not reachable
	var metricOutput io.Writer = os.Stdout
	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:25888", time.Millisecond*10000)
		if err != nil {
			log.Error().Err(err).Msgf("failed to connect to cloudwatch agent, attempt %d", i)
			// time.Sleep(1 * time.Second)

			if i == 9 {
				log.Warn().Msgf("Exhausted all attempts (%d) to connect to cwagent: %s", i, err)
			}

			continue
		}
		metricOutput = conn
		defer conn.Close()
		break
	}

	cmd.Run(config, metricOutput)
}
