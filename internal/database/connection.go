package database

import (
	"dialog/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Connect(cfg *config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBString()), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "core.", // schema name
		SingularTable: false,
	}})
	if err != nil {
		panic("failed to connect database")
	}

	return db, nil
}
