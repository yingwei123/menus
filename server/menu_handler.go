package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
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
	if strings.Contains(fileName, ".jpg") || strings.Contains(fileName, ".jpeg") {
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return "", err
		}

		image, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
		newImage := resize.Resize(300, 300, image, resize.Lanczos3)

		f, err := os.Create("resources/public/uploads/" + itemName + ".jpeg")
		if err != nil {
			return "", err
		}
		defer f.Close()
		err = jpeg.Encode(f, newImage, nil)
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
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return "", err
		}

		image, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
		newImage := resize.Resize(300, 300, image, resize.Lanczos3)

		f, err := os.Create("resources/public/uploads/" + itemName + ".png")
		if err != nil {
			return "", err
		}
		defer f.Close()
		err = png.Encode(f, newImage)
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
