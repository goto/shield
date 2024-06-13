// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	relation "github.com/goto/shield/core/relation"
	mock "github.com/stretchr/testify/mock"
)

// RelationTransformer is an autogenerated mock type for the RelationTransformer type
type RelationTransformer struct {
	mock.Mock
}

type RelationTransformer_Expecter struct {
	mock *mock.Mock
}

func (_m *RelationTransformer) EXPECT() *RelationTransformer_Expecter {
	return &RelationTransformer_Expecter{mock: &_m.Mock}
}

// TransformRelation provides a mock function with given fields: ctx, rlt
func (_m *RelationTransformer) TransformRelation(ctx context.Context, rlt relation.RelationV2) (relation.RelationV2, error) {
	ret := _m.Called(ctx, rlt)

	if len(ret) == 0 {
		panic("no return value specified for TransformRelation")
	}

	var r0 relation.RelationV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) (relation.RelationV2, error)); ok {
		return rf(ctx, rlt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, relation.RelationV2) relation.RelationV2); ok {
		r0 = rf(ctx, rlt)
	} else {
		r0 = ret.Get(0).(relation.RelationV2)
	}

	if rf, ok := ret.Get(1).(func(context.Context, relation.RelationV2) error); ok {
		r1 = rf(ctx, rlt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RelationTransformer_TransformRelation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransformRelation'
type RelationTransformer_TransformRelation_Call struct {
	*mock.Call
}

// TransformRelation is a helper method to define mock.On call
//   - ctx context.Context
//   - rlt relation.RelationV2
func (_e *RelationTransformer_Expecter) TransformRelation(ctx interface{}, rlt interface{}) *RelationTransformer_TransformRelation_Call {
	return &RelationTransformer_TransformRelation_Call{Call: _e.mock.On("TransformRelation", ctx, rlt)}
}

func (_c *RelationTransformer_TransformRelation_Call) Run(run func(ctx context.Context, rlt relation.RelationV2)) *RelationTransformer_TransformRelation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(relation.RelationV2))
	})
	return _c
}

func (_c *RelationTransformer_TransformRelation_Call) Return(_a0 relation.RelationV2, _a1 error) *RelationTransformer_TransformRelation_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RelationTransformer_TransformRelation_Call) RunAndReturn(run func(context.Context, relation.RelationV2) (relation.RelationV2, error)) *RelationTransformer_TransformRelation_Call {
	_c.Call.Return(run)
	return _c
}

// NewRelationTransformer creates a new instance of RelationTransformer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRelationTransformer(t interface {
	mock.TestingT
	Cleanup(func())
}) *RelationTransformer {
	mock := &RelationTransformer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
