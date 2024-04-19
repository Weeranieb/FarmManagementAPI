// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"

	repositories "boonmafarm/api/pkg/repositories"
)

// ISellDetailRepository is an autogenerated mock type for the ISellDetailRepository type
type ISellDetailRepository struct {
	mock.Mock
}

// BulkCreate provides a mock function with given fields: request
func (_m *ISellDetailRepository) BulkCreate(request []models.SellDetail) ([]models.SellDetail, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for BulkCreate")
	}

	var r0 []models.SellDetail
	var r1 error
	if rf, ok := ret.Get(0).(func([]models.SellDetail) ([]models.SellDetail, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func([]models.SellDetail) []models.SellDetail); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.SellDetail)
		}
	}

	if rf, ok := ret.Get(1).(func([]models.SellDetail) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: request
func (_m *ISellDetailRepository) Create(request *models.SellDetail) (*models.SellDetail, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.SellDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.SellDetail) (*models.SellDetail, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.SellDetail) *models.SellDetail); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SellDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.SellDetail) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *ISellDetailRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.SellDetail, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.SellDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.SellDetail, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.SellDetail); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SellDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByQuery provides a mock function with given fields: query, args
func (_m *ISellDetailRepository) ListByQuery(query interface{}, args ...interface{}) ([]models.SellDetail, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListByQuery")
	}

	var r0 []models.SellDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) ([]models.SellDetail, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) []models.SellDetail); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.SellDetail)
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
func (_m *ISellDetailRepository) TakeById(id int) (*models.SellDetail, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.SellDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.SellDetail, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.SellDetail); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SellDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: request
func (_m *ISellDetailRepository) Update(request *models.SellDetail) error {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SellDetail) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTrx provides a mock function with given fields: trxHandle
func (_m *ISellDetailRepository) WithTrx(trxHandle *gorm.DB) repositories.ISellDetailRepository {
	ret := _m.Called(trxHandle)

	if len(ret) == 0 {
		panic("no return value specified for WithTrx")
	}

	var r0 repositories.ISellDetailRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) repositories.ISellDetailRepository); ok {
		r0 = rf(trxHandle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.ISellDetailRepository)
		}
	}

	return r0
}

// NewISellDetailRepository creates a new instance of ISellDetailRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISellDetailRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISellDetailRepository {
	mock := &ISellDetailRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
