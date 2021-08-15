package server

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Credentials struct {
	UserName string
	Password string
	Token    string
}

func (cred Credentials) ValidTokenRequest(r *http.Request) bool {
	c, err := r.Cookie("cookie")
	if err != nil {
		return false
	}
	if c.Value == cred.Token {
		return true
	}

	return false
}

type Router struct {
	ResourcesPath string
	ServerURL     string
	MongoDBClient mongoDBClient
	Credentials   Credentials
}

type mongoDBClient interface {
	CreateMenuItem(menuItem MenuItem) (string, error)
	AllMenuItems() ([]MenuItem, error)
	MenuItemFromID(id string) (MenuItem, error)
}

func (rt Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	template := rt.NewTemplateHandlerFactory(filepath.Join(rt.ResourcesPath, "templates"))

	router.Handle("/dashboard", template.DashboardHandler("dashboard.gohtml"))
	router.Handle("/menu", template.PublicHandler("menu.gohtml"))
	router.Handle("/", template.PublicHandler("menu.gohtml"))
	router.Handle("/auth", template.AuthHandler("auth.gohtml"))
	router.Handle("/login", rt.LoginHandler())
	router.Handle("/dashboard/add-item", template.AddMenuItemHandler("addItem.gohtml"))
	router.Handle("/dashboard/add-menu-item", rt.AddMenuItem())
	router.Handle("/dashboard/edit-menu-item", template.ManageItemHandler("manageItems.gohtml"))
	router.Handle("/dashboard/edit-menu-item/edit", template.ManageEditHandler("edit-item.gohtml"))

	n := negroni.New()
	n.Use(negroni.NewLogger())

	faviconMiddleware := negroni.NewStatic(http.Dir(filepath.Join(rt.ResourcesPath, "favicon.ico")))
	faviconMiddleware.Prefix = "/favicon.ico"
	n.Use(faviconMiddleware)

	publicMiddleware := negroni.NewStatic(http.Dir(filepath.Join(rt.ResourcesPath, "public")))
	publicMiddleware.Prefix = "/public"

	n.Use(publicMiddleware)

	n.UseHandler(router)
	n.ServeHTTP(w, r)
}
