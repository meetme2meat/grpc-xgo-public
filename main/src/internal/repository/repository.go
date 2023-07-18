package repository

import (
	"context"
	"xgo/main/src/event"
	"xgo/main/src/internal/repository/mysql"
	"xgo/main/src/internal/repository/postgresql"
	"xgo/main/src/internal/repository/sqlite"
	model "xgo/main/src/models"

	"gorm.io/gorm"
)

func New(driver string, eventChan chan event.Event) *Repository {
	repo := Repository{
		db: findDriver(driver),
		ch: eventChan,
	}
	repo.migrate(model.Company{})
	updateHooks(repo)
	return &repo
}

type Repository struct {
	db *gorm.DB
	ch chan event.Event
}

func findDriver(driver string) *gorm.DB {
	switch driver {
	case "mysql":
		return mysql.New()
	case "postgresql":
		return postgresql.New()
	case "sqlite":
		return sqlite.New()
	default:
		panic("not driver found")
	}
}

func (r *Repository) migrate(table interface{}) {
	if err := r.db.AutoMigrate(table); err != nil {
		panic(err)
	}
}

func (r *Repository) GetCompany(ctx context.Context, id string) (*model.Company, error) {
	c := model.Company{ID: id}
	trx := r.db.Table("companies").First(&c)
	return &c, trx.Error
}

func (r *Repository) CreateCompany(ctx context.Context, c *model.Company) (string, error) {
	trx := r.db.Table("companies").Create(c)
	return c.ID, trx.Error
}

func (r *Repository) UpdateCompany(ctx context.Context, c *model.Company) (*model.Company, error) {
	trx := r.db.Table("companies").Save(c)
	return c, trx.Error
}

func (r *Repository) DeleteCompany(ctx context.Context, id string) (string, error) {
	trx := r.db.Table("companies").Delete(&model.Company{ID: id})
	return id, trx.Error
}

func (r *Repository) afterCreate(record event.EventRecord) {
	r.ch <- event.NewEvent(record, "create")
}

func (r *Repository) afterUpdate(record event.EventRecord) {
	r.ch <- event.NewEvent(record, "update")
}

func (r *Repository) afterDelete(record event.EventRecord) {
	r.ch <- event.NewEvent(record, "delete")
}

func updateHooks(repo Repository) {
	model.AddHook("afterCreate", repo.afterCreate)
	model.AddHook("afterUpdate", repo.afterUpdate)
	model.AddHook("afterDelete", repo.afterDelete)
}
