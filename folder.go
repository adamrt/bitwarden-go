package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Folder struct {
	UUID     string `db:"uuid" json:"Id"`
	UserUUID string `db:"user_uuid" json:"-"`
	Name     string `db:"name"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"RevisionDate"`
}

func FoldersGetByUserUUID(db *sqlx.DB, folders *[]Folder, uuid string) error {
	return db.Select(folders, "SELECT * FROM folders WHERE user_uuid=$1", uuid)
}

func FolderGetByUserUUIDAndUUID(db *sqlx.DB, folder *Folder, userUUID, uuid string) error {
	return db.Get(folder, "SELECT * FROM folders WHERE user_uuid=$1 AND uuid=$2", userUUID, uuid)
}
