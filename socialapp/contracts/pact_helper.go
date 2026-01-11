package contracts

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/stretchr/testify/require"
)

type contractSpec struct {
	providerState    string
	description      string
	requestMethod    consumer.Method
	requestPath      string
	requestBuilder   func(*consumer.V2RequestBuilder)
	responseBuilder  func(*consumer.V2ResponseBuilder)
	responseStatus   int
	executionHandler func(consumer.MockServerConfig) error
}

func runConsumerContract(t *testing.T, consumerName, providerName string, spec contractSpec) {
	pactDir, logDir := contractDirectoryPaths(t)
	pact := newV2Pact(t, consumerName, providerName, pactDir, logDir)
	interaction := pact.AddInteraction()

	testInteraction := interaction.Given(spec.providerState).
		UponReceiving(spec.description).
		WithRequest(spec.requestMethod, spec.requestPath, func(r *consumer.V2RequestBuilder) {
			if spec.requestBuilder != nil {
				spec.requestBuilder(r)
			}
		}).
		WillRespondWith(spec.responseStatus, func(r *consumer.V2ResponseBuilder) {
			if spec.responseBuilder != nil {
				spec.responseBuilder(r)
			}
		})

	if spec.executionHandler != nil {
		require.NoError(t, testInteraction.ExecuteTest(t, spec.executionHandler))
	}
}

func newV2Pact(t *testing.T, consumerName, providerName, pactDir, logDir string) *consumer.V2HTTPMockProvider {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: consumerName,
		Provider: providerName,
		PactDir:  pactDir,
		LogDir:   logDir,
	})
	require.NoError(t, err)
	return pact
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
