package admin

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"path"
	"strings"
)

type Page interface {
	AssetHost() string
}

type PageWithData struct {
	assetHost string
	Data interface{}
}

func (c Controller)NewPageWithData(data interface{}) Page {
	return &PageWithData{
		assetHost: c.AssetHost,
		Data: data,
	}
}

func (p PageWithData) AssetHost() string {
	return p.assetHost
}

const templateExtension = ".tmpl"

func (c Controller)ParseAndExecuteAdminTemplate(w io.Writer, shortPath string, data *Page, extraTemplates ...string) error {
	primaryTemplateFullPath := c.TemplatePath(shortPath)
	primaryTemplateName := path.Base(primaryTemplateFullPath)

	allTemplatesToParse := make([]string, 2+len(extraTemplates))
	allTemplatesToParse[0] = c.TemplatePath("layout")
	allTemplatesToParse[1] = primaryTemplateFullPath
	for index, path := range extraTemplates {
		allTemplatesToParse[index+2] = c.TemplatePath(path)
	}

	t, err := template.ParseFiles(allTemplatesToParse...)
	if err != nil {
		log.Printf("Template Parse Error: %+v\n", err)
		return err
	}

	err = ExecuteTemplateBuffered(t, w, primaryTemplateName, data)
	if err != nil {
		log.Printf("Template Execution Error: %+v\n", err)
		return err
	}

	return nil
}

// ExecuteTemplateBuffered Store an executed template in a buffer before writing it to "w".
// This allows you to gracefully recover from a template execution failure
// during a request, rather than send partial HTML.
func ExecuteTemplateBuffered(t *template.Template, w io.Writer, templateName string, data interface{}) error {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024)) // 1K buffer (should hold output of smallest templates easily)

	err := t.ExecuteTemplate(buffer, templateName, data)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (c Controller)TemplatePath(shortPath string) string {
	hasExtension := strings.HasSuffix(shortPath, templateExtension)
	if hasExtension {
		return c.TemplateRoot + "/admin/" + shortPath
	} else {
		return c.TemplateRoot + "/admin/" + shortPath + templateExtension
	}
}
