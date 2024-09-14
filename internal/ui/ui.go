package ui

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates
var templates embed.FS

var views map[string]*template.Template

func LoadTemplates() {
	views = make(map[string]*template.Template)

	layouts := template.Must(template.ParseFS(templates, "templates/layouts/*.html.tmpl"))

	dir, _ := templates.ReadDir("templates")

	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		view := template.Must(layouts.Clone())
		view = template.Must(view.ParseFS(templates, "templates/"+f.Name()))
		views[f.Name()] = view
	}
}

func RenderView(view string, w io.Writer, data interface{}) error {
	return views[view].ExecuteTemplate(w, view, data)
}
