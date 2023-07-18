package company

import (
	"context"
	"errors"
	"testing"

	repository "xgo/main/src/internal/controller/company/mock"
	model "xgo/main/src/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestControllerGetCompany(t *testing.T) {
	tests := []struct {
		name    string
		expRes  *model.Company
		expErr  error
		wantRes *model.Company
		wantErr error
	}{
		{
			name:    "not found error",
			expErr:  ErrNotFound,
			wantErr: ErrNotFound,
		},
		{
			name:    "unexpected error",
			expErr:  errors.New("unexpected error"),
			wantErr: errors.New("unexpected error"),
		},
		{
			name:    "success",
			expRes:  &model.Company{},
			wantRes: &model.Company{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := repository.NewMockcompanyRepository(ctrl)
			c := New(repoMock)
			ctx := context.Background()
			id := "id"
			repoMock.EXPECT().GetCompany(ctx, id).Return(tt.expRes, tt.expErr)
			res, err := c.GetCompany(ctx, id)
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}

func TestControllerCreateCompany(t *testing.T) {
	tests := []struct {
		name    string
		expRes  string
		expErr  error
		wantRes string
		wantErr error
	}{
		{
			name:    "unexpected error",
			expErr:  errors.New("unexpected error"),
			wantErr: errors.New("unexpected error"),
		},
		{
			name:    "success",
			expRes:  "1",
			wantRes: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := repository.NewMockcompanyRepository(ctrl)
			c := New(repoMock)
			ctx := context.Background()
			repoMock.EXPECT().CreateCompany(ctx, &model.Company{}).Return(tt.expRes, tt.expErr)
			res, err := c.CreateCompany(ctx, &model.Company{})
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}

func TestControllerPatchCompany(t *testing.T) {
	tests := []struct {
		name    string
		expRes  *model.Company
		wantRes *model.Company
		expErr  error
		wantErr error
	}{
		{
			name:    "not found error",
			expErr:  ErrNotFound,
			wantErr: ErrNotFound,
		},
		{
			name:    "unexpected error",
			expErr:  errors.New("unexpected error"),
			wantErr: errors.New("unexpected error"),
		},
		{
			name:    "success",
			expRes:  &model.Company{},
			wantRes: &model.Company{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := repository.NewMockcompanyRepository(ctrl)
			c := New(repoMock)
			ctx := context.Background()
			repoMock.EXPECT().UpdateCompany(ctx, &model.Company{}).Return(tt.expRes, tt.expErr)
			res, err := c.PatchCompany(ctx, &model.Company{})
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}
