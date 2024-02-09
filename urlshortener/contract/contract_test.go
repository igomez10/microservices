package contract

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/igomez10/microservices/urlshortener/generated/clients/go/client"
	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

func TestClientPact_Local(t *testing.T) {
	// initialize PACT DSL
	pact := dsl.Pact{
		Consumer: "example-client",
		Provider: "example-server",
	}

	// setup a PACT Mock Server
	pact.Setup(true)

	t.Run("get url by alias", func(t *testing.T) {
		alias := "example"

		pact.
			AddInteraction().                        // specify PACT interaction
			Given("URL exists in database").         // specify provider state
			UponReceiving("a request to get a URL"). // specify provider state
			WithRequest(dsl.Request{                 // specify expected request
				Method: http.MethodGet,
				// specify matching for endpoint
				Path: dsl.Term("/v1/urls/"+alias+"/data", "/v1/urls/.*/data"),
			}).
			WillRespondWith(dsl.Response{ // specify minimal expected response
				Status: 200,
				Body: dsl.Like(server.Url{
					Alias:     "example",
					Url:       "http://example.com",
					CreatedAt: time.Now().Add(-time.Hour),
					UpdatedAt: time.Now(),
				}),
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
			})

		// verify interaction on client side
		err := pact.Verify(func() error {
			// specify host anf post of PACT Mock Server as actual server
			host := fmt.Sprintf("%s:%d", pact.Host, pact.Server.Port)

			// execute function
			client := NewClient(host)
			url, res, err := client.URLAPI.GetUrlData(context.Background(), alias).Execute()
			if err != nil {
				return fmt.Errorf("error on client: %s", err)
			}

			// check if url is not nil and alias is expected
			if url == nil || url.Alias != alias {
				return errors.New("url is not expected")
			}

			// check if actual response is equal to expected
			if res.StatusCode != 200 {
				return errors.New("status code is not expected")
			}

			return nil
		})

		if err != nil {
			t.Fatal(err)
		}

		// Publish the Pacts...
		func() {
			p := dsl.Publisher{}

			fmt.Println("Publishing Pact files to broker", os.Getenv("PACT_DIR"), os.Getenv("PACT_BROKER_URL"))
			err := p.Publish(types.PublishRequest{
				PactURLs:        []string{filepath.FromSlash("./pacts/example-client-example-server.json")},
				PactBroker:      fmt.Sprintf("http://localhost:9292"),
				ConsumerVersion: "1.0.0",
				Tags:            []string{"master"},
				BrokerUsername:  "pactbroker",
				BrokerPassword:  "pactbroker",
			})

			if err != nil {
				log.Println("ERROR: ", err)
				os.Exit(1)
			}
		}()
	})

	// write Contract into file
	if err := pact.WritePact(); err != nil {
		t.Fatal(err)
	}

	// stop PACT mock server
	pact.Teardown()
}

func NewClient(serverURL string) *client.APIClient {
	configuration := client.NewConfiguration()

	os.Setenv("HTTP_PROXY", "http://localhost:9091")
	os.Setenv("HTTPS_PROXY", "http://localhost:9091")

	proxyURL, err := url.Parse("http://localhost:9091")
	if err != nil {
		panic(err)
	}

	configuration.HTTPClient = &http.Client{
		//proxy
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	configuration.Host = serverURL
	configuration.Scheme = "http"
	apiClient := client.NewAPIClient(configuration)
	return apiClient
}
