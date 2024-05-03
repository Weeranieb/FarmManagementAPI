// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// IMerchantRepository is an autogenerated mock type for the IMerchantRepository type
type IMerchantRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: merchant
func (_m *IMerchantRepository) Create(merchant *models.Merchant) (*models.Merchant, error) {
	ret := _m.Called(merchant)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.Merchant
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Merchant) (*models.Merchant, error)); ok {
		return rf(merchant)
	}
	if rf, ok := ret.Get(0).(func(*models.Merchant) *models.Merchant); ok {
		r0 = rf(merchant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Merchant)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Merchant) error); ok {
		r1 = rf(merchant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *IMerchantRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.Merchant, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.Merchant
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.Merchant, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.Merchant); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Merchant)
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
func (_m *IMerchantRepository) TakeById(id int) (*models.Merchant, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.Merchant
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.Merchant, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.Merchant); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Merchant)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: merchant
func (_m *IMerchantRepository) Update(merchant *models.Merchant) error {
	ret := _m.Called(merchant)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Merchant) error); ok {
		r0 = rf(merchant)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIMerchantRepository creates a new instance of IMerchantRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIMerchantRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IMerchantRepository {
	mock := &IMerchantRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
