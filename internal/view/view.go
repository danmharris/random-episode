package view

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strings"
)

type Views map[string]*template.Template

//go:embed templates
var templates embed.FS

func ParseViews() Views {
	views := make(Views)

	layout, err := templates.ReadFile("templates/layouts/main.html")
	if err != nil {
		panic(err)
	}

	dir, err := templates.ReadDir("templates")
	if err != nil {
		panic(err)
	}

	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		view, err := templates.ReadFile(fmt.Sprintf("templates/%s", f.Name()))
		if err != nil {
			panic(err)
		}

		var sb strings.Builder
		sb.Write(layout)
		sb.Write(view)

		t := template.Must(template.New(f.Name()).Parse(sb.String()))
		views[f.Name()] = t
	}

	return views
}

func (v Views) RenderView(w io.Writer, name string, data any) {
	file := fmt.Sprintf("%s.html", name)
	v[file].ExecuteTemplate(w, "base", data)
}
