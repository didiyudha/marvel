// Code generated by MockGen. DO NOT EDIT.
// Source: business/data/character/store.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	character "github.com/didiyudha/marvel/business/data/character"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CleanUp mocks base method.
func (m *MockStore) CleanUp(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanUp", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanUp indicates an expected call of CleanUp.
func (mr *MockStoreMockRecorder) CleanUp(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanUp", reflect.TypeOf((*MockStore)(nil).CleanUp), ctx)
}

// FindAll mocks base method.
func (m *MockStore) FindAll(ctx context.Context) ([]character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockStoreMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockStore)(nil).FindAll), ctx)
}

// FindByID mocks base method.
func (m *MockStore) FindByID(ctx context.Context, id int) (character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockStoreMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockStore)(nil).FindByID), ctx, id)
}

// Save mocks base method.
func (m *MockStore) Save(ctx context.Context, newCharacter ...character.Character) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range newCharacter {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Save", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockStoreMockRecorder) Save(ctx interface{}, newCharacter ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, newCharacter...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStore)(nil).Save), varargs...)
}
