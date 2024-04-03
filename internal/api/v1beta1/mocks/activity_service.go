// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	activity "github.com/goto/shield/core/activity"

	mock "github.com/stretchr/testify/mock"
)

// ActivityService is an autogenerated mock type for the ActivityService type
type ActivityService struct {
	mock.Mock
}

type ActivityService_Expecter struct {
	mock *mock.Mock
}

func (_m *ActivityService) EXPECT() *ActivityService_Expecter {
	return &ActivityService_Expecter{mock: &_m.Mock}
}

// List provides a mock function with given fields: ctx, filter
func (_m *ActivityService) List(ctx context.Context, filter activity.Filter) (activity.PagedActivity, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 activity.PagedActivity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, activity.Filter) (activity.PagedActivity, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, activity.Filter) activity.PagedActivity); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(activity.PagedActivity)
	}

	if rf, ok := ret.Get(1).(func(context.Context, activity.Filter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ActivityService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ActivityService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - filter activity.Filter
func (_e *ActivityService_Expecter) List(ctx interface{}, filter interface{}) *ActivityService_List_Call {
	return &ActivityService_List_Call{Call: _e.mock.On("List", ctx, filter)}
}

func (_c *ActivityService_List_Call) Run(run func(ctx context.Context, filter activity.Filter)) *ActivityService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(activity.Filter))
	})
	return _c
}

func (_c *ActivityService_List_Call) Return(_a0 activity.PagedActivity, _a1 error) *ActivityService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ActivityService_List_Call) RunAndReturn(run func(context.Context, activity.Filter) (activity.PagedActivity, error)) *ActivityService_List_Call {
	_c.Call.Return(run)
	return _c
}

// NewActivityService creates a new instance of ActivityService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActivityService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActivityService {
	mock := &ActivityService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}