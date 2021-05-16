package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"

	"gorm.io/gorm"
)

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) (string, error) {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded, nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	return ip, nil
}

func SetCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func SetJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func SetHTMLContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func SetJavaScriptContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
}

// Write JSON to response body or panic
func MustEncode(w io.Writer, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

// Write JSON to response body or panic, and set content type
func MustReturnJson(w http.ResponseWriter, v interface{}) {
	SetJSONContentType(w)
	MustEncode(w, v)
}

func MustParse(filenames ...string) *template.Template {
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		panic(err)
	}
	return t
}

type ErrorJson struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func JsonInternalServerError(msg string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	body := ErrorJson{500, msg}
	MustEncode(w, body)
}

func JsonBadReqest(msg string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	body := ErrorJson{400, msg}
	MustEncode(w, body)
}

func JsonNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	body := ErrorJson{404, "404 Not Found"}
	MustEncode(w, body)
}

func JsonRecordNotFound(recordName string, key string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	msg := fmt.Sprintf("%s with ID %s not found!", recordName, key)
	body := ErrorJson{404, msg}
	MustEncode(w, body)
}

func JsonHandleDbErr(err error, w http.ResponseWriter, r *http.Request) {
	if err == gorm.ErrRecordNotFound {
		JsonNotFound(w, r)
		return
	} else {
		JsonInternalServerError(err.Error(), w, r)
	}
}

func TurbolinksVisit(url string, clearCache bool, w http.ResponseWriter, r *http.Request) {
	clearCacheStep := ""
	if clearCache {
		clearCacheStep = "Turbolinks.clearCache();"
	}

	SetJavaScriptContentType(w)
	fmt.Fprintf(w, `;(function(){%sTurbolinks.visit("%s");})();`, clearCacheStep, url)
}
