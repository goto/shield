// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	user "github.com/goto/shield/core/user"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

type UserService_Expecter struct {
	mock *mock.Mock
}

func (_m *UserService) EXPECT() *UserService_Expecter {
	return &UserService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *UserService) Create(ctx context.Context, _a1 user.User) (user.User, error) {
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

// UserService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 user.User
func (_e *UserService_Expecter) Create(ctx interface{}, _a1 interface{}) *UserService_Create_Call {
	return &UserService_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *UserService_Create_Call) Run(run func(ctx context.Context, _a1 user.User)) *UserService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.User))
	})
	return _c
}

func (_c *UserService_Create_Call) Return(_a0 user.User, _a1 error) *UserService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_Create_Call) RunAndReturn(run func(context.Context, user.User) (user.User, error)) *UserService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// CreateMetadataKey provides a mock function with given fields: ctx, key
func (_m *UserService) CreateMetadataKey(ctx context.Context, key user.UserMetadataKey) (user.UserMetadataKey, error) {
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

// UserService_CreateMetadataKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMetadataKey'
type UserService_CreateMetadataKey_Call struct {
	*mock.Call
}

// CreateMetadataKey is a helper method to define mock.On call
//   - ctx context.Context
//   - key user.UserMetadataKey
func (_e *UserService_Expecter) CreateMetadataKey(ctx interface{}, key interface{}) *UserService_CreateMetadataKey_Call {
	return &UserService_CreateMetadataKey_Call{Call: _e.mock.On("CreateMetadataKey", ctx, key)}
}

func (_c *UserService_CreateMetadataKey_Call) Run(run func(ctx context.Context, key user.UserMetadataKey)) *UserService_CreateMetadataKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.UserMetadataKey))
	})
	return _c
}

func (_c *UserService_CreateMetadataKey_Call) Return(_a0 user.UserMetadataKey, _a1 error) *UserService_CreateMetadataKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_CreateMetadataKey_Call) RunAndReturn(run func(context.Context, user.UserMetadataKey) (user.UserMetadataKey, error)) *UserService_CreateMetadataKey_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *UserService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *UserService_Expecter) Delete(ctx interface{}, id interface{}) *UserService_Delete_Call {
	return &UserService_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *UserService_Delete_Call) Run(run func(ctx context.Context, id string)) *UserService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserService_Delete_Call) Return(_a0 error) *UserService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserService_Delete_Call) RunAndReturn(run func(context.Context, string) error) *UserService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FetchCurrentUser provides a mock function with given fields: ctx
func (_m *UserService) FetchCurrentUser(ctx context.Context) (user.User, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchCurrentUser")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (user.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) user.User); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserService_FetchCurrentUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchCurrentUser'
type UserService_FetchCurrentUser_Call struct {
	*mock.Call
}

// FetchCurrentUser is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserService_Expecter) FetchCurrentUser(ctx interface{}) *UserService_FetchCurrentUser_Call {
	return &UserService_FetchCurrentUser_Call{Call: _e.mock.On("FetchCurrentUser", ctx)}
}

func (_c *UserService_FetchCurrentUser_Call) Run(run func(ctx context.Context)) *UserService_FetchCurrentUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserService_FetchCurrentUser_Call) Return(_a0 user.User, _a1 error) *UserService_FetchCurrentUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_FetchCurrentUser_Call) RunAndReturn(run func(context.Context) (user.User, error)) *UserService_FetchCurrentUser_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, idOrEmail
func (_m *UserService) Get(ctx context.Context, idOrEmail string) (user.User, error) {
	ret := _m.Called(ctx, idOrEmail)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (user.User, error)); ok {
		return rf(ctx, idOrEmail)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) user.User); ok {
		r0 = rf(ctx, idOrEmail)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, idOrEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type UserService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - idOrEmail string
func (_e *UserService_Expecter) Get(ctx interface{}, idOrEmail interface{}) *UserService_Get_Call {
	return &UserService_Get_Call{Call: _e.mock.On("Get", ctx, idOrEmail)}
}

func (_c *UserService_Get_Call) Run(run func(ctx context.Context, idOrEmail string)) *UserService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserService_Get_Call) Return(_a0 user.User, _a1 error) *UserService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_Get_Call) RunAndReturn(run func(context.Context, string) (user.User, error)) *UserService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserService) GetByEmail(ctx context.Context, email string) (user.User, error) {
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

// UserService_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserService_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *UserService_Expecter) GetByEmail(ctx interface{}, email interface{}) *UserService_GetByEmail_Call {
	return &UserService_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *UserService_GetByEmail_Call) Run(run func(ctx context.Context, email string)) *UserService_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserService_GetByEmail_Call) Return(_a0 user.User, _a1 error) *UserService_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_GetByEmail_Call) RunAndReturn(run func(context.Context, string) (user.User, error)) *UserService_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetByIDs provides a mock function with given fields: ctx, userIDs
