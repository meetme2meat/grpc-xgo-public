package sqlite

import (
	"github.com/gookit/config/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	db, err := gorm.Open(
		sqlite.Open(config.String("auth.uri")),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
	return db
}
