// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/Van-programan/Forum_GO/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// CategoryUsecase is an autogenerated mock type for the CategoryUsecase type
type CategoryUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *CategoryUsecase) Create(_a0 context.Context, _a1 entity.Category) (int64, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Category) (int64, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Category) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Category) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *CategoryUsecase) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: _a0
func (_m *CategoryUsecase) GetAll(_a0 context.Context) ([]entity.Category, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []entity.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]entity.Category, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []entity.Category); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *CategoryUsecase) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *entity.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entity.Category, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Category); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, title, description
func (_m *CategoryUsecase) Update(ctx context.Context, id int64, title string, description string) error {
	ret := _m.Called(ctx, id, title, description)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, string) error); ok {
		r0 = rf(ctx, id, title, description)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCategoryUsecase creates a new instance of CategoryUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCategoryUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *CategoryUsecase {
	mock := &CategoryUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}