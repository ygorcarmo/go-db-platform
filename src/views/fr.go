package views

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
)

//go:embed templates/**
var assets embed.FS

var Templates map[string]*template.Template

func LoadTemplates() {
	Templates = make(map[string]*template.Template)

	templatesFolder, err := fs.Sub(assets, "templates")

	if err != nil {
		log.Fatal(err)
	}

	Templates["home"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "home.tmpl"))
	Templates["createDbUser"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbUser/createDbUser.tmpl"))
	Templates["deleteDbUser"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbUser/deleteDbUser.tmpl"))
	Templates["dbUserFormResponse"] = template.Must(template.ParseFS(templatesFolder, "dbUser/response.tmpl"))
	Templates["addDb"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "dbConfig/addDb.tmpl"))
	Templates["signIn"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "signIn.tmpl"))
	Templates["resetPassword"] = template.Must(template.ParseFS(templatesFolder, "base-layout.tmpl", "appUser/resetPassword.tmpl"))
}
