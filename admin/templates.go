package admin

import (
	"html/template"
	"io"
	"log"
	"path"
)

const templateExtension = ".tmpl"

func ParseAndExecuteAdminTemplate(w io.Writer, shortPath string, data interface{}, extraTemplates ...string) error {
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
		log.Println("Template Parse Error:", err.Error())
		return err
	}

	err = t.ExecuteTemplate(w, primaryTemplateName, data)
	if err != nil {
		log.Println("Template Execution Error:", err.Error())
		return err
	}

	return nil
}

func TemplatePath(shortPath string) string {
	return templateRoot + "/admin/" + shortPath + templateExtension
}
