// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockTimeProviderInterface is an autogenerated mock type for the TimeProviderInterface type
type MockTimeProviderInterface struct {
	mock.Mock
}

type MockTimeProviderInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTimeProviderInterface) EXPECT() *MockTimeProviderInterface_Expecter {
	return &MockTimeProviderInterface_Expecter{mock: &_m.Mock}
}

// Now provides a mock function with given fields:
func (_m *MockTimeProviderInterface) Now() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Now")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// MockTimeProviderInterface_Now_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Now'
type MockTimeProviderInterface_Now_Call struct {
	*mock.Call
}

// Now is a helper method to define mock.On call
func (_e *MockTimeProviderInterface_Expecter) Now() *MockTimeProviderInterface_Now_Call {
	return &MockTimeProviderInterface_Now_Call{Call: _e.mock.On("Now")}
}

func (_c *MockTimeProviderInterface_Now_Call) Run(run func()) *MockTimeProviderInterface_Now_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTimeProviderInterface_Now_Call) Return(_a0 time.Time) *MockTimeProviderInterface_Now_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTimeProviderInterface_Now_Call) RunAndReturn(run func() time.Time) *MockTimeProviderInterface_Now_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTimeProviderInterface creates a new instance of MockTimeProviderInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTimeProviderInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTimeProviderInterface {
	mock := &MockTimeProviderInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
