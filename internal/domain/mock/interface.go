// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	domain "hash-api/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHashKeeper is a mock of HashKeeper interface.
type MockHashKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockHashKeeperMockRecorder
}

// MockHashKeeperMockRecorder is the mock recorder for MockHashKeeper.
type MockHashKeeperMockRecorder struct {
	mock *MockHashKeeper
}

// NewMockHashKeeper creates a new mock instance.
func NewMockHashKeeper(ctrl *gomock.Controller) *MockHashKeeper {
	mock := &MockHashKeeper{ctrl: ctrl}
	mock.recorder = &MockHashKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHashKeeper) EXPECT() *MockHashKeeperMockRecorder {
	return m.recorder
}

// Store mocks base method.
func (m *MockHashKeeper) Store(arg0 context.Context, arg1 domain.Hash) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockHashKeeperMockRecorder) Store(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockHashKeeper)(nil).Store), arg0, arg1)
}

// MockHashGetter is a mock of HashGetter interface.
type MockHashGetter struct {
	ctrl     *gomock.Controller
	recorder *MockHashGetterMockRecorder
}

// MockHashGetterMockRecorder is the mock recorder for MockHashGetter.
type MockHashGetterMockRecorder struct {
	mock *MockHashGetter
}

// NewMockHashGetter creates a new mock instance.
func NewMockHashGetter(ctrl *gomock.Controller) *MockHashGetter {
	mock := &MockHashGetter{ctrl: ctrl}
	mock.recorder = &MockHashGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHashGetter) EXPECT() *MockHashGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockHashGetter) Get(arg0 context.Context) (domain.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(domain.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockHashGetterMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHashGetter)(nil).Get), arg0)
}
