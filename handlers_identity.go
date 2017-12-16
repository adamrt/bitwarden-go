package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
)

func tokenHandler(db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("tokenHandler")
		if err := r.ParseForm(); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		p := tokenPayload{}
		decoder := schema.NewDecoder()
		if err := decoder.Decode(&p, r.PostForm); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		device := Device{}
		user := User{}

		log.Printf("GrantType: %s", p.GrantType)
		if p.GrantType == "refresh_token" {
			if p.RefreshToken == "" {
				JSONError(w, http.StatusBadRequest, "Invalid refresh_token")
				return
			}

			err := DeviceGetByRefreshToken(db, &device, p.RefreshToken)
			if err != nil {
				JSONError(w, http.StatusBadRequest, "Invalid refresh_token")
				return
			}

			if err := UserGetByUUID(db, &user, device.UserUUID); err != nil {
				JSONError(w, http.StatusBadRequest, err.Error())
				return
			}

		} else if p.GrantType == "password" {
			if p.Scope != "api offline_access" {
				log.Println("scope not supported")
				JSONError(w, http.StatusBadRequest, "scope not supported")
				return
			}

			if err := UserGetByEmail(db, &user, p.Username); err != nil {
				log.Println(err.Error())
				JSONError(w, http.StatusBadRequest, err.Error())
				return
			}

			if !user.CheckPassword(p.Password) {
				log.Println("invalid password")
				JSONError(w, http.StatusBadRequest, "invalid password")
				return
			}

			if err := DeviceGetByUUID(db, &device, p.DeviceIdentifier); err != nil {
				log.Println(err.Error())
				JSONError(w, http.StatusBadRequest, err.Error())
				return
			}

			if device.UserUUID != user.UUID {
				query := `DELETE FROM devices WHERE uuid = :uuid`
				if _, err := db.NamedExec(query, device); err != nil {
					log.Println(err.Error())
					JSONError(w, http.StatusBadRequest, "Unknown error")
					return
				}
				device = Device{}
			}

			deviceTypeInt, err := strconv.Atoi(p.DeviceType)
			if err != nil {
				JSONError(w, http.StatusBadRequest, "device_type not integer")
				return
			}

			device.Type = deviceTypeInt
			device.Name = p.DeviceName
			if p.DevicePushToken != "" {
				device.PushToken = NullString{sql.NullString{String: p.DevicePushToken, Valid: true}}
			}

		} else {
			log.Println("grant_type must be 'refresh_token' or 'password'")
			JSONError(w, http.StatusBadRequest, "grant_type must be 'refresh_token' or 'password'")
			return
		}

		device.RegenerateTokens(db)

		query := `UPDATE devices SET type = :type, name = :name, push_token = :push_token`
		if _, err := db.NamedExec(query, device); err != nil {
			log.Println(err.Error())
			JSONError(w, http.StatusBadRequest, "Unknown error")
			return
		}

		// TODO: expires
		// expiresIn := time.Now().Sub(device.TokenExpiresAt)
		response := tokenResponse{
			AccessToken:  device.AccessToken,
			ExpiresIn:    360000,
			TokenType:    "Bearer",
			RefreshToken: device.RefreshToken,
			Key:          user.Key,
		}
		JSONResponse(w, response)
	})

}
