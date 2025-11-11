package data

import "database/sql"

// Models encloses all the DB Models for easy access using application struct
type Models struct {
	Users  UserModel
	Tokens TokenModel
}

// NewModels returns an Modles struct by
// initilizing it using the provided db connection
func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{db},
		Tokens: TokenModel{db},
	}
}
