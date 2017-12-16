package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

const DefaultTokenValidity = (60 * 60 * time.Minute)

type Device struct {
	UUID           string     `db:"uuid" json:"Id"`
	UserUUID       string     `db:"user_uuid"`
	Name           string     `db:"name"`
	Type           int        `db:"type"`
	PushToken      NullString `db:"push_token"`
	AccessToken    string     `db:"access_token"`
	RefreshToken   string     `db:"refresh_token"`
	TokenExpiresAt time.Time  `db:"token_expires_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

func (d *Device) RegenerateTokens(db *sqlx.DB) error {
	if d.RefreshToken == "" {
		token, err := GenerateRandomStringURLSafe(64)
		if err != nil {
			return err
		}
		d.RefreshToken = token
	}

	d.TokenExpiresAt = time.Now().Add(DefaultTokenValidity)

	user := User{}
	err := UserGetByUUID(db, &user, d.UserUUID)
	if err != nil {
		return err
	}

	nbf := time.Now().Add(time.Duration(-2) * time.Minute)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"nbf":            nbf,
			"exp":            d.TokenExpiresAt,
			"iss":            "/identity",
			"premium":        user.Premium,
			"email":          user.Email,
			"email_verified": user.EmailVerified,
			"sstamp":         user.SecurityStamp,
			"device":         d.UUID,
			"scope":          []string{"api", "offline_access"},
			"amr":            []string{"Application"},
		},
	)

	log.Printf("%v", token)

	// TODO: fix types
	// d.AccessToken = token

	return nil
}

func DeviceGetByUUID(db *sqlx.DB, device *Device, uuid string) error {
	return db.Get(device, "SELECT * FROM devices WHERE uuid=$1", uuid)
}

func DeviceGetByAccessToken(db *sqlx.DB, device *Device, token string) error {
	return db.Get(device, "SELECT * FROM devices WHERE access_token=$1", token)
}

func DeviceGetByRefreshToken(db *sqlx.DB, device *Device, token string) error {
	return db.Get(device, "SELECT * FROM devices WHERE refresh_token=$1", token)
}

func DeviceGetByBearer(db *sqlx.DB, device *Device, r *http.Request) error {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	err := DeviceGetByAccessToken(db, device, token)
	if err != nil {
		return err
	}

	if device.TokenExpiresAt.Before(time.Now()) {
		// TODO: fix this
		// return errors.New("invalid bearer")
		// return nil
	}

	return nil
}
