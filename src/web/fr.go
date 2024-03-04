package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
)

//go:embed templates/** assets/**
var assets embed.FS

var Templates map[string]*template.Template

func LoadTemplates() {
	Templates = make(map[string]*template.Template)

	templatesFolder, err := fs.Sub(assets, "templates")

	if err != nil {
		log.Fatal(err)
	}

	Templates["home"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "home.tmpl"))
	Templates["createUserForm"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbUser/createDbUserForm.tmpl"))
	Templates["deleteUserForm"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbUser/deleteDbUserForm.tmpl"))
	Templates["dbUserFormResponse"] = template.Must(template.ParseFS(templatesFolder, "dbUser/response.tmpl"))
	Templates["addDbForm"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbConfig/addDbForm.tmpl"))
	Templates["signIn"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "signIn.tmpl"))
}
