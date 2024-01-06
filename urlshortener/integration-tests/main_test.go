package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	userClient "github.com/igomez10/microservices/socialapp/client"
	urlClient "github.com/igomez10/microservices/urlshortener/generated/clients/go/client"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	RENDER_SERVER_URL          = 0
	LOCALHOST_SERVER_URL       = 1
	LOCALHOST_DEBUG_SERVER_URL = 2

	CONTEXT_SERVER       int
	urlClnt              *urlClient.APIClient
	ENDPOINT_OAUTH_TOKEN string
)

// add setup function
func Setup() {
	//  set the endpoint for the oauth token
	testSetup := os.Getenv("TEST_SETUP")
	if testSetup == "" {
		testSetup = "LOCALHOST_DEBUG"
	}

	switch testSetup {
	case "LOCALHOST":
		CONTEXT_SERVER = LOCALHOST_SERVER_URL
		ENDPOINT_OAUTH_TOKEN = "http://localhost:8085/v1/oauth/token"
	case "LOCALHOST_DEBUG":
		CONTEXT_SERVER = LOCALHOST_DEBUG_SERVER_URL
		ENDPOINT_OAUTH_TOKEN = "http://localhost:8085/v1/oauth/token"
	default:
		CONTEXT_SERVER = RENDER_SERVER_URL
		ENDPOINT_OAUTH_TOKEN = "https://urlshortener.gomezignacio.com/v1/oauth/token"
	}
}

func TestMain(m *testing.M) {
	// run tests
	// add jitter
	if os.Getenv("ADD_TEST_JITTER") != "" {
		jitterInSeconds := uuid.New().ID() % 60
		log.Printf("Adding test jitter of %d seconds", jitterInSeconds)
		time.Sleep(time.Duration(jitterInSeconds))
	}

	Setup()
	code := m.Run()
	os.Exit(code)
}

func getHTTPClient() *http.Client {
	if os.Getenv("USE_PROXY") == "true" || true {
		proxyStr := "http://localhost:9091"
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return http.DefaultClient
		}
		// Setup http client with proxy to capture traffic
		httpClient := &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

		return httpClient
	}

	return http.DefaultClient
}

func TestURLLifeCycle(t *testing.T) {
	Setup()
	// setup url client
	urlClientConfiguration := urlClient.NewConfiguration()
	httpClient := getHTTPClient()
	urlClientConfiguration.HTTPClient = httpClient
	urlClnt = urlClient.NewAPIClient(urlClientConfiguration)

	// setup user client
	userClientConfiguration := userClient.NewConfiguration()
	userClientConfiguration.HTTPClient = httpClient
	userClnt := userClient.NewAPIClient(userClientConfiguration)

	// setup general context
	noAuthCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	noAuthCtx = context.WithValue(noAuthCtx, urlClient.ContextServerIndex, CONTEXT_SERVER)
	noAuthCtx = context.WithValue(noAuthCtx, userClient.ContextServerIndex, CONTEXT_SERVER)

	// create a user to run test on behalf of
	username1 := fmt.Sprintf("Test-%d1", time.Now().UnixNano())
	password := fmt.Sprintf("Password-%d1", time.Now().UnixNano())
	func() {
		createUsrReq := userClient.NewCreateUserRequest(username1, password, "FirstName_example", "LastName_example", username1)
		_, _, err := userClnt.UserAPI.CreateUser(noAuthCtx).CreateUserRequest(*createUsrReq).Execute()
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
	}()

	conf := clientcredentials.Config{
		ClientID:     username1,
		ClientSecret: password,
		Scopes:       []string{"shortly.url.create", "shortly.url.delete"},
		TokenURL:     ENDPOINT_OAUTH_TOKEN,
	}

	oauth2Ctx, err := getOuath2Context(noAuthCtx, conf)
	if err != nil {
		t.Fatalf("Error getting oauth2 context: %v", err)
	}

	// create url
	urlAPICtx := context.WithValue(oauth2Ctx, urlClient.ContextServerIndex, CONTEXT_SERVER)
	newURL := urlClient.NewURL("https://www.google.com", "google")
	_, r, err := urlClnt.URLAPI.CreateUrl(urlAPICtx).URL(*newURL).Execute()
	if err != nil {
		t.Errorf("Error when calling `URLAPI.CreateURL`: %v\n", err)
		t.Errorf("Full HTTP response: %v ", r)
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, r.StatusCode)
	}

	// get url
	getUrlRes, err := urlClnt.URLAPI.GetUrl(noAuthCtx, "google").Execute()
	if err != nil {
		t.Errorf("Error when calling `URLAPI.GetURL`: %v\n", err)
		t.Errorf("Full HTTP response: %v ", r)
		t.Fatalf("Error getting url: %v", err)
	}

	if getUrlRes.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, getUrlRes.StatusCode)
	}

	// delete url
	deleteUrlRes, err := urlClnt.URLAPI.DeleteUrl(urlAPICtx, "google").Execute()
	if err != nil {
		t.Errorf("Error when calling `URLAPI.DeleteURL`: %v\n", err)
		t.Errorf("Full HTTP response: %v ", r)
		t.Fatalf("Error deleting url: %v", err)
	}

	if deleteUrlRes.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, deleteUrlRes.StatusCode)
	}

	// get url
	getUrlRes, err = urlClnt.URLAPI.GetUrl(noAuthCtx, "google").Execute()
	if err == nil {
		t.Errorf("Expected error when calling `URLAPI.GetURL`, got none")
	}
	if getUrlRes.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, getUrlRes.StatusCode)
	}
}

func getOuath2Context(initialContext context.Context, config clientcredentials.Config) (context.Context, error) {
	tokenSource := config.TokenSource(initialContext)
	initialContext = context.WithValue(initialContext, urlClient.ContextOAuth2, tokenSource)

	return initialContext, nil
}
