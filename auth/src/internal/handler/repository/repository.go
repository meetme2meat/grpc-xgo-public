package repository

import (
	"context"
	"errors"
	"xgo/auth/src/internal/handler/repository/sqlite"
	"xgo/auth/src/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	driver *gorm.DB
}

func New(name string) *Repository {
	var driver *gorm.DB
	switch name {
	case "sqlite":
		driver = sqlite.New()
	default:
		panic("unknown driver")
	}

	repo := &Repository{driver: driver}
	return repo.runMigration(model.User{})
}

func (r *Repository) runMigration(table interface{}) *Repository {
	if err := r.driver.AutoMigrate(table); err != nil {
		panic(err)
	}

	return r
}

func (r *Repository) Login(ctx context.Context, user model.User, password string) (bool, error) {
	if password == "" {
		return false, errors.New("no password provided")
	}
	trx := r.driver.Table("users").Where("username = ?", user.Username).First(&user)
	if trx.Error != nil {
		return false, trx.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("password mismatch")
	}

	return true, nil
}

func (r *Repository) Create(ctx context.Context, user model.User, password string) error {
	if password == "" {
		return errors.New("no password provided")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.New("unable to generate password hash")
	}

	user.PasswordHash = string(passwordHash)
	trx := r.driver.Table("users").Create(&user)
	return trx.Error
}
