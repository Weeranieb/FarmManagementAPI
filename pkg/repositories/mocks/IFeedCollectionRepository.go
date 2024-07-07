// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models "boonmafarm/api/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// IFeedCollectionRepository is an autogenerated mock type for the IFeedCollectionRepository type
type IFeedCollectionRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: feedCollection
func (_m *IFeedCollectionRepository) Create(feedCollection *models.FeedCollection) (*models.FeedCollection, error) {
	ret := _m.Called(feedCollection)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.FeedCollection
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.FeedCollection) (*models.FeedCollection, error)); ok {
		return rf(feedCollection)
	}
	if rf, ok := ret.Get(0).(func(*models.FeedCollection) *models.FeedCollection); ok {
		r0 = rf(feedCollection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FeedCollection)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.FeedCollection) error); ok {
		r1 = rf(feedCollection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstByQuery provides a mock function with given fields: query, args
func (_m *IFeedCollectionRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FeedCollection, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FirstByQuery")
	}

	var r0 *models.FeedCollection
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*models.FeedCollection, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *models.FeedCollection); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FeedCollection)
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
func (_m *IFeedCollectionRepository) TakeById(id int) (*models.FeedCollection, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for TakeById")
	}

	var r0 *models.FeedCollection
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.FeedCollection, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.FeedCollection); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FeedCollection)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TakePage provides a mock function with given fields: clientId, page, pageSize, orderBy, keyword
func (_m *IFeedCollectionRepository) TakePage(clientId int, page int, pageSize int, orderBy string, keyword string) (*[]models.FeedCollection, int64, error) {
	ret := _m.Called(clientId, page, pageSize, orderBy, keyword)

	if len(ret) == 0 {
		panic("no return value specified for TakePage")
	}

	var r0 *[]models.FeedCollection
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(int, int, int, string, string) (*[]models.FeedCollection, int64, error)); ok {
		return rf(clientId, page, pageSize, orderBy, keyword)
	}
	if rf, ok := ret.Get(0).(func(int, int, int, string, string) *[]models.FeedCollection); ok {
		r0 = rf(clientId, page, pageSize, orderBy, keyword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.FeedCollection)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, int, string, string) int64); ok {
		r1 = rf(clientId, page, pageSize, orderBy, keyword)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(int, int, int, string, string) error); ok {
		r2 = rf(clientId, page, pageSize, orderBy, keyword)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Update provides a mock function with given fields: feedCollection
func (_m *IFeedCollectionRepository) Update(feedCollection *models.FeedCollection) error {
	ret := _m.Called(feedCollection)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.FeedCollection) error); ok {
		r0 = rf(feedCollection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIFeedCollectionRepository creates a new instance of IFeedCollectionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIFeedCollectionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IFeedCollectionRepository {
	mock := &IFeedCollectionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}