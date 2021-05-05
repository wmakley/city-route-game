package admin

import (
	"city-route-game/util"
	"errors"
	"html/template"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type ErrorPage struct {
	StatusCode int
	Message    string
	Details    string
}

func genericNotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TemplatePath("error"))
	if err != nil {
		log.Printf("Showing generic error page due to template parse error: %+v\n", err)
		http.NotFound(w, r)
		return
	}

	errorPage := ErrorPage{
		StatusCode: 404,
		Message:    "Not Found",
		Details:    "The resource you were looking for was not found on this server. :(",
	}

	w.WriteHeader(http.StatusNotFound)
	util.MustExecute(t, w, &errorPage)
}

func internalServerError(err error, w http.ResponseWriter, r *http.Request) {
	log.Printf("Internal Server Error: %+v\n", err)

	t, err := template.ParseFiles(TemplatePath("error"))
	if err != nil {
		log.Printf("Showing generic error page due to template parse error: %+v\n", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	errorPage := ErrorPage{
		StatusCode: 500,
		Message:    "Internal Server Error",
		Details:    "Something went wrong. :(",
	}

	w.WriteHeader(http.StatusInternalServerError)
	util.MustExecute(t, w, &errorPage)
}

func handleDBErr(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		genericNotFound(w, r)
	} else {
		log.Printf("Database error: %+v\n", err)
		internalServerError(err, w, r)
	}
}
