package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loadgenerator/cmd"
	"log"
	"os"
)

var FileName = "./exampleConfig.json"

func main() {

	jsonFile, err := os.Open(FileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalf("failed to open file: %+v", err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	fmt.Println("Successfully Opened exampleConfig.json")

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config cmd.LoadGeneratorConfig
	// we unmarshal our byteArray which contains our
	// jsonFile's content into which we defined above
	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatalf("Failed to unmarshal file: %+v", err)
	}

	cmd.Run(config)
}
