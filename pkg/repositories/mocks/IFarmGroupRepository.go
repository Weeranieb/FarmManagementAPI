// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// IFarmGroupRepository is an autogenerated mock type for the IFarmGroupRepository type
type IFarmGroupRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: request
func (_m *IFarmGroupRepository) Create(request *models.FarmGroup) (*models.FarmGroup, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.FarmGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.FarmGroup) (*models.FarmGroup, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.FarmGroup) *models.FarmGroup); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FarmGroup)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.FarmGroup) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *IFarmGroupRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.FarmGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.FarmGroup, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.FarmGroup); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FarmGroup)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFarmList provides a mock function with given fields: farmGroupId
func (_m *IFarmGroupRepository) GetFarmList(farmGroupId int) (*[]models.Farm, error) {
	ret := _m.Called(farmGroupId)

	if len(ret) == 0 {
		panic("no return value specified for GetFarmList")
	}

	var r0 *[]models.Farm
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*[]models.Farm, error)); ok {
		return rf(farmGroupId)
	}
	if rf, ok := ret.Get(0).(func(int) *[]models.Farm); ok {
		r0 = rf(farmGroupId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Farm)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(farmGroupId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TakeById provides a mock function with given fields: id
func (_m *IFarmGroupRepository) TakeById(id int) (*models.FarmGroup, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.FarmGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.FarmGroup, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.FarmGroup); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FarmGroup)
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
func (_m *IFarmGroupRepository) Update(request *models.FarmGroup) error {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.FarmGroup) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIFarmGroupRepository creates a new instance of IFarmGroupRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIFarmGroupRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IFarmGroupRepository {
	mock := &IFarmGroupRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}