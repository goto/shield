// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	metadata "google.golang.org/grpc/metadata"

	mock "github.com/stretchr/testify/mock"
)

// AggregatedDiscoveryService_StreamAggregatedResourcesServer is an autogenerated mock type for the AggregatedDiscoveryService_StreamAggregatedResourcesServer type
type AggregatedDiscoveryService_StreamAggregatedResourcesServer struct {
	mock.Mock
}

type AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter struct {
	mock *mock.Mock
}

func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) EXPECT() *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter{mock: &_m.Mock}
}

// Context provides a mock function with given fields:
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) Context() *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call) Run(run func()) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call) Return(_a0 context.Context) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call) RunAndReturn(run func() context.Context) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Context_Call {
	_c.Call.Return(run)
	return _c
}

// Recv provides a mock function with given fields:
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) Recv() (*discoveryv3.DiscoveryRequest, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Recv")
	}

	var r0 *discoveryv3.DiscoveryRequest
	var r1 error
	if rf, ok := ret.Get(0).(func() (*discoveryv3.DiscoveryRequest, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *discoveryv3.DiscoveryRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv3.DiscoveryRequest)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Recv'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call struct {
	*mock.Call
}

// Recv is a helper method to define mock.On call
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) Recv() *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call{Call: _e.mock.On("Recv")}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call) Run(run func()) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call) Return(_a0 *discoveryv3.DiscoveryRequest, _a1 error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call) RunAndReturn(run func() (*discoveryv3.DiscoveryRequest, error)) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Recv_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for RecvMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) RecvMsg(m interface{}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call) Run(run func(m interface{})) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call) Return(_a0 error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call) RunAndReturn(run func(interface{}) error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) Send(_a0 *discoveryv3.DiscoveryResponse) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*discoveryv3.DiscoveryResponse) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 *discoveryv3.DiscoveryResponse
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) Send(_a0 interface{}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call{Call: _e.mock.On("Send", _a0)}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call) Run(run func(_a0 *discoveryv3.DiscoveryResponse)) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*discoveryv3.DiscoveryResponse))
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call) Return(_a0 error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call) RunAndReturn(run func(*discoveryv3.DiscoveryResponse) error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Send_Call {
	_c.Call.Return(run)
	return _c
}

// SendHeader provides a mock function with given fields: _a0
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHeader'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call struct {
	*mock.Call
}

// SendHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) SendHeader(_a0 interface{}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call{Call: _e.mock.On("SendHeader", _a0)}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call) Run(run func(_a0 metadata.MD)) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call) Return(_a0 error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call) RunAndReturn(run func(metadata.MD) error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for SendMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) SendMsg(m interface{}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call) Run(run func(m interface{})) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call) Return(_a0 error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call) RunAndReturn(run func(interface{}) error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SetHeader provides a mock function with given fields: _a0
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeader'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call struct {
	*mock.Call
}

// SetHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) SetHeader(_a0 interface{}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call{Call: _e.mock.On("SetHeader", _a0)}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call) Run(run func(_a0 metadata.MD)) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call) Return(_a0 error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call) RunAndReturn(run func(metadata.MD) error) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *AggregatedDiscoveryService_StreamAggregatedResourcesServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrailer'
type AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call struct {
	*mock.Call
}

// SetTrailer is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *AggregatedDiscoveryService_StreamAggregatedResourcesServer_Expecter) SetTrailer(_a0 interface{}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call {
	return &AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call{Call: _e.mock.On("SetTrailer", _a0)}
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call) Run(run func(_a0 metadata.MD)) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call) Return() *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call {
	_c.Call.Return()
	return _c
}

func (_c *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call) RunAndReturn(run func(metadata.MD)) *AggregatedDiscoveryService_StreamAggregatedResourcesServer_SetTrailer_Call {
	_c.Call.Return(run)
	return _c
}

// NewAggregatedDiscoveryService_StreamAggregatedResourcesServer creates a new instance of AggregatedDiscoveryService_StreamAggregatedResourcesServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAggregatedDiscoveryService_StreamAggregatedResourcesServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *AggregatedDiscoveryService_StreamAggregatedResourcesServer {
	mock := &AggregatedDiscoveryService_StreamAggregatedResourcesServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
