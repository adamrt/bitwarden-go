package main

import (
	"crypto/subtle"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	UUID          string     `db:"uuid" json:"Id"`
	Name          NullString `db:"name"`
	Email         string     `db:"email"`
	EmailVerified bool       `db:"email_verified"`
	Premium       bool       `db:"premium"`
	PasswordHash  string     `db:"password_hash"`
	PasswordHint  NullString `db:"password_hint" json:"MasterPasswordHint"`
	Key           string     `db:"key"`
	PrivateKey    NullString `db:"private_key"`
	PublicKey     NullString `db:"public_key"`
	TOTPSecret    NullString `db:"totp_secret"`
	SecurityStamp NullString `db:"security_stamp"`
	Culture       NullString `db:"culture"`

	// For Response, not DB
	TwoFactorEnabled bool
	Organizations    []string
	Object           string

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u User) CheckPassword(pw string) bool {
	pass := []byte(u.PasswordHash)
	guess := []byte(pw)
	// ConstantTimeCompare returns 1 instead of bool
	return subtle.ConstantTimeCompare(pass, guess) == 1
}

func UserGetByEmail(db *sqlx.DB, user *User, email string) error {
	return db.Get(user, "SELECT * FROM users WHERE email=$1", email)
}

func UserGetByUUID(db *sqlx.DB, user *User, uuid string) error {
	return db.Get(user, "SELECT * FROM users WHERE uuid=$1", uuid)
}
