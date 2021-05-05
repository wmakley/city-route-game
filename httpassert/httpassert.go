package httpassert

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func Success(t *testing.T, resp *httptest.ResponseRecorder) {
	if resp.Code != 200 {
		t.Error("Status code is not 200")
	}
}

func NotFound(t *testing.T, resp *httptest.ResponseRecorder) {
	if resp.Code != 404 {
		t.Error("Status code is not 404")
	}
}

// Assert that the content-type header is "application/json; charset=utf-8"
// (Not providing the charset may cause issues for some clients.)
func JsonContentType(t *testing.T, resp *httptest.ResponseRecorder) {
	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Error("Content-Type is not 'application/json'")
	}
}

// Assert that the response content type is JSON, and the first character is "["
func JsonArray(t *testing.T, resp *httptest.ResponseRecorder) {
	JsonContentType(t, resp)

	if !strings.HasPrefix(resp.Body.String(), "[") {
		t.Error("body is not a JSON array")
	}
}

// Assert that the response content type is JSON, and the first character is "{"
func JsonObject(t *testing.T, resp *httptest.ResponseRecorder) {
	JsonContentType(t, resp)

	if !strings.HasPrefix(resp.Body.String(), "{") {
		t.Error("body is not a JSON object")
	}
}
