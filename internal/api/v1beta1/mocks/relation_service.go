// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	relation "github.com/goto/shield/core/relation"
	mock "github.com/stretchr/testify/mock"
)

// RelationService is an autogenerated mock type for the RelationService type
type RelationService struct {
	mock.Mock
}

type RelationService_Expecter struct {
	mock *mock.Mock
}

func (_m *RelationService) EXPECT() *RelationService_Expecter {
	return &RelationService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, rel
func (_m *RelationService) Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error) {
	ret := _m.Called(ctx, rel)

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

// RelationService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type RelationService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - rel relation.RelationV2
func (_e *RelationService_Expecter) Create(ctx interface{}, rel interface{}) *RelationService_Create_Call {
	return &RelationService_Create_Call{Call: _e.mock.On("Create", ctx, rel)}
}

func (_c *RelationService_Create_Call) Run(run func(ctx context.Context, rel relation.RelationV2)) *RelationService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.RelationV2))
	})
	return _c
}

func (_c *RelationService_Create_Call) Return(_a0 relation.RelationV2, _a1 error) *RelationService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RelationService_Create_Call) RunAndReturn(run func(context.Context, relation.RelationV2) (relation.RelationV2, error)) *RelationService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteV2 provides a mock function with given fields: ctx, rel
func (_m *RelationService) DeleteV2(ctx context.Context, rel relation.RelationV2) error {
	ret := _m.Called(ctx, rel)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) error); ok {
		r0 = rf(ctx, rel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RelationService_DeleteV2_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteV2'
type RelationService_DeleteV2_Call struct {
	*mock.Call
}

// DeleteV2 is a helper method to define mock.On call
//   - ctx context.Context
//   - rel relation.RelationV2
func (_e *RelationService_Expecter) DeleteV2(ctx interface{}, rel interface{}) *RelationService_DeleteV2_Call {
	return &RelationService_DeleteV2_Call{Call: _e.mock.On("DeleteV2", ctx, rel)}
}

func (_c *RelationService_DeleteV2_Call) Run(run func(ctx context.Context, rel relation.RelationV2)) *RelationService_DeleteV2_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.RelationV2))
	})
	return _c
}

func (_c *RelationService_DeleteV2_Call) Return(_a0 error) *RelationService_DeleteV2_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RelationService_DeleteV2_Call) RunAndReturn(run func(context.Context, relation.RelationV2) error) *RelationService_DeleteV2_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *RelationService) Get(ctx context.Context, id string) (relation.RelationV2, error) {
	ret := _m.Called(ctx, id)

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

// RelationService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type RelationService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *RelationService_Expecter) Get(ctx interface{}, id interface{}) *RelationService_Get_Call {
	return &RelationService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *RelationService_Get_Call) Run(run func(ctx context.Context, id string)) *RelationService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *RelationService_Get_Call) Return(_a0 relation.RelationV2, _a1 error) *RelationService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RelationService_Get_Call) RunAndReturn(run func(context.Context, string) (relation.RelationV2, error)) *RelationService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetRelationByFields provides a mock function with given fields: ctx, rel
func (_m *RelationService) GetRelationByFields(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error) {
	ret := _m.Called(ctx, rel)

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

// RelationService_GetRelationByFields_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRelationByFields'
type RelationService_GetRelationByFields_Call struct {
	*mock.Call
}

// GetRelationByFields is a helper method to define mock.On call
//   - ctx context.Context
//   - rel relation.RelationV2
func (_e *RelationService_Expecter) GetRelationByFields(ctx interface{}, rel interface{}) *RelationService_GetRelationByFields_Call {
	return &RelationService_GetRelationByFields_Call{Call: _e.mock.On("GetRelationByFields", ctx, rel)}
}

func (_c *RelationService_GetRelationByFields_Call) Run(run func(ctx context.Context, rel relation.RelationV2)) *RelationService_GetRelationByFields_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.RelationV2))
	})
	return _c
}

func (_c *RelationService_GetRelationByFields_Call) Return(_a0 relation.RelationV2, _a1 error) *RelationService_GetRelationByFields_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RelationService_GetRelationByFields_Call) RunAndReturn(run func(context.Context, relation.RelationV2) (relation.RelationV2, error)) *RelationService_GetRelationByFields_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *RelationService) List(ctx context.Context) ([]relation.RelationV2, error) {
	ret := _m.Called(ctx)

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

// RelationService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type RelationService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RelationService_Expecter) List(ctx interface{}) *RelationService_List_Call {
	return &RelationService_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *RelationService_List_Call) Run(run func(ctx context.Context)) *RelationService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RelationService_List_Call) Return(_a0 []relation.RelationV2, _a1 error) *RelationService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RelationService_List_Call) RunAndReturn(run func(context.Context) ([]relation.RelationV2, error)) *RelationService_List_Call {
	_c.Call.Return(run)
	return _c
}

// NewRelationService creates a new instance of RelationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRelationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *RelationService {
	mock := &RelationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
