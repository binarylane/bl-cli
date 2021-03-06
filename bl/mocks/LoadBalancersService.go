// Code generated by MockGen. DO NOT EDIT.
// Source: load_balancers.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	bl "github.com/binarylane/bl-cli/bl"
	binarylane "github.com/binarylane/go-binarylane"
	gomock "github.com/golang/mock/gomock"
)

// MockLoadBalancersService is a mock of LoadBalancersService interface.
type MockLoadBalancersService struct {
	ctrl     *gomock.Controller
	recorder *MockLoadBalancersServiceMockRecorder
}

// MockLoadBalancersServiceMockRecorder is the mock recorder for MockLoadBalancersService.
type MockLoadBalancersServiceMockRecorder struct {
	mock *MockLoadBalancersService
}

// NewMockLoadBalancersService creates a new mock instance.
func NewMockLoadBalancersService(ctrl *gomock.Controller) *MockLoadBalancersService {
	mock := &MockLoadBalancersService{ctrl: ctrl}
	mock.recorder = &MockLoadBalancersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoadBalancersService) EXPECT() *MockLoadBalancersServiceMockRecorder {
	return m.recorder
}

// AddForwardingRules mocks base method.
func (m *MockLoadBalancersService) AddForwardingRules(lbID int, rules ...binarylane.ForwardingRule) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{lbID}
	for _, a := range rules {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddForwardingRules", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddForwardingRules indicates an expected call of AddForwardingRules.
func (mr *MockLoadBalancersServiceMockRecorder) AddForwardingRules(lbID interface{}, rules ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{lbID}, rules...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddForwardingRules", reflect.TypeOf((*MockLoadBalancersService)(nil).AddForwardingRules), varargs...)
}

// AddServers mocks base method.
func (m *MockLoadBalancersService) AddServers(lbID int, sIDs ...int) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{lbID}
	for _, a := range sIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddServers", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddServers indicates an expected call of AddServers.
func (mr *MockLoadBalancersServiceMockRecorder) AddServers(lbID interface{}, sIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{lbID}, sIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddServers", reflect.TypeOf((*MockLoadBalancersService)(nil).AddServers), varargs...)
}

// Create mocks base method.
func (m *MockLoadBalancersService) Create(lbr *binarylane.LoadBalancerRequest) (*bl.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", lbr)
	ret0, _ := ret[0].(*bl.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockLoadBalancersServiceMockRecorder) Create(lbr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLoadBalancersService)(nil).Create), lbr)
}

// Delete mocks base method.
func (m *MockLoadBalancersService) Delete(lbID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", lbID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockLoadBalancersServiceMockRecorder) Delete(lbID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLoadBalancersService)(nil).Delete), lbID)
}

// Get mocks base method.
func (m *MockLoadBalancersService) Get(lbID int) (*bl.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", lbID)
	ret0, _ := ret[0].(*bl.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockLoadBalancersServiceMockRecorder) Get(lbID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLoadBalancersService)(nil).Get), lbID)
}

// List mocks base method.
func (m *MockLoadBalancersService) List() (bl.LoadBalancers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].(bl.LoadBalancers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockLoadBalancersServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockLoadBalancersService)(nil).List))
}

// RemoveForwardingRules mocks base method.
func (m *MockLoadBalancersService) RemoveForwardingRules(lbID int, rules ...binarylane.ForwardingRule) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{lbID}
	for _, a := range rules {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveForwardingRules", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveForwardingRules indicates an expected call of RemoveForwardingRules.
func (mr *MockLoadBalancersServiceMockRecorder) RemoveForwardingRules(lbID interface{}, rules ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{lbID}, rules...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveForwardingRules", reflect.TypeOf((*MockLoadBalancersService)(nil).RemoveForwardingRules), varargs...)
}

// RemoveServers mocks base method.
func (m *MockLoadBalancersService) RemoveServers(lbID int, sIDs ...int) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{lbID}
	for _, a := range sIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveServers", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveServers indicates an expected call of RemoveServers.
func (mr *MockLoadBalancersServiceMockRecorder) RemoveServers(lbID interface{}, sIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{lbID}, sIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServers", reflect.TypeOf((*MockLoadBalancersService)(nil).RemoveServers), varargs...)
}

// Update mocks base method.
func (m *MockLoadBalancersService) Update(lbID int, lbr *binarylane.LoadBalancerRequest) (*bl.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", lbID, lbr)
	ret0, _ := ret[0].(*bl.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockLoadBalancersServiceMockRecorder) Update(lbID, lbr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLoadBalancersService)(nil).Update), lbID, lbr)
}
