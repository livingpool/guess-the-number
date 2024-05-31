// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// MockTemplatesInterface is an autogenerated mock type for the TemplatesInterface type
type MockTemplatesInterface struct {
	mock.Mock
}

type MockTemplatesInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTemplatesInterface) EXPECT() *MockTemplatesInterface_Expecter {
	return &MockTemplatesInterface_Expecter{mock: &_m.Mock}
}

// Render provides a mock function with given fields: w, name, data
func (_m *MockTemplatesInterface) Render(w io.Writer, name string, data interface{}) error {
	ret := _m.Called(w, name, data)

	if len(ret) == 0 {
		panic("no return value specified for Render")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, string, interface{}) error); ok {
		r0 = rf(w, name, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTemplatesInterface_Render_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Render'
type MockTemplatesInterface_Render_Call struct {
	*mock.Call
}

// Render is a helper method to define mock.On call
//   - w io.Writer
//   - name string
//   - data interface{}
func (_e *MockTemplatesInterface_Expecter) Render(w interface{}, name interface{}, data interface{}) *MockTemplatesInterface_Render_Call {
	return &MockTemplatesInterface_Render_Call{Call: _e.mock.On("Render", w, name, data)}
}

func (_c *MockTemplatesInterface_Render_Call) Run(run func(w io.Writer, name string, data interface{})) *MockTemplatesInterface_Render_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(io.Writer), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *MockTemplatesInterface_Render_Call) Return(_a0 error) *MockTemplatesInterface_Render_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemplatesInterface_Render_Call) RunAndReturn(run func(io.Writer, string, interface{}) error) *MockTemplatesInterface_Render_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTemplatesInterface creates a new instance of MockTemplatesInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTemplatesInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTemplatesInterface {
	mock := &MockTemplatesInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
