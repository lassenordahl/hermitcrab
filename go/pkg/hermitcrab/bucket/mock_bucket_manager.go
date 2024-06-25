// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/hermitcrab/bucket/interface.go

// Package bucket is a generated GoMock package.
package bucket

import (
	context "context"
	io "io"
	reflect "reflect"

	semver "github.com/Masterminds/semver/v3"
	gomock "github.com/golang/mock/gomock"
)

// MockBucketManager is a mock of BucketManager interface.
type MockBucketManager struct {
	ctrl     *gomock.Controller
	recorder *MockBucketManagerMockRecorder
}

// MockBucketManagerMockRecorder is the mock recorder for MockBucketManager.
type MockBucketManagerMockRecorder struct {
	mock *MockBucketManager
}

// NewMockBucketManager creates a new mock instance.
func NewMockBucketManager(ctrl *gomock.Controller) *MockBucketManager {
	mock := &MockBucketManager{ctrl: ctrl}
	mock.recorder = &MockBucketManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBucketManager) EXPECT() *MockBucketManagerMockRecorder {
	return m.recorder
}

// DownloadPatchVersion mocks base method.
func (m *MockBucketManager) DownloadPatchVersion(ctx context.Context, version *semver.Version) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadPatchVersion", ctx, version)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadPatchVersion indicates an expected call of DownloadPatchVersion.
func (mr *MockBucketManagerMockRecorder) DownloadPatchVersion(ctx, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadPatchVersion", reflect.TypeOf((*MockBucketManager)(nil).DownloadPatchVersion), ctx, version)
}

// GetLatestPatchVersion mocks base method.
func (m *MockBucketManager) GetLatestPatchVersion(ctx context.Context, majorVersion string) (*semver.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestPatchVersion", ctx, majorVersion)
	ret0, _ := ret[0].(*semver.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestPatchVersion indicates an expected call of GetLatestPatchVersion.
func (mr *MockBucketManagerMockRecorder) GetLatestPatchVersion(ctx, majorVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestPatchVersion", reflect.TypeOf((*MockBucketManager)(nil).GetLatestPatchVersion), ctx, majorVersion)
}
