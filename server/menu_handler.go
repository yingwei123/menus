package server

import (
	"encoding/json"
	"errors"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type MenuItem struct {
	ImageSource     string
	ItemPrice       string
	ItemDescription string
	ItemName        string
	ItemIngredients string
	ID              string
}

type ObjectID struct {
	OID string
}

func (rt Router) AddMenuItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		valid := rt.Credentials.ValidTokenRequest(r)
		if !valid {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, fileName, err := r.FormFile("image-upload-input")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		itemPrice := r.FormValue("item-price")
		ingredients := r.FormValue("item-ingredients")
		description := r.FormValue("item-description")
		itemName := r.FormValue("item-name")

		oid, err := rt.MenuAdd(fileName.Filename, itemPrice, itemName, description, "/public/uploads/"+itemName+".jpeg", ingredients, file)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(UserID{OID: oid})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		return
	}
}

func (rt Router) MenuAdd(fileName string, itemPrice string, itemName string, itemDescription string, imageSource string, itemIngredients string, file multipart.File) (string, error) {
	println(fileName)
	if strings.Contains(fileName, ".jpg") || strings.Contains(fileName, ".jpeg") {
		image, err := jpeg.Decode(file)
		f, err := os.Create("resources/public/uploads/" + itemName + ".jpeg")
		if err != nil {
			return "", err
		}
		defer f.Close()
		err = png.Encode(f, image)
		if err != nil {

			return "", err
		}
		menuItem := MenuItem{ImageSource: "/public/uploads/" + itemName + ".jpeg", ItemPrice: itemPrice, ItemDescription: itemDescription, ItemName: itemName, ItemIngredients: itemIngredients}

		oid, err := rt.MongoDBClient.CreateMenuItem(menuItem)
		if err != nil {

			return "", err
		}

		return oid, nil
	}

	if strings.Contains(fileName, ".png") {
		image, err := png.Decode(file)
		f, err := os.Create("resources/public/uploads/" + itemName + ".png")
		if err != nil {
			return "", err
		}
		defer f.Close()
		err = png.Encode(f, image)
		if err != nil {

			return "", err
		}
		menuItem := MenuItem{ImageSource: "/public/uploads/" + itemName + ".png", ItemPrice: itemPrice, ItemDescription: itemDescription, ItemName: itemName, ItemIngredients: itemIngredients}

		oid, err := rt.MongoDBClient.CreateMenuItem(menuItem)
		if err != nil {

			return "", err
		}

		return oid, nil
	}

	return "", errors.New("wrong extention")
}
