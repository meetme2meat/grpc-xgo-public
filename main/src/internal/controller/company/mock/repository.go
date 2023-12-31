// Code generated by MockGen. DO NOT EDIT.
// Source: src/internal/controller/company/controller.go

// Package mock_company is a generated GoMock package.
package mock_company

import (
	context "context"
	reflect "reflect"
	model "xgo/main/src/models"

	gomock "go.uber.org/mock/gomock"
)

// MockcompanyRepository is a mock of companyRepository interface.
type MockcompanyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockcompanyRepositoryMockRecorder
}

// MockcompanyRepositoryMockRecorder is the mock recorder for MockcompanyRepository.
type MockcompanyRepositoryMockRecorder struct {
	mock *MockcompanyRepository
}

// NewMockcompanyRepository creates a new mock instance.
func NewMockcompanyRepository(ctrl *gomock.Controller) *MockcompanyRepository {
	mock := &MockcompanyRepository{ctrl: ctrl}
	mock.recorder = &MockcompanyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcompanyRepository) EXPECT() *MockcompanyRepositoryMockRecorder {
	return m.recorder
}

// CreateCompany mocks base method.
func (m *MockcompanyRepository) CreateCompany(arg0 context.Context, arg1 *model.Company) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockcompanyRepositoryMockRecorder) CreateCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockcompanyRepository)(nil).CreateCompany), arg0, arg1)
}

// DeleteCompany mocks base method.
func (m *MockcompanyRepository) DeleteCompany(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCompany", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCompany indicates an expected call of DeleteCompany.
func (mr *MockcompanyRepositoryMockRecorder) DeleteCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCompany", reflect.TypeOf((*MockcompanyRepository)(nil).DeleteCompany), arg0, arg1)
}

// GetCompany mocks base method.
func (m *MockcompanyRepository) GetCompany(arg0 context.Context, arg1 string) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompany", arg0, arg1)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompany indicates an expected call of GetCompany.
func (mr *MockcompanyRepositoryMockRecorder) GetCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompany", reflect.TypeOf((*MockcompanyRepository)(nil).GetCompany), arg0, arg1)
}

// UpdateCompany mocks base method.
func (m *MockcompanyRepository) UpdateCompany(arg0 context.Context, arg1 *model.Company) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCompany", arg0, arg1)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCompany indicates an expected call of UpdateCompany.
func (mr *MockcompanyRepositoryMockRecorder) UpdateCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCompany", reflect.TypeOf((*MockcompanyRepository)(nil).UpdateCompany), arg0, arg1)
}
