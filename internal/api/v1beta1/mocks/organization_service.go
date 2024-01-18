// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	organization "github.com/goto/shield/core/organization"
	mock "github.com/stretchr/testify/mock"

	user "github.com/goto/shield/core/user"
)

// OrganizationService is an autogenerated mock type for the OrganizationService type
type OrganizationService struct {
	mock.Mock
}

type OrganizationService_Expecter struct {
	mock *mock.Mock
}

func (_m *OrganizationService) EXPECT() *OrganizationService_Expecter {
	return &OrganizationService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, org
func (_m *OrganizationService) Create(ctx context.Context, org organization.Organization) (organization.Organization, error) {
	ret := _m.Called(ctx, org)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 organization.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, organization.Organization) (organization.Organization, error)); ok {
		return rf(ctx, org)
	}
	if rf, ok := ret.Get(0).(func(context.Context, organization.Organization) organization.Organization); ok {
		r0 = rf(ctx, org)
	} else {
		r0 = ret.Get(0).(organization.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, organization.Organization) error); ok {
		r1 = rf(ctx, org)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type OrganizationService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - org organization.Organization
func (_e *OrganizationService_Expecter) Create(ctx interface{}, org interface{}) *OrganizationService_Create_Call {
	return &OrganizationService_Create_Call{Call: _e.mock.On("Create", ctx, org)}
}

func (_c *OrganizationService_Create_Call) Run(run func(ctx context.Context, org organization.Organization)) *OrganizationService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(organization.Organization))
	})
	return _c
}

func (_c *OrganizationService_Create_Call) Return(_a0 organization.Organization, _a1 error) *OrganizationService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_Create_Call) RunAndReturn(run func(context.Context, organization.Organization) (organization.Organization, error)) *OrganizationService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, idOrSlug
func (_m *OrganizationService) Get(ctx context.Context, idOrSlug string) (organization.Organization, error) {
	ret := _m.Called(ctx, idOrSlug)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 organization.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (organization.Organization, error)); ok {
		return rf(ctx, idOrSlug)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) organization.Organization); ok {
		r0 = rf(ctx, idOrSlug)
	} else {
		r0 = ret.Get(0).(organization.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, idOrSlug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type OrganizationService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - idOrSlug string
func (_e *OrganizationService_Expecter) Get(ctx interface{}, idOrSlug interface{}) *OrganizationService_Get_Call {
	return &OrganizationService_Get_Call{Call: _e.mock.On("Get", ctx, idOrSlug)}
}

func (_c *OrganizationService_Get_Call) Run(run func(ctx context.Context, idOrSlug string)) *OrganizationService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *OrganizationService_Get_Call) Return(_a0 organization.Organization, _a1 error) *OrganizationService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_Get_Call) RunAndReturn(run func(context.Context, string) (organization.Organization, error)) *OrganizationService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *OrganizationService) List(ctx context.Context) ([]organization.Organization, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []organization.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]organization.Organization, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []organization.Organization); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]organization.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type OrganizationService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *OrganizationService_Expecter) List(ctx interface{}) *OrganizationService_List_Call {
	return &OrganizationService_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *OrganizationService_List_Call) Run(run func(ctx context.Context)) *OrganizationService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *OrganizationService_List_Call) Return(_a0 []organization.Organization, _a1 error) *OrganizationService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_List_Call) RunAndReturn(run func(context.Context) ([]organization.Organization, error)) *OrganizationService_List_Call {
	_c.Call.Return(run)
	return _c
}

// ListAdmins provides a mock function with given fields: ctx, id
func (_m *OrganizationService) ListAdmins(ctx context.Context, id string) ([]user.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for ListAdmins")
	}

	var r0 []user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]user.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []user.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_ListAdmins_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAdmins'
type OrganizationService_ListAdmins_Call struct {
	*mock.Call
}

// ListAdmins is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *OrganizationService_Expecter) ListAdmins(ctx interface{}, id interface{}) *OrganizationService_ListAdmins_Call {
	return &OrganizationService_ListAdmins_Call{Call: _e.mock.On("ListAdmins", ctx, id)}
}

func (_c *OrganizationService_ListAdmins_Call) Run(run func(ctx context.Context, id string)) *OrganizationService_ListAdmins_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *OrganizationService_ListAdmins_Call) Return(_a0 []user.User, _a1 error) *OrganizationService_ListAdmins_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_ListAdmins_Call) RunAndReturn(run func(context.Context, string) ([]user.User, error)) *OrganizationService_ListAdmins_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, toUpdate
func (_m *OrganizationService) Update(ctx context.Context, toUpdate organization.Organization) (organization.Organization, error) {
	ret := _m.Called(ctx, toUpdate)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 organization.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, organization.Organization) (organization.Organization, error)); ok {
		return rf(ctx, toUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, organization.Organization) organization.Organization); ok {
		r0 = rf(ctx, toUpdate)
	} else {
		r0 = ret.Get(0).(organization.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, organization.Organization) error); ok {
		r1 = rf(ctx, toUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type OrganizationService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate organization.Organization
func (_e *OrganizationService_Expecter) Update(ctx interface{}, toUpdate interface{}) *OrganizationService_Update_Call {
	return &OrganizationService_Update_Call{Call: _e.mock.On("Update", ctx, toUpdate)}
}

func (_c *OrganizationService_Update_Call) Run(run func(ctx context.Context, toUpdate organization.Organization)) *OrganizationService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(organization.Organization))
	})
	return _c
}

func (_c *OrganizationService_Update_Call) Return(_a0 organization.Organization, _a1 error) *OrganizationService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_Update_Call) RunAndReturn(run func(context.Context, organization.Organization) (organization.Organization, error)) *OrganizationService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrganizationService creates a new instance of OrganizationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrganizationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrganizationService {
	mock := &OrganizationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
