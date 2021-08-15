package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type UserID struct {
	OID string
}

type User struct {
	email     string
	FirstName string `firstname:"required"`
	lastName  string
}

// func (t Router) AddUser2() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 			return
// 		}
// 		var user User
// 		err := json.NewDecoder(r.Body).Decode(&user)
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			return
// 		}

// 		oid, err := t.MongoDBClient.CreateUser(user)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		err = json.NewEncoder(w).Encode(UserID{OID: oid})
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			return
// 		}

// 	}
// }

func (t Router) AddUser2() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var user User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		err = user.Unmarshal(body)
		if err != nil {
			println(err.Error())
		}

	}
}

func (p *User) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	fields := reflect.ValueOf(p).Elem()
	for i := 0; i < fields.NumField(); i++ {
		yourpojectTags := fields.Type().Field(i).Tag.Get(fields.Field(i).String())
		if strings.Contains(yourpojectTags, "required") && fields.Field(i).IsZero() {
			return errors.New("required field is missing")
		}

	}

	return nil
}
