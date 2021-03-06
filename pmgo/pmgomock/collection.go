// Automatically generated by MockGen. DO NOT EDIT!
// Source: collection.go

package pmgomock

import (
	. "github.com/percona/toolkit-go/pmgo"
	gomock "github.com/vikstrous/mock/gomock"
	mgo_v2 "gopkg.in/mgo.v2"
)

// Mock of CollectionManager interface
type MockCollectionManager struct {
	ctrl     *gomock.Controller
	recorder *_MockCollectionManagerRecorder
}

// Recorder for MockCollectionManager (not exported)
type _MockCollectionManagerRecorder struct {
	mock *MockCollectionManager
}

func NewMockCollectionManager(ctrl *gomock.Controller) *MockCollectionManager {
	mock := &MockCollectionManager{ctrl: ctrl}
	mock.recorder = &_MockCollectionManagerRecorder{mock}
	return mock
}

func (_m *MockCollectionManager) EXPECT() *_MockCollectionManagerRecorder {
	return _m.recorder
}

func (_m *MockCollectionManager) Count() (int, error) {
	ret := _m.ctrl.Call(_m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCollectionManagerRecorder) Count() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Count")
}

func (_m *MockCollectionManager) Create(_param0 *mgo_v2.CollectionInfo) error {
	ret := _m.ctrl.Call(_m, "Create", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCollectionManagerRecorder) Create(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0)
}

func (_m *MockCollectionManager) Find(_param0 interface{}) QueryManager {
	ret := _m.ctrl.Call(_m, "Find", _param0)
	ret0, _ := ret[0].(QueryManager)
	return ret0
}

func (_mr *_MockCollectionManagerRecorder) Find(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Find", arg0)
}

func (_m *MockCollectionManager) Pipe(_param0 interface{}) *mgo_v2.Pipe {
	ret := _m.ctrl.Call(_m, "Pipe", _param0)
	ret0, _ := ret[0].(*mgo_v2.Pipe)
	return ret0
}

func (_mr *_MockCollectionManagerRecorder) Pipe(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Pipe", arg0)
}
