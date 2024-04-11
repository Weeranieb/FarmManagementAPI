// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"

	repositories "boonmafarm/api/pkg/repositories"
)

// IUserRepository is an autogenerated mock type for the IUserRepository type
type IUserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: user
func (_m *IUserRepository) Create(user *models.User) (*models.User, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.User) (*models.User, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*models.User) *models.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *IUserRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.User, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.User, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.User); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TakeById provides a mock function with given fields: id
func (_m *IUserRepository) TakeById(id int) (*models.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTrx provides a mock function with given fields: trxHandle
func (_m *IUserRepository) WithTrx(trxHandle *gorm.DB) repositories.IUserRepository {
	ret := _m.Called(trxHandle)

	if len(ret) == 0 {
		panic("no return value specified for WithTrx")
	}

	var r0 repositories.IUserRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) repositories.IUserRepository); ok {
		r0 = rf(trxHandle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.IUserRepository)
		}
	}

	return r0
}

// NewIUserRepository creates a new instance of IUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUserRepository {
	mock := &IUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
