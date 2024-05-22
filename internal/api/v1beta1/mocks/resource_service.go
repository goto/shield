// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	action "github.com/goto/shield/core/action"

	mock "github.com/stretchr/testify/mock"

	resource "github.com/goto/shield/core/resource"
)

// ResourceService is an autogenerated mock type for the ResourceService type
type ResourceService struct {
	mock.Mock
}

type ResourceService_Expecter struct {
	mock *mock.Mock
}

func (_m *ResourceService) EXPECT() *ResourceService_Expecter {
	return &ResourceService_Expecter{mock: &_m.Mock}
}

// CheckAuthz provides a mock function with given fields: ctx, _a1, _a2
func (_m *ResourceService) CheckAuthz(ctx context.Context, _a1 resource.Resource, _a2 action.Action) (bool, error) {
	ret := _m.Called(ctx, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for CheckAuthz")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, resource.Resource, action.Action) (bool, error)); ok {
		return rf(ctx, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, resource.Resource, action.Action) bool); ok {
		r0 = rf(ctx, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, resource.Resource, action.Action) error); ok {
		r1 = rf(ctx, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceService_CheckAuthz_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckAuthz'
type ResourceService_CheckAuthz_Call struct {
	*mock.Call
}

// CheckAuthz is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 resource.Resource
//   - _a2 action.Action
func (_e *ResourceService_Expecter) CheckAuthz(ctx interface{}, _a1 interface{}, _a2 interface{}) *ResourceService_CheckAuthz_Call {
	return &ResourceService_CheckAuthz_Call{Call: _e.mock.On("CheckAuthz", ctx, _a1, _a2)}
}

func (_c *ResourceService_CheckAuthz_Call) Run(run func(ctx context.Context, _a1 resource.Resource, _a2 action.Action)) *ResourceService_CheckAuthz_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(resource.Resource), args[2].(action.Action))
	})
	return _c
}

func (_c *ResourceService_CheckAuthz_Call) Return(_a0 bool, _a1 error) *ResourceService_CheckAuthz_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResourceService_CheckAuthz_Call) RunAndReturn(run func(context.Context, resource.Resource, action.Action) (bool, error)) *ResourceService_CheckAuthz_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *ResourceService) Get(ctx context.Context, id string) (resource.Resource, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 resource.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (resource.Resource, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) resource.Resource); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(resource.Resource)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ResourceService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *ResourceService_Expecter) Get(ctx interface{}, id interface{}) *ResourceService_Get_Call {
	return &ResourceService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *ResourceService_Get_Call) Run(run func(ctx context.Context, id string)) *ResourceService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ResourceService_Get_Call) Return(_a0 resource.Resource, _a1 error) *ResourceService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResourceService_Get_Call) RunAndReturn(run func(context.Context, string) (resource.Resource, error)) *ResourceService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, flt
func (_m *ResourceService) List(ctx context.Context, flt resource.Filter) (resource.PagedResources, error) {
	ret := _m.Called(ctx, flt)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 resource.PagedResources
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, resource.Filter) (resource.PagedResources, error)); ok {
		return rf(ctx, flt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, resource.Filter) resource.PagedResources); ok {
		r0 = rf(ctx, flt)
	} else {
		r0 = ret.Get(0).(resource.PagedResources)
	}

	if rf, ok := ret.Get(1).(func(context.Context, resource.Filter) error); ok {
		r1 = rf(ctx, flt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ResourceService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - flt resource.Filter
func (_e *ResourceService_Expecter) List(ctx interface{}, flt interface{}) *ResourceService_List_Call {
	return &ResourceService_List_Call{Call: _e.mock.On("List", ctx, flt)}
}

func (_c *ResourceService_List_Call) Run(run func(ctx context.Context, flt resource.Filter)) *ResourceService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(resource.Filter))
	})
	return _c
}

func (_c *ResourceService_List_Call) Return(_a0 resource.PagedResources, _a1 error) *ResourceService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResourceService_List_Call) RunAndReturn(run func(context.Context, resource.Filter) (resource.PagedResources, error)) *ResourceService_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, id, _a2
func (_m *ResourceService) Update(ctx context.Context, id string, _a2 resource.Resource) (resource.Resource, error) {
	ret := _m.Called(ctx, id, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 resource.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, resource.Resource) (resource.Resource, error)); ok {
		return rf(ctx, id, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, resource.Resource) resource.Resource); ok {
		r0 = rf(ctx, id, _a2)
	} else {
		r0 = ret.Get(0).(resource.Resource)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, resource.Resource) error); ok {
		r1 = rf(ctx, id, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ResourceService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - _a2 resource.Resource
func (_e *ResourceService_Expecter) Update(ctx interface{}, id interface{}, _a2 interface{}) *ResourceService_Update_Call {
	return &ResourceService_Update_Call{Call: _e.mock.On("Update", ctx, id, _a2)}
}

func (_c *ResourceService_Update_Call) Run(run func(ctx context.Context, id string, _a2 resource.Resource)) *ResourceService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(resource.Resource))
	})
	return _c
}

func (_c *ResourceService_Update_Call) Return(_a0 resource.Resource, _a1 error) *ResourceService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResourceService_Update_Call) RunAndReturn(run func(context.Context, string, resource.Resource) (resource.Resource, error)) *ResourceService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// Upsert provides a mock function with given fields: ctx, _a1
func (_m *ResourceService) Upsert(ctx context.Context, _a1 resource.Resource) (resource.Resource, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Upsert")
	}

	var r0 resource.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, resource.Resource) (resource.Resource, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, resource.Resource) resource.Resource); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(resource.Resource)
	}

	if rf, ok := ret.Get(1).(func(context.Context, resource.Resource) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceService_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type ResourceService_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 resource.Resource
func (_e *ResourceService_Expecter) Upsert(ctx interface{}, _a1 interface{}) *ResourceService_Upsert_Call {
	return &ResourceService_Upsert_Call{Call: _e.mock.On("Upsert", ctx, _a1)}
}

func (_c *ResourceService_Upsert_Call) Run(run func(ctx context.Context, _a1 resource.Resource)) *ResourceService_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(resource.Resource))
	})
	return _c
}

func (_c *ResourceService_Upsert_Call) Return(_a0 resource.Resource, _a1 error) *ResourceService_Upsert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResourceService_Upsert_Call) RunAndReturn(run func(context.Context, resource.Resource) (resource.Resource, error)) *ResourceService_Upsert_Call {
	_c.Call.Return(run)
	return _c
}

// NewResourceService creates a new instance of ResourceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResourceService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ResourceService {
	mock := &ResourceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
