package postgresql

import (
	"github.com/gookit/config/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(config.String("main.uri")),
		&gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
