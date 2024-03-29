// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/datastore/notification.go

// Package mock_datastore is a generated GoMock package.
package mock_datastore

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	datastore "github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
)

// MockNotificationDataStore is a mock of NotificationDataStore interface.
type MockNotificationDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationDataStoreMockRecorder
}

// MockNotificationDataStoreMockRecorder is the mock recorder for MockNotificationDataStore.
type MockNotificationDataStoreMockRecorder struct {
	mock *MockNotificationDataStore
}

// NewMockNotificationDataStore creates a new mock instance.
func NewMockNotificationDataStore(ctrl *gomock.Controller) *MockNotificationDataStore {
	mock := &MockNotificationDataStore{ctrl: ctrl}
	mock.recorder = &MockNotificationDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationDataStore) EXPECT() *MockNotificationDataStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNotificationDataStore) Create(notification *datastore.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", notification)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockNotificationDataStoreMockRecorder) Create(notification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNotificationDataStore)(nil).Create), notification)
}
