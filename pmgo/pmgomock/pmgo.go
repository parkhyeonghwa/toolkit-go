// Automatically generated by MockGen. DO NOT EDIT!
// Source: pmgo.go

package pmgomock

import (
	time "time"

	"github.com/percona/toolkit-go/pmgo"
	gomock "github.com/vikstrous/mock/gomock"
	mgo_v2 "gopkg.in/mgo.v2"
)

// Mock of DialerInterface interface
type MockDialerInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockDialerInterfaceRecorder
}

// Recorder for MockDialerInterface (not exported)
type _MockDialerInterfaceRecorder struct {
	mock *MockDialerInterface
}

func NewMockDialerInterface(ctrl *gomock.Controller) *MockDialerInterface {
	mock := &MockDialerInterface{ctrl: ctrl}
	mock.recorder = &_MockDialerInterfaceRecorder{mock}
	return mock
}

func (_m *MockDialerInterface) EXPECT() *_MockDialerInterfaceRecorder {
	return _m.recorder
}

func (_m *MockDialerInterface) Dial(_param0 string) (pmgo.SessionManager, error) {
	ret := _m.ctrl.Call(_m, "Dial", _param0)
	ret0, _ := ret[0].(pmgo.SessionManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDialerInterfaceRecorder) Dial(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Dial", arg0)
}

func (_m *MockDialerInterface) DialWithInfo(_param0 *mgo_v2.DialInfo) (pmgo.SessionManager, error) {
	ret := _m.ctrl.Call(_m, "DialWithInfo", _param0)
	ret0, _ := ret[0].(pmgo.SessionManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDialerInterfaceRecorder) DialWithInfo(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DialWithInfo", arg0)
}

func (_m *MockDialerInterface) DialWithTimeout(_param0 string, _param1 time.Duration) (pmgo.SessionManager, error) {
	ret := _m.ctrl.Call(_m, "DialWithTimeout", _param0, _param1)
	ret0, _ := ret[0].(pmgo.SessionManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDialerInterfaceRecorder) DialWithTimeout(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DialWithTimeout", arg0, arg1)
}
