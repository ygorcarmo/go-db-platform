package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
)

//go:embed templates assets
var assets embed.FS

var Templates map[string]*template.Template

func LoadTemplates() {
	Templates = make(map[string]*template.Template)

	templatesFolder, err := fs.Sub(assets, "templates")

	if err != nil {
		log.Fatal(err)
	}

	Templates["home"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "home.tmpl"))
	Templates["createUser"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "createUser.tmpl"))
	Templates["deleteUser"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "deleteUser.tmpl"))
	Templates["config"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "config.tmpl"))
}
