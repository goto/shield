// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	user "github.com/goto/shield/core/user"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *Repository) Create(ctx context.Context, _a1 user.User) (user.User, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.User) (user.User, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.User) user.User); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.User) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Repository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 user.User
func (_e *Repository_Expecter) Create(ctx interface{}, _a1 interface{}) *Repository_Create_Call {
	return &Repository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *Repository_Create_Call) Run(run func(ctx context.Context, _a1 user.User)) *Repository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.User))
	})
	return _c
}

func (_c *Repository_Create_Call) Return(_a0 user.User, _a1 error) *Repository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Create_Call) RunAndReturn(run func(context.Context, user.User) (user.User, error)) *Repository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// CreateMetadataKey provides a mock function with given fields: ctx, key
func (_m *Repository) CreateMetadataKey(ctx context.Context, key user.UserMetadataKey) (user.UserMetadataKey, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for CreateMetadataKey")
	}

	var r0 user.UserMetadataKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.UserMetadataKey) (user.UserMetadataKey, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.UserMetadataKey) user.UserMetadataKey); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(user.UserMetadataKey)
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.UserMetadataKey) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_CreateMetadataKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMetadataKey'
type Repository_CreateMetadataKey_Call struct {
	*mock.Call
}

// CreateMetadataKey is a helper method to define mock.On call
//   - ctx context.Context
//   - key user.UserMetadataKey
func (_e *Repository_Expecter) CreateMetadataKey(ctx interface{}, key interface{}) *Repository_CreateMetadataKey_Call {
	return &Repository_CreateMetadataKey_Call{Call: _e.mock.On("CreateMetadataKey", ctx, key)}
}

func (_c *Repository_CreateMetadataKey_Call) Run(run func(ctx context.Context, key user.UserMetadataKey)) *Repository_CreateMetadataKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.UserMetadataKey))
	})
	return _c
}

func (_c *Repository_CreateMetadataKey_Call) Return(_a0 user.UserMetadataKey, _a1 error) *Repository_CreateMetadataKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_CreateMetadataKey_Call) RunAndReturn(run func(context.Context, user.UserMetadataKey) (user.UserMetadataKey, error)) *Repository_CreateMetadataKey_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *Repository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (user.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) user.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type Repository_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *Repository_Expecter) GetByEmail(ctx interface{}, email interface{}) *Repository_GetByEmail_Call {
	return &Repository_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *Repository_GetByEmail_Call) Run(run func(ctx context.Context, email string)) *Repository_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_GetByEmail_Call) Return(_a0 user.User, _a1 error) *Repository_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetByEmail_Call) RunAndReturn(run func(context.Context, string) (user.User, error)) *Repository_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Repository) GetByID(ctx context.Context, id string) (user.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (user.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) user.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type Repository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *Repository_Expecter) GetByID(ctx interface{}, id interface{}) *Repository_GetByID_Call {
	return &Repository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *Repository_GetByID_Call) Run(run func(ctx context.Context, id string)) *Repository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_GetByID_Call) Return(_a0 user.User, _a1 error) *Repository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetByID_Call) RunAndReturn(run func(context.Context, string) (user.User, error)) *Repository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByIDs provides a mock function with given fields: ctx, userIds
