// Code generated by mockery v2.38.0. DO NOT EDIT.

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

// Log provides a mock function with given fields: ctx, action, actor, data
func (_m *ActivityService) Log(ctx context.Context, action string, actor activity.Actor, data interface{}) error {
	ret := _m.Called(ctx, action, actor, data)

	if len(ret) == 0 {
		panic("no return value specified for Log")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, activity.Actor, interface{}) error); ok {
		r0 = rf(ctx, action, actor, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ActivityService_Log_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Log'
type ActivityService_Log_Call struct {
	*mock.Call
}

// Log is a helper method to define mock.On call
//   - ctx context.Context
//   - action string
//   - actor activity.Actor
//   - data interface{}
func (_e *ActivityService_Expecter) Log(ctx interface{}, action interface{}, actor interface{}, data interface{}) *ActivityService_Log_Call {
	return &ActivityService_Log_Call{Call: _e.mock.On("Log", ctx, action, actor, data)}
}

func (_c *ActivityService_Log_Call) Run(run func(ctx context.Context, action string, actor activity.Actor, data interface{})) *ActivityService_Log_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(activity.Actor), args[3].(interface{}))
	})
	return _c
}

func (_c *ActivityService_Log_Call) Return(_a0 error) *ActivityService_Log_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ActivityService_Log_Call) RunAndReturn(run func(context.Context, string, activity.Actor, interface{}) error) *ActivityService_Log_Call {
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
