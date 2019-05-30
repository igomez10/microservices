package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	portEnviron := os.Getenv("PORT")

	if portEnviron == "" {
		portEnviron = "8080"
	}

	port := flag.String("port", portEnviron, "port to listen")
	*port = ":" + *port
	flag.Parse()

	http.HandleFunc("/bitcoinPrice", func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("%+v\n", *r)

		price, err := getBitcoinPrice()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, price)
	})

	err := http.ListenAndServe(*port, nil)
	log.Fatal(err)

}

func getBitcoinPrice() (float64, error) {

	resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		return 0, fmt.Errorf("error requesting data from website: %+v", err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("invalid response from website %s", err)
	}

	var btcprice coindeskBTCResponse
	jsonError := json.Unmarshal(contents, &btcprice)
	if jsonError != nil {
		return 0, fmt.Errorf("invalid response from website or modified struct def - %s", jsonError)
	}

	return btcprice.Bpi.USD.RateFloat, nil

}
