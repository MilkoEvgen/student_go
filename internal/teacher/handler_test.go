package teacher_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"student_go/internal/dto/request"
	"student_go/internal/dto/response"
	"student_go/internal/mocks"
	"student_go/internal/teacher"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
)

func setupHandlerTest() (*gin.Engine, *mocks.TeacherServiceMock, *teacher.TeacherHandler) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.TeacherServiceMock)
	handler := &teacher.TeacherHandler{Service: mockService}
	r := gin.Default()
	return r, mockService, handler
}

func TestCreateTeacherHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.TeacherRequest{Name: "Alice"}
	expected := &response.TeacherResponse{ID: 1, Name: "Alice"}
	mockService.On("CreateTeacher", input).Return(expected, nil)

	r.POST("/teachers", handler.CreateTeacher)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/teachers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateTeacherHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.TeacherRequest{Name: "Updated"}
	expected := &response.TeacherResponse{ID: 1, Name: "Updated"}
	mockService.On("UpdateTeacher", uint(1), input).Return(expected, nil)

	r.PATCH("/teachers/:id", handler.UpdateTeacher)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPatch, "/teachers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindTeacherByIdHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.TeacherResponse{ID: 2, Name: "Marie"}
	mockService.On("FindTeacherById", uint(2)).Return(expected, nil)

	r.GET("/teachers/:id", handler.FindTeacherById)
	req := httptest.NewRequest(http.MethodGet, "/teachers/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindAllTeachersHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	teachers := []*response.TeacherResponse{{ID: 1, Name: "X"}}
	mockService.On("Count").Return(1, nil)
	mockService.On("FindAllTeachers", 1, 10).Return(teachers, nil)

	r.GET("/teachers", handler.FindAllTeachers)
	req := httptest.NewRequest(http.MethodGet, "/teachers?page=1&per_page=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteTeacherHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	mockService.On("DeleteTeacherById", uint(3)).Return(nil)

	r.DELETE("/teachers/:id", handler.DeleteTeacherById)
	req := httptest.NewRequest(http.MethodDelete, "/teachers/3", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockService.AssertExpectations(t)
}
