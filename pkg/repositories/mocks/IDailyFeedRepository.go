// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// IDailyFeedRepository is an autogenerated mock type for the IDailyFeedRepository type
type IDailyFeedRepository struct {
	mock.Mock
}

// BulkCreate provides a mock function with given fields: dailyFeeds
func (_m *IDailyFeedRepository) BulkCreate(dailyFeeds []*models.DailyFeed) ([]*models.DailyFeed, error) {
	ret := _m.Called(dailyFeeds)

	if len(ret) == 0 {
		panic("no return value specified for BulkCreate")
	}

	var r0 []*models.DailyFeed
	var r1 error
	if rf, ok := ret.Get(0).(func([]*models.DailyFeed) ([]*models.DailyFeed, error)); ok {
		return rf(dailyFeeds)
	}
	if rf, ok := ret.Get(0).(func([]*models.DailyFeed) []*models.DailyFeed); ok {
		r0 = rf(dailyFeeds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.DailyFeed)
		}
	}

	if rf, ok := ret.Get(1).(func([]*models.DailyFeed) error); ok {
		r1 = rf(dailyFeeds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: dailyFeed
func (_m *IDailyFeedRepository) Create(dailyFeed *models.DailyFeed) (*models.DailyFeed, error) {
	ret := _m.Called(dailyFeed)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.DailyFeed
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.DailyFeed) (*models.DailyFeed, error)); ok {
		return rf(dailyFeed)
	}
	if rf, ok := ret.Get(0).(func(*models.DailyFeed) *models.DailyFeed); ok {
		r0 = rf(dailyFeed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DailyFeed)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DailyFeed) error); ok {
		r1 = rf(dailyFeed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *IDailyFeedRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.DailyFeed, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.DailyFeed
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.DailyFeed, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.DailyFeed); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DailyFeed)
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
func (_m *IDailyFeedRepository) TakeById(id int) (*models.DailyFeed, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.DailyFeed
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.DailyFeed, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.DailyFeed); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DailyFeed)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: dailyFeed
func (_m *IDailyFeedRepository) Update(dailyFeed *models.DailyFeed) error {
	ret := _m.Called(dailyFeed)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.DailyFeed) error); ok {
		r0 = rf(dailyFeed)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIDailyFeedRepository creates a new instance of IDailyFeedRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDailyFeedRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDailyFeedRepository {
	mock := &IDailyFeedRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
