// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	relation "github.com/goto/shield/core/relation"
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
func (_m *Repository) Create(ctx context.Context, _a1 relation.RelationV2) (relation.RelationV2, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 relation.RelationV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) (relation.RelationV2, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) relation.RelationV2); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(relation.RelationV2)
	}

	if rf, ok := ret.Get(1).(func(context.Context, relation.RelationV2) error); ok {
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
//   - _a1 relation.RelationV2
func (_e *Repository_Expecter) Create(ctx interface{}, _a1 interface{}) *Repository_Create_Call {
	return &Repository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *Repository_Create_Call) Run(run func(ctx context.Context, _a1 relation.RelationV2)) *Repository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.RelationV2))
	})
	return _c
}

func (_c *Repository_Create_Call) Return(_a0 relation.RelationV2, _a1 error) *Repository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Create_Call) RunAndReturn(run func(context.Context, relation.RelationV2) (relation.RelationV2, error)) *Repository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteByID provides a mock function with given fields: ctx, id
func (_m *Repository) DeleteByID(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_DeleteByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteByID'
type Repository_DeleteByID_Call struct {
	*mock.Call
}

// DeleteByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *Repository_Expecter) DeleteByID(ctx interface{}, id interface{}) *Repository_DeleteByID_Call {
	return &Repository_DeleteByID_Call{Call: _e.mock.On("DeleteByID", ctx, id)}
}

func (_c *Repository_DeleteByID_Call) Run(run func(ctx context.Context, id string)) *Repository_DeleteByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_DeleteByID_Call) Return(_a0 error) *Repository_DeleteByID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_DeleteByID_Call) RunAndReturn(run func(context.Context, string) error) *Repository_DeleteByID_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *Repository) Get(ctx context.Context, id string) (relation.RelationV2, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 relation.RelationV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (relation.RelationV2, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) relation.RelationV2); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(relation.RelationV2)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type Repository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *Repository_Expecter) Get(ctx interface{}, id interface{}) *Repository_Get_Call {
	return &Repository_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *Repository_Get_Call) Run(run func(ctx context.Context, id string)) *Repository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_Get_Call) Return(_a0 relation.RelationV2, _a1 error) *Repository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Get_Call) RunAndReturn(run func(context.Context, string) (relation.RelationV2, error)) *Repository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByFields provides a mock function with given fields: ctx, rel
func (_m *Repository) GetByFields(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error) {
	ret := _m.Called(ctx, rel)

	if len(ret) == 0 {
		panic("no return value specified for GetByFields")
	}

	var r0 relation.RelationV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) (relation.RelationV2, error)); ok {
		return rf(ctx, rel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) relation.RelationV2); ok {
		r0 = rf(ctx, rel)
	} else {
		r0 = ret.Get(0).(relation.RelationV2)
	}

	if rf, ok := ret.Get(1).(func(context.Context, relation.RelationV2) error); ok {
		r1 = rf(ctx, rel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetByFields_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByFields'
type Repository_GetByFields_Call struct {
	*mock.Call
}

// GetByFields is a helper method to define mock.On call
//   - ctx context.Context
//   - rel relation.RelationV2
func (_e *Repository_Expecter) GetByFields(ctx interface{}, rel interface{}) *Repository_GetByFields_Call {
	return &Repository_GetByFields_Call{Call: _e.mock.On("GetByFields", ctx, rel)}
}

func (_c *Repository_GetByFields_Call) Run(run func(ctx context.Context, rel relation.RelationV2)) *Repository_GetByFields_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.RelationV2))
	})
	return _c
}

func (_c *Repository_GetByFields_Call) Return(_a0 relation.RelationV2, _a1 error) *Repository_GetByFields_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetByFields_Call) RunAndReturn(run func(context.Context, relation.RelationV2) (relation.RelationV2, error)) *Repository_GetByFields_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *Repository) List(ctx context.Context) ([]relation.RelationV2, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []relation.RelationV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]relation.RelationV2, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []relation.RelationV2); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]relation.RelationV2)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
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
func (_e *Repository_Expecter) List(ctx interface{}) *Repository_List_Call {
	return &Repository_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *Repository_List_Call) Run(run func(ctx context.Context)) *Repository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Repository_List_Call) Return(_a0 []relation.RelationV2, _a1 error) *Repository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_List_Call) RunAndReturn(run func(context.Context) ([]relation.RelationV2, error)) *Repository_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, toUpdate
func (_m *Repository) Update(ctx context.Context, toUpdate relation.Relation) (relation.Relation, error) {
	ret := _m.Called(ctx, toUpdate)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 relation.Relation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, relation.Relation) (relation.Relation, error)); ok {
		return rf(ctx, toUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, relation.Relation) relation.Relation); ok {
		r0 = rf(ctx, toUpdate)
	} else {
		r0 = ret.Get(0).(relation.Relation)
	}

	if rf, ok := ret.Get(1).(func(context.Context, relation.Relation) error); ok {
		r1 = rf(ctx, toUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type Repository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate relation.Relation
func (_e *Repository_Expecter) Update(ctx interface{}, toUpdate interface{}) *Repository_Update_Call {
	return &Repository_Update_Call{Call: _e.mock.On("Update", ctx, toUpdate)}
}

func (_c *Repository_Update_Call) Run(run func(ctx context.Context, toUpdate relation.Relation)) *Repository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.Relation))
	})
	return _c
}

func (_c *Repository_Update_Call) Return(_a0 relation.Relation, _a1 error) *Repository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Update_Call) RunAndReturn(run func(context.Context, relation.Relation) (relation.Relation, error)) *Repository_Update_Call {
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
