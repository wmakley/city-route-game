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
	Data interface{}
}

func NewPageWithData(data interface{}) Page {
	return &PageWithData{
		Data: data,
	}
}

func (p *PageWithData) AssetHost() string {
	return assetHost
}

const templateExtension = ".tmpl"

func ParseAndExecuteAdminTemplate(w io.Writer, shortPath string, data *Page, extraTemplates ...string) error {
	primaryTemplateFullPath := TemplatePath(shortPath)
	primaryTemplateName := path.Base(primaryTemplateFullPath)

	allTemplatesToParse := make([]string, 2+len(extraTemplates))
	allTemplatesToParse[0] = TemplatePath("layout")
	allTemplatesToParse[1] = primaryTemplateFullPath
	for index, path := range extraTemplates {
		allTemplatesToParse[index+2] = TemplatePath(path)
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

// Store an executed template in a buffer before writing it to "w".
// This allows you to gracefully recover from a template execution failure
// during an http request, rather than send partial HTML.
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

func TemplatePath(shortPath string) string {
	hasExtension := strings.HasSuffix(shortPath, templateExtension)
	if hasExtension {
		return templateRoot + "/admin/" + shortPath
	} else {
		return templateRoot + "/admin/" + shortPath + templateExtension
	}
}