func (_m *UserService) GetByIDs(ctx context.Context, userIDs []string) ([]user.User, error) {
	ret := _m.Called(ctx, userIDs)

	if len(ret) == 0 {
		panic("no return value specified for GetByIDs")
	}

	var r0 []user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) ([]user.User, error)); ok {
		return rf(ctx, userIDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) []user.User); ok {
		r0 = rf(ctx, userIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, userIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserService_GetByIDs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByIDs'
type UserService_GetByIDs_Call struct {
	*mock.Call
}

// GetByIDs is a helper method to define mock.On call
//   - ctx context.Context
//   - userIDs []string
func (_e *UserService_Expecter) GetByIDs(ctx interface{}, userIDs interface{}) *UserService_GetByIDs_Call {
	return &UserService_GetByIDs_Call{Call: _e.mock.On("GetByIDs", ctx, userIDs)}
}

func (_c *UserService_GetByIDs_Call) Run(run func(ctx context.Context, userIDs []string)) *UserService_GetByIDs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *UserService_GetByIDs_Call) Return(_a0 []user.User, _a1 error) *UserService_GetByIDs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_GetByIDs_Call) RunAndReturn(run func(context.Context, []string) ([]user.User, error)) *UserService_GetByIDs_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, flt
func (_m *UserService) List(ctx context.Context, flt user.Filter) (user.PagedUsers, error) {
	ret := _m.Called(ctx, flt)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 user.PagedUsers
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.Filter) (user.PagedUsers, error)); ok {
		return rf(ctx, flt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.Filter) user.PagedUsers); ok {
		r0 = rf(ctx, flt)
	} else {
		r0 = ret.Get(0).(user.PagedUsers)
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.Filter) error); ok {
		r1 = rf(ctx, flt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type UserService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - flt user.Filter
func (_e *UserService_Expecter) List(ctx interface{}, flt interface{}) *UserService_List_Call {
	return &UserService_List_Call{Call: _e.mock.On("List", ctx, flt)}
}

func (_c *UserService_List_Call) Run(run func(ctx context.Context, flt user.Filter)) *UserService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.Filter))
	})
	return _c
}

func (_c *UserService_List_Call) Return(_a0 user.PagedUsers, _a1 error) *UserService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_List_Call) RunAndReturn(run func(context.Context, user.Filter) (user.PagedUsers, error)) *UserService_List_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateByEmail provides a mock function with given fields: ctx, toUpdate
func (_m *UserService) UpdateByEmail(ctx context.Context, toUpdate user.User) (user.User, error) {
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

// UserService_UpdateByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateByEmail'
type UserService_UpdateByEmail_Call struct {
	*mock.Call
}

// UpdateByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate user.User
func (_e *UserService_Expecter) UpdateByEmail(ctx interface{}, toUpdate interface{}) *UserService_UpdateByEmail_Call {
	return &UserService_UpdateByEmail_Call{Call: _e.mock.On("UpdateByEmail", ctx, toUpdate)}
}

func (_c *UserService_UpdateByEmail_Call) Run(run func(ctx context.Context, toUpdate user.User)) *UserService_UpdateByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.User))
	})
	return _c
}

func (_c *UserService_UpdateByEmail_Call) Return(_a0 user.User, _a1 error) *UserService_UpdateByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_UpdateByEmail_Call) RunAndReturn(run func(context.Context, user.User) (user.User, error)) *UserService_UpdateByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateByID provides a mock function with given fields: ctx, toUpdate
func (_m *UserService) UpdateByID(ctx context.Context, toUpdate user.User) (user.User, error) {
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

// UserService_UpdateByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateByID'
type UserService_UpdateByID_Call struct {
	*mock.Call
}

// UpdateByID is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate user.User
func (_e *UserService_Expecter) UpdateByID(ctx interface{}, toUpdate interface{}) *UserService_UpdateByID_Call {
	return &UserService_UpdateByID_Call{Call: _e.mock.On("UpdateByID", ctx, toUpdate)}
}

func (_c *UserService_UpdateByID_Call) Run(run func(ctx context.Context, toUpdate user.User)) *UserService_UpdateByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.User))
	})
	return _c
}

func (_c *UserService_UpdateByID_Call) Return(_a0 user.User, _a1 error) *UserService_UpdateByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_UpdateByID_Call) RunAndReturn(run func(context.Context, user.User) (user.User, error)) *UserService_UpdateByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
