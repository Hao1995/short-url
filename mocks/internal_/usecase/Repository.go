// Code generated by mockery v2.52.1. DO NOT EDIT.

package usecase

import (
	context "context"

	domain "github.com/Hao1995/short-url/internal/domain"
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

// Create provides a mock function with given fields: ctx, createDto
func (_m *Repository) Create(ctx context.Context, createDto *domain.CreateDto) (string, error) {
	ret := _m.Called(ctx, createDto)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CreateDto) (string, error)); ok {
		return rf(ctx, createDto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CreateDto) string); ok {
		r0 = rf(ctx, createDto)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.CreateDto) error); ok {
		r1 = rf(ctx, createDto)
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
//   - createDto *domain.CreateDto
func (_e *Repository_Expecter) Create(ctx interface{}, createDto interface{}) *Repository_Create_Call {
	return &Repository_Create_Call{Call: _e.mock.On("Create", ctx, createDto)}
}

func (_c *Repository_Create_Call) Run(run func(ctx context.Context, createDto *domain.CreateDto)) *Repository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.CreateDto))
	})
	return _c
}

func (_c *Repository_Create_Call) Return(_a0 string, _a1 error) *Repository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Create_Call) RunAndReturn(run func(context.Context, *domain.CreateDto) (string, error)) *Repository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *Repository) Get(ctx context.Context, id string) (*domain.ShortUrlDto, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.ShortUrlDto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.ShortUrlDto, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.ShortUrlDto); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ShortUrlDto)
		}
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

func (_c *Repository_Get_Call) Return(_a0 *domain.ShortUrlDto, _a1 error) *Repository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Get_Call) RunAndReturn(run func(context.Context, string) (*domain.ShortUrlDto, error)) *Repository_Get_Call {
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
