package server

import (
	"image/png"
	"net/http"
	"os"
)

func (rt Router) ImageHandler() http.HandlerFunc {
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

		file, _, err := r.FormFile("image-upload-input")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.FormValue("item-price")
		r.FormValue("item-ingredients")
		r.FormValue("item-description")
		itemName := r.FormValue("item-name")

		// buffer := make([]byte, 512)

		// _, err = file.Read(buffer)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// mimeType := http.DetectContentType(buffer)
		// println(mimeType)

		// if mimeType == "image/png" {
		image, err := png.Decode(file)
		f, err := os.Create("resources/public/uploads/" + itemName + ".png")
		if err != nil {
			// Handle error
		}
		defer f.Close()

		// Encode to `PNG` with `DefaultCompression` level
		// then save to file
		err = png.Encode(f, image)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// }

		w.WriteHeader(http.StatusOK)

		return
	}
}
