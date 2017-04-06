package mocks

import gin "gopkg.in/gin-gonic/gin.v1"
import handlers "github.com/jakelong95/TownCenter/handlers"
import mock "github.com/stretchr/testify/mock"

// RoasterI is an autogenerated mock type for the RoasterI type
type RoasterI struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx
func (_m *RoasterI) Delete(ctx *gin.Context) {
	_m.Called(ctx)
}

// GetJWT provides a mock function with given fields:
func (_m *RoasterI) GetJWT() gin.HandlerFunc {
	ret := _m.Called()

	var r0 gin.HandlerFunc
	if rf, ok := ret.Get(0).(func() gin.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.HandlerFunc)
		}
	}

	return r0
}

// New provides a mock function with given fields: ctx
func (_m *RoasterI) New(ctx *gin.Context) {
	_m.Called(ctx)
}

// Time provides a mock function with given fields:
func (_m *RoasterI) Time() gin.HandlerFunc {
	ret := _m.Called()

	var r0 gin.HandlerFunc
	if rf, ok := ret.Get(0).(func() gin.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.HandlerFunc)
		}
	}

	return r0
}

// Update provides a mock function with given fields: ctx
func (_m *RoasterI) Update(ctx *gin.Context) {
	_m.Called(ctx)
}

// Upload provides a mock function with given fields: ctx
func (_m *RoasterI) Upload(ctx *gin.Context) {
	_m.Called(ctx)
}

// View provides a mock function with given fields: ctx
func (_m *RoasterI) View(ctx *gin.Context) {
	_m.Called(ctx)
}

// ViewAll provides a mock function with given fields: ctx
func (_m *RoasterI) ViewAll(ctx *gin.Context) {
	_m.Called(ctx)
}

var _ handlers.RoasterI = (*RoasterI)(nil)