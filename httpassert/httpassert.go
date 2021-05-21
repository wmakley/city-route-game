package httpassert

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func Success(t *testing.T, resp *httptest.ResponseRecorder) bool {
	if resp.Code != 200 {
		t.Errorf("Status code is not 200 (was %d)", resp.Code)
		return false
	}
	return true
}

func NotFound(t *testing.T, resp *httptest.ResponseRecorder) bool {
	if resp.Code != 404 {
		t.Errorf("Status code is not 404 (was %d)", resp.Code)
		return false
	}
	return true
}

func HtmlContentType(t *testing.T, resp *httptest.ResponseRecorder) bool {
	contentType := resp.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Error("Content-Type is not 'text/html; charset=utf-8'; was:", contentType)
		return false
	}
	return true
}

// JsonContentType Assert that the content-type header is "application/json; charset=utf-8"
// (Not providing the charset may cause issues for some clients.)
func JsonContentType(t *testing.T, resp *httptest.ResponseRecorder) bool {
	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Error("Content-Type is not 'application/json; charset=utf-8'; was:", contentType)
		return false
	}
	return true
}

func JavascriptContentType(t *testing.T, resp *httptest.ResponseRecorder) bool {
	contentType := resp.Header().Get("Content-Type")
	if contentType != "text/javascript; charset=utf-8" {
		t.Error("Content-Type is not 'text/javascript; charset=utf-8'; was:", contentType)
		return false
	}
	return true
}

// JsonArray Assert that the response content type is JSON, and the first character is "["
func JsonArray(t *testing.T, resp *httptest.ResponseRecorder) bool {
	if !JsonContentType(t, resp) {
		return false
	}

	if !strings.HasPrefix(resp.Body.String(), "[") {
		t.Error("body is not a JSON array")
		return false
	}

	return true
}

// JsonObject Assert that the response content type is JSON, and the first character is "{"
func JsonObject(t *testing.T, resp *httptest.ResponseRecorder) bool {
	if !JsonContentType(t, resp) {
		return false
	}

	if !strings.HasPrefix(resp.Body.String(), "{") {
		t.Error("body is not a JSON object")
		return false
	}

	return true
}
