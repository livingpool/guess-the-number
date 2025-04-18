// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
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

// Now provides a mock function with given fields: timeZone
func (_m *MockTimeProviderInterface) Now(timeZone string) time.Time {
	ret := _m.Called(timeZone)

	if len(ret) == 0 {
		panic("no return value specified for Now")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func(string) time.Time); ok {
		r0 = rf(timeZone)
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
//   - timeZone string
func (_e *MockTimeProviderInterface_Expecter) Now(timeZone interface{}) *MockTimeProviderInterface_Now_Call {
	return &MockTimeProviderInterface_Now_Call{Call: _e.mock.On("Now", timeZone)}
}

func (_c *MockTimeProviderInterface_Now_Call) Run(run func(timeZone string)) *MockTimeProviderInterface_Now_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockTimeProviderInterface_Now_Call) Return(_a0 time.Time) *MockTimeProviderInterface_Now_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTimeProviderInterface_Now_Call) RunAndReturn(run func(string) time.Time) *MockTimeProviderInterface_Now_Call {
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
