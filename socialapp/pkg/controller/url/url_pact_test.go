package url

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/require"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	urlClient "github.com/igomez10/microservices/urlshortener/generated/clients/go/client"
)

func TestURLApiService_GetUrlDataContract(t *testing.T) {
	pactDir, logDir := contractDirectoryPaths(t)

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "socialapp",
		Provider: "urlshortener",
		PactDir:  pactDir,
		LogDir:   logDir,
	})
	require.NoError(t, err)

	alias := "contract-alias"
	requestID := "contract-id-1234"
	expectedURL := "https://pact.socialapp/golden-url"

	interaction := pact.AddInteraction()
	testInteraction := interaction.Given("a url alias exists").
		UponReceiving("a request for url data").
		WithRequestPathMatcher(http.MethodGet, matchers.S(fmt.Sprintf("/v1/urls/%s/data", alias)),
			func(r *consumer.V2RequestBuilder) {
				r.Header("X-Request-ID", matchers.String(requestID))
			},
		).
		WillRespondWith(http.StatusOK,
			func(r *consumer.V2ResponseBuilder) {
				r.Header("Content-Type", matchers.String("application/json"))
				r.JSONBody(matchers.StructMatcher{
					"url":        matchers.String(expectedURL),
					"alias":      matchers.String(alias),
					"created_at": matchers.Timestamp(),
					"updated_at": matchers.Timestamp(),
					"deleted_at": matchers.Timestamp(),
				})
			},
		)

	err = testInteraction.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		return verifyGetUrlDataRequest(config, alias, requestID, expectedURL)
	})
	require.NoError(t, err)
}

func verifyGetUrlDataRequest(config consumer.MockServerConfig, alias, requestID, expectedURL string) error {
	service := newURLApiServiceForMock(config)

	ctx := contexthelper.SetRequestIDInContext(context.Background(), requestID)
	res, err := service.GetUrlData(ctx, alias)
	if err != nil {
		return err
	}

	if res.Code != http.StatusOK {
		return fmt.Errorf("unexpected status %d, expected %d", res.Code, http.StatusOK)
	}

	body, ok := res.Body.(openapi.Url)
	if !ok {
		return fmt.Errorf("invalid body type %T, expected %T", res.Body, openapi.Url{})
	}

	if body.Alias != alias {
		return fmt.Errorf("unexpected alias %q, expected %q", body.Alias, alias)
	}

	if body.Url != expectedURL {
		return fmt.Errorf("unexpected url %q, expected %q", body.Url, expectedURL)
	}

	return nil
}

func newURLApiServiceForMock(config consumer.MockServerConfig) *URLApiService {
	apiConfig := urlClient.NewConfiguration()
	apiConfig.Servers = urlClient.ServerConfigurations{
		{URL: fmt.Sprintf("http://%s:%d", config.Host, config.Port)},
	}
	return &URLApiService{
		Client: urlClient.NewAPIClient(apiConfig),
	}
}

func contractDirectoryPaths(t *testing.T) (string, string) {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "could not determine caller file")

	root := filepath.Join(filepath.Dir(file), "..", "..", "..")
	pactDir := filepath.Join(root, "contracts", "pacts")
	logDir := filepath.Join(root, "contracts", "logs")

	require.NoError(t, os.MkdirAll(pactDir, 0o755))
	require.NoError(t, os.MkdirAll(logDir, 0o755))

	return pactDir, logDir
}