func (_m *Repository) GetByIDs(ctx context.Context, userIds []string) ([]user.User, error) {
	ret := _m.Called(ctx, userIds)

	if len(ret) == 0 {
		panic("no return value specified for GetByIDs")
	}

	var r0 []user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) ([]user.User, error)); ok {
		return rf(ctx, userIds)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) []user.User); ok {
		r0 = rf(ctx, userIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, userIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetByIDs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByIDs'
type Repository_GetByIDs_Call struct {
	*mock.Call
}

// GetByIDs is a helper method to define mock.On call
//   - ctx context.Context
//   - userIds []string
func (_e *Repository_Expecter) GetByIDs(ctx interface{}, userIds interface{}) *Repository_GetByIDs_Call {
	return &Repository_GetByIDs_Call{Call: _e.mock.On("GetByIDs", ctx, userIds)}
}

func (_c *Repository_GetByIDs_Call) Run(run func(ctx context.Context, userIds []string)) *Repository_GetByIDs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *Repository_GetByIDs_Call) Return(_a0 []user.User, _a1 error) *Repository_GetByIDs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetByIDs_Call) RunAndReturn(run func(context.Context, []string) ([]user.User, error)) *Repository_GetByIDs_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, flt
func (_m *Repository) List(ctx context.Context, flt user.Filter) ([]user.User, error) {
	ret := _m.Called(ctx, flt)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.Filter) ([]user.User, error)); ok {
		return rf(ctx, flt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.Filter) []user.User); ok {
		r0 = rf(ctx, flt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.Filter) error); ok {
		r1 = rf(ctx, flt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type Repository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - flt user.Filter
func (_e *Repository_Expecter) List(ctx interface{}, flt interface{}) *Repository_List_Call {
	return &Repository_List_Call{Call: _e.mock.On("List", ctx, flt)}
}

func (_c *Repository_List_Call) Run(run func(ctx context.Context, flt user.Filter)) *Repository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.Filter))
	})
	return _c
}

func (_c *Repository_List_Call) Return(_a0 []user.User, _a1 error) *Repository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_List_Call) RunAndReturn(run func(context.Context, user.Filter) ([]user.User, error)) *Repository_List_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateByEmail provides a mock function with given fields: ctx, toUpdate
func (_m *Repository) UpdateByEmail(ctx context.Context, toUpdate user.User) (user.User, error) {
	ret := _m.Called(ctx, toUpdate)

	if len(ret) == 0 {
		panic("no return value specified for UpdateByEmail")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.User) (user.User, error)); ok {
		return rf(ctx, toUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.User) user.User); ok {
		r0 = rf(ctx, toUpdate)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.User) error); ok {
		r1 = rf(ctx, toUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_UpdateByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateByEmail'
type Repository_UpdateByEmail_Call struct {
	*mock.Call
}

// UpdateByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate user.User
func (_e *Repository_Expecter) UpdateByEmail(ctx interface{}, toUpdate interface{}) *Repository_UpdateByEmail_Call {
	return &Repository_UpdateByEmail_Call{Call: _e.mock.On("UpdateByEmail", ctx, toUpdate)}
}

func (_c *Repository_UpdateByEmail_Call) Run(run func(ctx context.Context, toUpdate user.User)) *Repository_UpdateByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.User))
	})
	return _c
}

func (_c *Repository_UpdateByEmail_Call) Return(_a0 user.User, _a1 error) *Repository_UpdateByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_UpdateByEmail_Call) RunAndReturn(run func(context.Context, user.User) (user.User, error)) *Repository_UpdateByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateByID provides a mock function with given fields: ctx, toUpdate
func (_m *Repository) UpdateByID(ctx context.Context, toUpdate user.User) (user.User, error) {
	ret := _m.Called(ctx, toUpdate)

	if len(ret) == 0 {
		panic("no return value specified for UpdateByID")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.User) (user.User, error)); ok {
		return rf(ctx, toUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.User) user.User); ok {
		r0 = rf(ctx, toUpdate)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.User) error); ok {
		r1 = rf(ctx, toUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_UpdateByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateByID'
type Repository_UpdateByID_Call struct {
	*mock.Call
}

// UpdateByID is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate user.User
func (_e *Repository_Expecter) UpdateByID(ctx interface{}, toUpdate interface{}) *Repository_UpdateByID_Call {
	return &Repository_UpdateByID_Call{Call: _e.mock.On("UpdateByID", ctx, toUpdate)}
}

func (_c *Repository_UpdateByID_Call) Run(run func(ctx context.Context, toUpdate user.User)) *Repository_UpdateByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.User))
	})
	return _c
}

func (_c *Repository_UpdateByID_Call) Return(_a0 user.User, _a1 error) *Repository_UpdateByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_UpdateByID_Call) RunAndReturn(run func(context.Context, user.User) (user.User, error)) *Repository_UpdateByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
