package middleware

import (
	"city-route-game/util"
	"log"
	"net/http"
	"strings"
	"time"
)

// Automatically parse form data for POST, PUT, and PATCH requests
// Hopefully this handles both multipart form data and url-encoded?
func ParseFormData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {

			// Do nothing if content type is JSON
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {

				err := r.ParseMultipartForm(32 << 20) // 32 MB
				if err != nil {
					err = r.ParseForm()
					if err != nil {
						panic(err)
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Require all requests that can make changes to have the "X-Requested-With"
// header set. (This header may not be sent by cross-origin requests.)
func CSRFMitigation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
				log.Println("Aborting due to CSRF mitigation.")
				http.NotFound(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore uninteresting static files
		if strings.HasPrefix(r.RequestURI, "/favicon.ico") || strings.HasPrefix(r.RequestURI, "/robots.txt") {
			next.ServeHTTP(w, r)
			return
		}

		ip, err := util.GetIP(r)
		if err != nil {
			ip = "UNKOWN_IP"
		}

		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] [%s] %q %v\n", ip, r.Method, r.URL.String(), t2.Sub(t1))
	})
}

func JsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if w.Header().Get("Content-Type") == "" {
			util.SetJSONContentType(w)
		}
	})
}

func HtmlContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if w.Header().Get("Content-Type") == "" {
			util.SetHTMLContentType(w)
		}
	})
}

func PreventCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// Allow the handler to specify own Cache-Control header
		if r.Header.Get("Cache-Control") == "" {
			r.Header.Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
			r.Header.Set("Pragma", "no-cache")
			r.Header.Set("Expires", "0")
		}
	})
}

func CorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func NewIPWhiteListMiddleware(ips []string, verbose bool) func(http.Handler) http.Handler {
	ipWhitelistMap := make(map[string]bool, len(ips))
	for _, ip := range ips {
		ipWhitelistMap[ip] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := util.GetIP(r)
			if err != nil {
				log.Printf("Error getting request IP: %+v\n", err)
				ip = ""
			}

			_, ipFound := ipWhitelistMap[ip]

			if ipFound {
				next.ServeHTTP(w, r)
			} else {
				if verbose {
					log.Println("IP not found in whitelist")
				}
				http.NotFound(w, r)
			}
		})
	}
}
