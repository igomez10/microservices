package contexthelper

import (
	"net/http"
	"testing"
)

func TestRequestIDLifecycle(t *testing.T) {
	r := &http.Request{}
	expectedRequestID := "1234"
	r = r.WithContext(SetRequestIDInContext(r.Context(), &expectedRequestID))

	actualRequestID := GetRequestIDInContext(r.Context())
	if *actualRequestID != expectedRequestID {
		t.Errorf("expected request ID to be %q, got %q", expectedRequestID, *actualRequestID)
	}

	// Now change the request ID and make sure it's updated
	expectedRequestID = "5678"
	r = r.WithContext(SetRequestIDInContext(r.Context(), &expectedRequestID))
	actualRequestID = GetRequestIDInContext(r.Context())
	if *actualRequestID != expectedRequestID {
		t.Errorf("expected request ID to be %q, got %q", expectedRequestID, *actualRequestID)
	}
}

func TestPatternLifecycle(t *testing.T) {
	r := &http.Request{}
	expectedPattern := "/users/id"
	r = r.WithContext(SetRequestPatternInContext(r.Context(), &expectedPattern))

	actualPattern := GetRequestPatternInContext(r.Context())
	if *actualPattern != expectedPattern {
		t.Errorf("expected pattern to be %q, got %q", expectedPattern, *actualPattern)
	}

	// Now change the pattern and make sure it's updated
	expectedPattern = "/users/id/username"
	r = r.WithContext(SetRequestPatternInContext(r.Context(), &expectedPattern))
	actualPattern = GetRequestPatternInContext(r.Context())
	if *actualPattern != expectedPattern {
		t.Errorf("expected pattern to be %q, got %q", expectedPattern, *actualPattern)
	}
}
