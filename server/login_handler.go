package server

import (
	"encoding/json"
	"net/http"
)

func (t Router) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var user Credentials
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if user.UserName == t.Credentials.UserName && user.Password == t.Credentials.Password {
			cookie := &http.Cookie{
				Name:  "cookie",
				Value: t.Credentials.Token,
			}
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
