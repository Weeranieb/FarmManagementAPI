// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// IActivePondRepository is an autogenerated mock type for the IActivePondRepository type
type IActivePondRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: activePond
func (_m *IActivePondRepository) Create(activePond *models.ActivePond) (*models.ActivePond, error) {
	ret := _m.Called(activePond)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.ActivePond
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.ActivePond) (*models.ActivePond, error)); ok {
		return rf(activePond)
	}
	if rf, ok := ret.Get(0).(func(*models.ActivePond) *models.ActivePond); ok {
		r0 = rf(activePond)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ActivePond)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.ActivePond) error); ok {
		r1 = rf(activePond)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *IActivePondRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.ActivePond, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.ActivePond
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.ActivePond, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.ActivePond); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ActivePond)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetListWithActive provides a mock function with given fields: farmId
func (_m *IActivePondRepository) GetListWithActive(farmId int) ([]*models.PondWithActive, error) {
	ret := _m.Called(farmId)

	if len(ret) == 0 {
		panic("no return value specified for GetListWithActive")
	}

	var r0 []*models.PondWithActive
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*models.PondWithActive, error)); ok {
		return rf(farmId)
	}
	if rf, ok := ret.Get(0).(func(int) []*models.PondWithActive); ok {
		r0 = rf(farmId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.PondWithActive)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(farmId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TakeById provides a mock function with given fields: id
func (_m *IActivePondRepository) TakeById(id int) (*models.ActivePond, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.ActivePond
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.ActivePond, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.ActivePond); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ActivePond)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: activePond
func (_m *IActivePondRepository) Update(activePond *models.ActivePond) error {
	ret := _m.Called(activePond)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.ActivePond) error); ok {
		r0 = rf(activePond)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIActivePondRepository creates a new instance of IActivePondRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIActivePondRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IActivePondRepository {
	mock := &IActivePondRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
