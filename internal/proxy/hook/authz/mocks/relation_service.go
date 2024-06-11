// Code generated by mockery v2.43.2. DO NOT EDIT.

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

// Create provides a mock function with given fields: ctx, _a1
func (_m *RelationService) Create(ctx context.Context, _a1 relation.RelationV2) (relation.RelationV2, error) {
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

// RelationService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type RelationService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 relation.RelationV2
func (_e *RelationService_Expecter) Create(ctx interface{}, _a1 interface{}) *RelationService_Create_Call {
	return &RelationService_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *RelationService_Create_Call) Run(run func(ctx context.Context, _a1 relation.RelationV2)) *RelationService_Create_Call {
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
