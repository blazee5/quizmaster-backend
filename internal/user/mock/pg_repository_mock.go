// Code generated by MockGen. DO NOT EDIT.
// Source: internal/user/pg_repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/user/pg_repository.go -destination internal/user/mock/pg_repository_mock.go
//
// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	reflect "reflect"

	domain "github.com/blazee5/quizmaster-backend/internal/domain"
	models "github.com/blazee5/quizmaster-backend/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// ChangeAvatar mocks base method.
func (m *MockRepository) ChangeAvatar(ctx context.Context, userId int, file string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAvatar", ctx, userId, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAvatar indicates an expected call of ChangeAvatar.
func (mr *MockRepositoryMockRecorder) ChangeAvatar(ctx, userId, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAvatar", reflect.TypeOf((*MockRepository)(nil).ChangeAvatar), ctx, userId, file)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAnswer", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAnswer", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, userId)
}

// GetById mocks base method.
func (m *MockRepository) GetById(ctx context.Context, userId int) (models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, userId)
	ret0, _ := ret[0].(models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockRepositoryMockRecorder) GetById(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockRepository)(nil).GetById), ctx, userId)
}

// GetQuizzes mocks base method.
func (m *MockRepository) GetQuizzes(ctx context.Context, userId int) ([]models.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuizzes", ctx, userId)
	ret0, _ := ret[0].([]models.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuizzes indicates an expected call of GetQuizzes.
func (mr *MockRepositoryMockRecorder) GetQuizzes(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuizzes", reflect.TypeOf((*MockRepository)(nil).GetQuizzes), ctx, userId)
}

// GetResults mocks base method.
func (m *MockRepository) GetResults(ctx context.Context, userId int) ([]models.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResults", ctx, userId)
	ret0, _ := ret[0].([]models.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResults indicates an expected call of GetResults.
func (mr *MockRepositoryMockRecorder) GetResults(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResults", reflect.TypeOf((*MockRepository)(nil).GetResults), ctx, userId)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, userId int, input domain.UpdateUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAnswer", ctx, userId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, userId, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAnswer", reflect.TypeOf((*MockRepository)(nil).Update), ctx, userId, input)
}