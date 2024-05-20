package database

import (
	"avito-backend-2024-trainee/pkg/postgresql"
	"gorm.io/gorm"
)

func GetTableName(postgresql postgresql.Postgresql, ent any) (string, error) {
	statement := &gorm.Statement{
		DB: postgresql.DB,
	}
	err := statement.Parse(ent)
	if err != nil {
		return "", err
	}
	return statement.Schema.Table, nil
}
