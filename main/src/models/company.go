package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Company struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	EmployeeCount uint32  `json:"employee_count"`
	Registered    bool    `json:"registered"`
	Type          string  `json:"type"`
}

func (c *Company) Table() string {
	return "companies"
}

func (c *Company) AfterCreate(tx *gorm.DB) error {
	if hook, ok := hooks["afterCreate"]; ok {
		zap.L().Debug("after create", zap.Any("company", c))
		hook(c)
	}
	return nil
}

func (c *Company) AfterUpdate(tx *gorm.DB) error {
	if hook, ok := hooks["afterUpdate"]; ok {
		zap.L().Debug("after update", zap.Any("company", c))
		hook(c)
	}
	return nil
}

func (c *Company) AfterSave(tx *gorm.DB) error {
	if hook, ok := hooks["afterSave"]; ok {
		zap.L().Debug("after save", zap.Any("company", c))
		hook(c)
	}
	return nil
}

func (c *Company) AfterDelete(tx *gorm.DB) error {
	if hook, ok := hooks["afterDelete"]; ok {
		zap.L().Debug("after delete", zap.Any("company", c))
		hook(c)
	}
	return nil
}

func (c *Company) RecordId() string   { return c.ID }
func (c *Company) RecordType() string { return c.Table() }
