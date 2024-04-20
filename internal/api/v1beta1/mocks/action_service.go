// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	action "github.com/goto/shield/core/action"

	mock "github.com/stretchr/testify/mock"
)

// ActionService is an autogenerated mock type for the ActionService type
type ActionService struct {
	mock.Mock
}

type ActionService_Expecter struct {
	mock *mock.Mock
}

func (_m *ActionService) EXPECT() *ActionService_Expecter {
	return &ActionService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1, opts
func (_m *ActionService) Create(ctx context.Context, _a1 action.Action, opts ...action.ServiceOption) (action.Action, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 action.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, action.Action, ...action.ServiceOption) (action.Action, error)); ok {
		return rf(ctx, _a1, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, action.Action, ...action.ServiceOption) action.Action); ok {
		r0 = rf(ctx, _a1, opts...)
	} else {
		r0 = ret.Get(0).(action.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, action.Action, ...action.ServiceOption) error); ok {
		r1 = rf(ctx, _a1, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ActionService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type ActionService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 action.Action
//   - opts ...action.ServiceOption
func (_e *ActionService_Expecter) Create(ctx interface{}, _a1 interface{}, opts ...interface{}) *ActionService_Create_Call {
	return &ActionService_Create_Call{Call: _e.mock.On("Create",
		append([]interface{}{ctx, _a1}, opts...)...)}
}

func (_c *ActionService_Create_Call) Run(run func(ctx context.Context, _a1 action.Action, opts ...action.ServiceOption)) *ActionService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]action.ServiceOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(action.ServiceOption)
			}
		}
		run(args[0].(context.Context), args[1].(action.Action), variadicArgs...)
	})
	return _c
}

func (_c *ActionService_Create_Call) Return(_a0 action.Action, _a1 error) *ActionService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ActionService_Create_Call) RunAndReturn(run func(context.Context, action.Action, ...action.ServiceOption) (action.Action, error)) *ActionService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *ActionService) Get(ctx context.Context, id string) (action.Action, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 action.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (action.Action, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) action.Action); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(action.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ActionService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ActionService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *ActionService_Expecter) Get(ctx interface{}, id interface{}) *ActionService_Get_Call {
	return &ActionService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *ActionService_Get_Call) Run(run func(ctx context.Context, id string)) *ActionService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ActionService_Get_Call) Return(_a0 action.Action, _a1 error) *ActionService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ActionService_Get_Call) RunAndReturn(run func(context.Context, string) (action.Action, error)) *ActionService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *ActionService) List(ctx context.Context) ([]action.Action, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []action.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]action.Action, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []action.Action); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]action.Action)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ActionService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ActionService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ActionService_Expecter) List(ctx interface{}) *ActionService_List_Call {
	return &ActionService_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *ActionService_List_Call) Run(run func(ctx context.Context)) *ActionService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ActionService_List_Call) Return(_a0 []action.Action, _a1 error) *ActionService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ActionService_List_Call) RunAndReturn(run func(context.Context) ([]action.Action, error)) *ActionService_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, id, _a2
func (_m *ActionService) Update(ctx context.Context, id string, _a2 action.Action) (action.Action, error) {
	ret := _m.Called(ctx, id, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 action.Action
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, action.Action) (action.Action, error)); ok {
		return rf(ctx, id, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, action.Action) action.Action); ok {
		r0 = rf(ctx, id, _a2)
	} else {
		r0 = ret.Get(0).(action.Action)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, action.Action) error); ok {
		r1 = rf(ctx, id, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ActionService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ActionService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - _a2 action.Action
func (_e *ActionService_Expecter) Update(ctx interface{}, id interface{}, _a2 interface{}) *ActionService_Update_Call {
	return &ActionService_Update_Call{Call: _e.mock.On("Update", ctx, id, _a2)}
}

func (_c *ActionService_Update_Call) Run(run func(ctx context.Context, id string, _a2 action.Action)) *ActionService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(action.Action))
	})
	return _c
}

func (_c *ActionService_Update_Call) Return(_a0 action.Action, _a1 error) *ActionService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ActionService_Update_Call) RunAndReturn(run func(context.Context, string, action.Action) (action.Action, error)) *ActionService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewActionService creates a new instance of ActionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActionService {
	mock := &ActionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
