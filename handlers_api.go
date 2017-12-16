package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func syncHandler(db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var device Device
		if err := DeviceGetByBearer(db, &device, r); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		var user User
		if err := UserGetByUUID(db, &user, device.UserUUID); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		var folders []Folder
		if err := FoldersGetByUserUUID(db, &folders, user.UUID); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		var ciphers []Cipher
		if err := CiphersGetByUserUUID(db, &ciphers, user.UUID); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp := syncResponse{
			Profile: user,
			Folders: folders,
			Ciphers: ciphers,
			Object:  "sync",
			Domains: struct {
				EquivalentDomains       *[]string
				GlobalEquivalentDomains []string
				Object                  string
			}{
				EquivalentDomains:       nil,
				GlobalEquivalentDomains: []string{},
				Object:                  "domains",
			},
		}
		JSONResponse(w, resp)
	})
}

func registerHandler(db *sqlx.DB, allowSignups bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: ensure POST

		if !allowSignups {
			JSONError(w, http.StatusBadRequest, "Signups are not permitted")
			return

		}

		var j registerJSON
		err := json.NewDecoder(r.Body).Decode(&j)
		if err != nil {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: Incomplete
	})
}
