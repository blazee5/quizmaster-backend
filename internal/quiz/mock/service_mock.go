// Code generated by MockGen. DO NOT EDIT.
// Source: internal/quiz/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/quiz/service.go -destination internal/quiz/mock/service_mock.go
//
// Package mock_quiz is a generated GoMock package.
package mock_quiz

import (
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	domain "github.com/blazee5/quizmaster-backend/internal/domain"
	models "github.com/blazee5/quizmaster-backend/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockService) Create(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userID, input)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder) Create(ctx, userID, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), ctx, userID, input)
}

// Delete mocks base method.
func (m *MockService) Delete(ctx context.Context, userID, quizID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, userID, quizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(ctx, userID, quizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), ctx, userID, quizID)
}

// DeleteImage mocks base method.
func (m *MockService) DeleteImage(ctx context.Context, userID, quizID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", ctx, userID, quizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockServiceMockRecorder) DeleteImage(ctx, userID, quizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockService)(nil).DeleteImage), ctx, userID, quizID)
}

// GetAll mocks base method.
func (m *MockService) GetAll(ctx context.Context, title, sortBy, sortDir string, page, size int) (models.QuizList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, title, sortBy, sortDir, page, size)
	ret0, _ := ret[0].(models.QuizList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockServiceMockRecorder) GetAll(ctx, title, sortBy, sortDir, page, size any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockService)(nil).GetAll), ctx, title, sortBy, sortDir, page, size)
}

// GetByID mocks base method.
func (m *MockService) GetByID(ctx context.Context, id int) (models.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(models.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServiceMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockService)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *MockService) Update(ctx context.Context, userID, quizID int, input domain.Quiz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, userID, quizID, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockServiceMockRecorder) Update(ctx, userID, quizID, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockService)(nil).Update), ctx, userID, quizID, input)
}

// UploadImage mocks base method.
func (m *MockService) UploadImage(ctx context.Context, userID, quizID int, fileHeader *multipart.FileHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", ctx, userID, quizID, fileHeader)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockServiceMockRecorder) UploadImage(ctx, userID, quizID, fileHeader any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockService)(nil).UploadImage), ctx, userID, quizID, fileHeader)
}
