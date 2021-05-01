package admin

import (
	"html/template"
	"io"
	"log"
	"path"
)

func ParseAndExecuteAdminTemplate(w io.Writer, shortPath string, data interface{}, extraTemplates ...string) error {
	primaryTemplateName := path.Base(shortPath) + ".tmpl"
	primaryTemplateFullPath := toFullPath(shortPath)

	allTemplatesToParse := make([]string, 2+len(extraTemplates))
	allTemplatesToParse[0] = layoutTemplatePath
	allTemplatesToParse[1] = primaryTemplateFullPath
	for index, path := range extraTemplates {
		allTemplatesToParse[index+2] = toFullPath(path)
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

const layoutTemplatePath = "./templates/admin/layout.tmpl"

func toFullPath(shortPath string) string {
	return "./templates/admin/" + shortPath + ".tmpl"
}
