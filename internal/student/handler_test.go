package student

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"student_go/internal/dto/request"
	"student_go/internal/dto/response"
	"student_go/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupHandlerTest() (*gin.Engine, *mocks.StudentServiceMock, *StudentHandler) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.StudentServiceMock)
	handler := &StudentHandler{Service: mockService}
	r := gin.Default()
	return r, mockService, handler
}

func TestCreateStudentHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.StudentRequest{Name: "John", Email: "john@example.com"}
	expected := &response.StudentResponse{ID: 1, Name: "John", Email: "john@example.com"}
	mockService.On("CreateStudent", input).Return(expected, nil)

	r.POST("/students", handler.CreateStudent)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStudentHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.StudentRequest{Name: "Updated", Email: "updated@example.com"}
	expected := &response.StudentResponse{ID: 1, Name: "Updated", Email: "updated@example.com"}
	mockService.On("UpdateStudent", uint(1), input).Return(expected, nil)

	r.PATCH("/students/:id", handler.UpdateStudent)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPatch, "/students/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindStudentByIdHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.StudentResponse{ID: 2, Name: "Alice", Email: "alice@example.com"}
	mockService.On("FindStudentById", uint(2)).Return(expected, nil)

	r.GET("/students/:id", handler.FindStudentById)
	req := httptest.NewRequest(http.MethodGet, "/students/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindAllStudentsHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	students := []*response.StudentResponse{{ID: 1, Name: "X", Email: "x@example.com"}}
	mockService.On("Count").Return(1, nil)
	mockService.On("FindAllStudent", 1, 10).Return(students, nil)

	r.GET("/students", handler.FindAllStudents)
	req := httptest.NewRequest(http.MethodGet, "/students?page=1&per_page=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteStudentHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	mockService.On("DeleteStudentById", uint(3)).Return(nil)

	r.DELETE("/students/:id", handler.DeleteStudentById)
	req := httptest.NewRequest(http.MethodDelete, "/students/3", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindAllCoursesByStudentIdHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	courses := []response.CourseResponse{{ID: 1, Title: "Math"}}
	studentResp := &response.StudentResponse{ID: 1, Name: "John", Courses: courses}
	mockService.On("FindStudentById", uint(1)).Return(studentResp, nil)

	r.GET("/students/:id/courses", handler.FindAllCoursesByStudentId)
	req := httptest.NewRequest(http.MethodGet, "/students/1/courses", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestStudentAddCourseHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.StudentResponse{ID: 1, Name: "John"}
	mockService.On("AddCourseToStudent", uint(1), uint(2)).Return(expected, nil)

	r.POST("/students/:studentId/courses/:courseId", handler.StudentAddCourse)
	req := httptest.NewRequest(http.MethodPost, "/students/1/courses/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
