package server

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type templateHandlerFactory struct {
	tmpl          *template.Template
	Credentials   Credentials
	MongoDBClient mongoDBClient
}

func (rt Router) NewTemplateHandlerFactory(templateDirPath string) templateHandlerFactory {
	var tmpl *template.Template
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join(templateDirPath, "*.gohtml")))
	template.Must(tmpl.ParseGlob(filepath.Join(templateDirPath, "partials/*.gohtml")))

	return templateHandlerFactory{tmpl, rt.Credentials, rt.MongoDBClient}
}

type templateData struct {
	PageName  string
	MenuItems []MenuItem
}

func (t templateHandlerFactory) AuthHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		valid := t.Credentials.ValidTokenRequest(r)
		if valid {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		data := &templateData{
			PageName: "Auth",
		}
		err := t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.NotFound(w, r)
		}

		return
	}
}

func (t templateHandlerFactory) DashboardHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		if t.Credentials.ValidTokenRequest(r) {
			var pageName string
			if templateName == "dashboard.gohtml" {
				pageName = "Dashboard"
			}

			data := &templateData{
				PageName: pageName,
			}

			err := t.tmpl.ExecuteTemplate(w, templateName, data)
			if err != nil {
				http.NotFound(w, r)
			}
			return
		}

		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return

	}
}

func (t templateHandlerFactory) PublicHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		pageName := "Menu"

		menuItems, err := t.MongoDBClient.AllMenuItems()
		if err != nil {
			println(err.Error())
		}

		data := &templateData{
			PageName:  pageName,
			MenuItems: menuItems,
		}

		err = t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.NotFound(w, r)
			return
		}

	}
}

func (t templateHandlerFactory) AddMenuItemHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		valid := t.Credentials.ValidTokenRequest(r)
		if !valid {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		pageName := "Add Menu Item"

		data := &templateData{
			PageName: pageName,
		}

		err := t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.NotFound(w, r)
			return
		}

	}
}

func (t templateHandlerFactory) ManageItemHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		valid := t.Credentials.ValidTokenRequest(r)
		if !valid {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		pageName := "Menu Edit"

		menuItems, err := t.MongoDBClient.AllMenuItems()
		if err != nil {
			println(err.Error())
		}

		data := &templateData{
			PageName:  pageName,
			MenuItems: menuItems,
		}

		err = t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.NotFound(w, r)
			return
		}

	}
}

type EditMenuItem struct {
	PageName        string
	ImageSource     string
	ItemPrice       string
	ItemDescription string
	ItemName        string
	ItemIngredients string
	ID              string
}

func (t templateHandlerFactory) ManageEditHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		valid := t.Credentials.ValidTokenRequest(r)
		if !valid {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		objID := r.URL.Query().Get("item")

		item, err := t.MongoDBClient.MenuItemFromID(objID)
		if err != nil {
			println(err.Error())
			http.NotFound(w, r)
			return
		}

		pageName := "Edit " + item.ItemName

		data := &EditMenuItem{
			PageName:        pageName,
			ImageSource:     item.ImageSource,
			ItemPrice:       item.ItemPrice,
			ItemDescription: item.ItemDescription,
			ItemName:        item.ItemName,
			ItemIngredients: item.ItemIngredients,
			ID:              item.ID,
		}

		err = t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.NotFound(w, r)
			return
		}

	}
}
