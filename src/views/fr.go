package views

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
	Templates["createUserPage"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbUser/createDbUserPage.tmpl"))
	Templates["deleteUserPage"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbUser/deleteDbUserPage.tmpl"))
	Templates["dbUserFormResponse"] = template.Must(template.ParseFS(templatesFolder, "dbUser/response.tmpl"))
	Templates["addDbPage"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbConfig/addDbPage.tmpl"))
	Templates["signIn"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "signIn.tmpl"))
}
