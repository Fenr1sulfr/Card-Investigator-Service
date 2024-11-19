package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRecordFound = errors.New("record not found")
	ErrEditConflict  = errors.New("edit conflict")
)

type Models struct {
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions: PermissionModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}
