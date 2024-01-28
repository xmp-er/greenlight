package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found") //for when using get method there is no entry
	ErrEditConflict   = errors.New("edit conflict")    //to prevent race condition, if same version but two different goroutines then after first goroutine, this will be implemented
)

type Models struct {
	Movies MovieModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
