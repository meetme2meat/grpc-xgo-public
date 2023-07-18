package company

import (
	"context"
	"errors"
	model "xgo/main/src/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("company data not found")

type companyRepository interface {
	GetCompany(context.Context, string) (*model.Company, error)
	CreateCompany(context.Context, *model.Company) (string, error)
	UpdateCompany(context.Context, *model.Company) (*model.Company, error)
	DeleteCompany(context.Context, string) (string, error)
}

type Controller struct {
	repo companyRepository
}

func New(repo companyRepository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) GetCompany(ctx context.Context, id string) (*model.Company, error) {
	resp, err := c.repo.GetCompany(ctx, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return resp, err
}

func (c *Controller) CreateCompany(ctx context.Context, obj *model.Company) (string, error) {
	zap.L().Debug("inside company controller create function")
	return c.repo.CreateCompany(ctx, obj)
}

func (c *Controller) PatchCompany(ctx context.Context, obj *model.Company) (*model.Company, error) {
	zap.L().Debug("inside company controller patch function")
	resp, err := c.repo.UpdateCompany(ctx, obj)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return resp, err
}

func (c *Controller) DeleteCompany(ctx context.Context, id string) (string, error) {
	zap.L().Debug("inside company controller delete function")
	return c.repo.DeleteCompany(ctx, id)
}
