package admin

import (
	"city-route-game/internal/app"
	"city-route-game/util"
	"errors"
	"html/template"
	"log"
	"net/http"
)

type ErrorPage struct {
	AssetHost  string
	StatusCode int
	Message    string
	Details    string
}

func (c Controller)GenericNotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(c.TemplatePath("error.tmpl"))
	if err != nil {
		log.Printf("Showing generic error page due to template parse error: %+v\n", err)
		http.NotFound(w, r)
		return
	}

	errorPage := ErrorPage{
		AssetHost:  c.AssetHost,
		StatusCode: 404,
		Message:    "Not Found",
		Details:    "The resource you were looking for was not found on this server. :(",
	}

	w.WriteHeader(http.StatusNotFound)
	util.SetHTMLContentType(w)
	err = ExecuteTemplateBuffered(t, w, "error.tmpl", &errorPage)
	if err != nil {
		panic(err)
	}
}

func (c Controller)InternalServerError(err error, w http.ResponseWriter, r *http.Request) {
	log.Printf("Internal Server Error: %+v\n", err)

	t, err := template.ParseFiles(c.TemplatePath("error.tmpl"))
	if err != nil {
		log.Printf("Showing generic error page due to template parse error: %+v\n", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	errorPage := ErrorPage{
		AssetHost:  c.AssetHost,
		StatusCode: 500,
		Message:    "Internal Server Error",
		Details:    "Something went wrong. :(",
	}

	w.WriteHeader(http.StatusInternalServerError)
	util.SetHTMLContentType(w)
	err = ExecuteTemplateBuffered(t, w, "error.tmpl", &errorPage)
	if err != nil {
		panic(err)
	}
}

func (c Controller)HandleServiceError(err error, w http.ResponseWriter, r *http.Request) {
	// TODO: switch on accept header to render HTML vs JSON
	if errors.Is(app.RecordNotFound{}, err) {
		c.GenericNotFound(w, r)
	} else if errors.Is(app.ErrInvalidIDString{}, err) {
		log.Println(err.Error())
		c.GenericNotFound(w, r)
	} else {
		log.Printf("Unknown error: %+v\n", err)
		c.InternalServerError(err, w, r)
	}
}
