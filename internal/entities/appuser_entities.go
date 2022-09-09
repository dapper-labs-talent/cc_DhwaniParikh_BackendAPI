package entities

import (
	uuid "github.com/satori/go.uuid"
)

type AppUserEntity struct {
	tableName struct{} `pg:"app_users,alias:t,discard_unknown_columns"`

	ID        uuid.UUID `pg:"id,pk,type:uuid"`
	Email     string    `pg:"email"`
	Password  string    `pg:"password"`
	FirstName string    `pg:"first_name"`
	LastName  string    `pg:"last_name"`
}

type UserDetailsEntity struct {
	tableName struct{} `pg:"app_users,alias:t,discard_unknown_columns"`

	Email     string `pg:"email"`
	FirstName string `pg:"first_name"`
	LastName  string `pg:"last_name"`
}
