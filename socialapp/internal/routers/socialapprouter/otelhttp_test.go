package socialapprouter

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TestOTelHTTPMiddlewareDoesNotWriteHeaderTwice(t *testing.T) {
	var serverLogs bytes.Buffer
	handler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}), "test")

	server := httptest.NewUnstartedServer(handler)
	server.Config.ErrorLog = log.New(&serverLogs, "", 0)
	server.Start()
	t.Cleanup(server.Close)

	response, err := server.Client().Get(server.URL)
	require.NoError(t, err)
	defer response.Body.Close()

	_, err = io.Copy(io.Discard, response.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotContains(t, serverLogs.String(), "superfluous response.WriteHeader")
}
