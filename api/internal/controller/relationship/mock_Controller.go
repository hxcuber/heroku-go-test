// Code generated by mockery v2.33.0. DO NOT EDIT.

package relationship

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockController is an autogenerated mock type for the Controller type
type MockController struct {
	mock.Mock
}

// Befriend provides a mock function with given fields: ctx, email1, email2
func (_m *MockController) Befriend(ctx context.Context, email1 string, email2 string) error {
	ret := _m.Called(ctx, email1, email2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email1, email2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Block provides a mock function with given fields: ctx, requestorEmail, targetEmail
func (_m *MockController) Block(ctx context.Context, requestorEmail string, targetEmail string) error {
	ret := _m.Called(ctx, requestorEmail, targetEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, requestorEmail, targetEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CommonFriends provides a mock function with given fields: ctx, email1, email2
func (_m *MockController) CommonFriends(ctx context.Context, email1 string, email2 string) ([]string, error) {
	ret := _m.Called(ctx, email1, email2)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]string, error)); ok {
		return rf(ctx, email1, email2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []string); ok {
		r0 = rf(ctx, email1, email2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email1, email2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Friends provides a mock function with given fields: ctx, email
func (_m *MockController) Friends(ctx context.Context, email string) ([]string, error) {
	ret := _m.Called(ctx, email)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Receivers provides a mock function with given fields: ctx, email, text
func (_m *MockController) Receivers(ctx context.Context, email string, text string) ([]string, error) {
	ret := _m.Called(ctx, email, text)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]string, error)); ok {
		return rf(ctx, email, text)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []string); ok {
		r0 = rf(ctx, email, text)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subscribe provides a mock function with given fields: ctx, requestorEmail, targetEmail
func (_m *MockController) Subscribe(ctx context.Context, requestorEmail string, targetEmail string) error {
	ret := _m.Called(ctx, requestorEmail, targetEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, requestorEmail, targetEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockController creates a new instance of MockController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockController {
	mock := &MockController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
