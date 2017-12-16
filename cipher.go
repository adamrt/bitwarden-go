package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	typeLogin = 1
	typeNote  = 2
	typeCard  = 3
)

type Field struct {
	Type  int
	Name  string
	Value string
}

type CipherData struct {
	Uri      string
	Username string
	Password string
	Totp     *string
	Name     string
	Notes    string
	Fields   *[]Field
}

type Cipher struct {
	UUID             string     `db:"uuid" json:"Id"`
	UserUUID         string     `db:"user_uuid" json:"-"`
	FolderUUID       NullString `db:"folder_uuid" schema:"FolderId"`
	OrganizationUUID NullString `db:"organization_uuid"`
	Type             int        `db:"type" schema:"type"`
	Data             string     `db:"data"`
	Favorite         bool       `db:"favorite"`
	Attachments      *string    `db:"attachments"`

	Edit                bool
	OrganizationId      *string
	OrganizationUseTotp bool
	Object              string

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"RevisionDate"`
}

func (c *Cipher) MarshalJSON() ([]byte, error) {

	var cd CipherData
	err := json.Unmarshal([]byte(c.Data), &cd)
	if err != nil {
		log.Printf("BAD ONE: %+v", c.Data)
		return []byte{}, err
	}

	type Alias Cipher
	return json.Marshal(&struct {
		Data CipherData
		*Alias
	}{
		Data:  cd,
		Alias: (*Alias)(c),
	})
}

func (c *Cipher) UpdateFromPayload(p cipherPayload) error {
	c.FolderUUID = NullString{sql.NullString{String: p.FolderUUID, Valid: true}}
	c.OrganizationUUID = NullString{sql.NullString{String: p.OrganizationUUID, Valid: true}}
	c.Favorite = p.Favorite
	t, err := strconv.Atoi(p.Type)
	if err != nil {
		log.Printf("Atoi failed in UpdateFromPayload: %v", err)
		return err
	}
	c.Type = t

	cData := CipherData{}
	cData.Name = p.Name

	// TODO: Incomplete
	if t == typeLogin {
	} else if t == typeCard {
	} else if t == typeNote {
	}

	return nil
}

func CiphersGetByUserUUID(db *sqlx.DB, ciphers *[]Cipher, uuid string) error {
	return db.Select(ciphers, "SELECT * FROM ciphers WHERE user_uuid=$1", uuid)
}

func ParseCipherString(s string) error {
	// TODO: add proper regex
	matched, err := regexp.MatchString("foo.*", s)
	if err != nil || matched == false {
		return errors.New("invalid cipher string")
	}
	return nil
}
