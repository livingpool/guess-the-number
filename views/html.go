package views

import (
	"embed"
	"html/template"
	"io"
)

// The go embed directive statement must be outside of function body
// embed only works for current & sub modules (https://github.com/golang/go/issues/46056)

//go:embed assets/* css/* html/* scripts/*
var StaticFiles embed.FS

type TemplatesRepository interface {
	Render(w io.Writer, name string, data interface{}) error
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("./views/html/*.tmpl")),
	}
}
