package main

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
)

func createCipherHandler(db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("createCipherHandler")

		device := Device{}
		err := DeviceGetByBearer(db, &device, r)
		if err != nil {
			JSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// TODO: will all cases fail above. can device be nil?
		if &device == nil {
			JSONError(w, http.StatusUnauthorized, "invalid bearer")
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		p := cipherPayload{}
		decoder := schema.NewDecoder()
		if err := decoder.Decode(&p, r.PostForm); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if p.Name == "" || p.Type == "" {
			log.Println("name and type are requeired")
			JSONError(w, http.StatusBadRequest, "name and type are required")
			return
		}

		if err := ParseCipherString(p.Name); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		folder := Folder{}
		if p.FolderUUID != "" {
			err := FolderGetByUserUUIDAndUUID(db, &folder, device.UserUUID, p.FolderUUID)
			if err != nil {
				JSONError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		cipher := Cipher{}
		cipher.UserUUID = device.UserUUID
		cipher.UpdateFromPayload(p)

		// TODO: Incomplete
	})
}
