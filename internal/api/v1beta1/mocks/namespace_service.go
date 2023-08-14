// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	namespace "github.com/goto/shield/core/namespace"
	mock "github.com/stretchr/testify/mock"
)

// NamespaceService is an autogenerated mock type for the NamespaceService type
type NamespaceService struct {
	mock.Mock
}

type NamespaceService_Expecter struct {
	mock *mock.Mock
}

func (_m *NamespaceService) EXPECT() *NamespaceService_Expecter {
	return &NamespaceService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, ns
func (_m *NamespaceService) Create(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	ret := _m.Called(ctx, ns)

	var r0 namespace.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, namespace.Namespace) (namespace.Namespace, error)); ok {
		return rf(ctx, ns)
	}
	if rf, ok := ret.Get(0).(func(context.Context, namespace.Namespace) namespace.Namespace); ok {
		r0 = rf(ctx, ns)
	} else {
		r0 = ret.Get(0).(namespace.Namespace)
	}

	if rf, ok := ret.Get(1).(func(context.Context, namespace.Namespace) error); ok {
		r1 = rf(ctx, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type NamespaceService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - ns namespace.Namespace
func (_e *NamespaceService_Expecter) Create(ctx interface{}, ns interface{}) *NamespaceService_Create_Call {
	return &NamespaceService_Create_Call{Call: _e.mock.On("Create", ctx, ns)}
}

func (_c *NamespaceService_Create_Call) Run(run func(ctx context.Context, ns namespace.Namespace)) *NamespaceService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(namespace.Namespace))
	})
	return _c
}

func (_c *NamespaceService_Create_Call) Return(_a0 namespace.Namespace, _a1 error) *NamespaceService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NamespaceService_Create_Call) RunAndReturn(run func(context.Context, namespace.Namespace) (namespace.Namespace, error)) *NamespaceService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *NamespaceService) Get(ctx context.Context, id string) (namespace.Namespace, error) {
	ret := _m.Called(ctx, id)

	var r0 namespace.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (namespace.Namespace, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) namespace.Namespace); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(namespace.Namespace)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type NamespaceService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *NamespaceService_Expecter) Get(ctx interface{}, id interface{}) *NamespaceService_Get_Call {
	return &NamespaceService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *NamespaceService_Get_Call) Run(run func(ctx context.Context, id string)) *NamespaceService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *NamespaceService_Get_Call) Return(_a0 namespace.Namespace, _a1 error) *NamespaceService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NamespaceService_Get_Call) RunAndReturn(run func(context.Context, string) (namespace.Namespace, error)) *NamespaceService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *NamespaceService) List(ctx context.Context) ([]namespace.Namespace, error) {
	ret := _m.Called(ctx)

	var r0 []namespace.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]namespace.Namespace, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []namespace.Namespace); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]namespace.Namespace)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type NamespaceService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *NamespaceService_Expecter) List(ctx interface{}) *NamespaceService_List_Call {
	return &NamespaceService_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *NamespaceService_List_Call) Run(run func(ctx context.Context)) *NamespaceService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *NamespaceService_List_Call) Return(_a0 []namespace.Namespace, _a1 error) *NamespaceService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NamespaceService_List_Call) RunAndReturn(run func(context.Context) ([]namespace.Namespace, error)) *NamespaceService_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, ns
func (_m *NamespaceService) Update(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	ret := _m.Called(ctx, ns)

	var r0 namespace.Namespace
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, namespace.Namespace) (namespace.Namespace, error)); ok {
		return rf(ctx, ns)
	}
	if rf, ok := ret.Get(0).(func(context.Context, namespace.Namespace) namespace.Namespace); ok {
		r0 = rf(ctx, ns)
	} else {
		r0 = ret.Get(0).(namespace.Namespace)
	}

	if rf, ok := ret.Get(1).(func(context.Context, namespace.Namespace) error); ok {
		r1 = rf(ctx, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type NamespaceService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - ns namespace.Namespace
func (_e *NamespaceService_Expecter) Update(ctx interface{}, ns interface{}) *NamespaceService_Update_Call {
	return &NamespaceService_Update_Call{Call: _e.mock.On("Update", ctx, ns)}
}

func (_c *NamespaceService_Update_Call) Run(run func(ctx context.Context, ns namespace.Namespace)) *NamespaceService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(namespace.Namespace))
	})
	return _c
}

func (_c *NamespaceService_Update_Call) Return(_a0 namespace.Namespace, _a1 error) *NamespaceService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NamespaceService_Update_Call) RunAndReturn(run func(context.Context, namespace.Namespace) (namespace.Namespace, error)) *NamespaceService_Update_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewNamespaceService interface {
	mock.TestingT
	Cleanup(func())
}

// NewNamespaceService creates a new instance of NamespaceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNamespaceService(t mockConstructorTestingTNewNamespaceService) *NamespaceService {
	mock := &NamespaceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
