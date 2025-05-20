package course

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"student_go/internal/dto/request"
	"student_go/internal/dto/response"
	"student_go/internal/mocks"
	"testing"
)

func setupHandlerTest() (*gin.Engine, *mocks.CourseServiceMock, *Handler) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.CourseServiceMock)
	handler := &Handler{Service: mockService}
	r := gin.Default()
	return r, mockService, handler
}

func TestCreateCourseHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.CourseRequest{Title: "Physics"}
	expected := &response.CourseResponse{ID: 1, Title: "Physics"}
	mockService.On("CreateCourse", input).Return(expected, nil)

	r.POST("/courses", handler.CreateCourse)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateCourseHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.CourseRequest{Title: "Updated"}
	expected := &response.CourseResponse{ID: 1, Title: "Updated"}
	mockService.On("UpdateCourse", uint(1), input).Return(expected, nil)

	r.PATCH("/courses/:id", handler.UpdateCourse)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPatch, "/courses/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindCourseByIdHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.CourseResponse{ID: 2, Title: "Math"}
	mockService.On("FindCourseById", uint(2)).Return(expected, nil)

	r.GET("/courses/:id", handler.FindCourseById)
	req := httptest.NewRequest(http.MethodGet, "/courses/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindAllCoursesHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	courses := []*response.CourseResponse{{ID: 1, Title: "X"}}
	mockService.On("Count").Return(1, nil)
	mockService.On("FindAllCourse", 1, 10).Return(courses, nil)

	r.GET("/courses", handler.FindAllCourses)
	req := httptest.NewRequest(http.MethodGet, "/courses?page=1&per_page=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteCourseHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	mockService.On("DeleteCourseById", uint(3)).Return(nil)

	r.DELETE("/courses/:id", handler.DeleteCourseById)
	req := httptest.NewRequest(http.MethodDelete, "/courses/3", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockService.AssertExpectations(t)
}

func TestSetTeacherToCourseHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.CourseResponse{ID: 1, Title: "Physics"}
	mockService.On("SetTeacherToCourse", uint(1), uint(2)).Return(expected, nil)

	r.POST("/courses/:courseId/teachers/:teacherId", handler.SetTeacherToCourse)
	req := httptest.NewRequest(http.MethodPost, "/courses/1/teachers/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
