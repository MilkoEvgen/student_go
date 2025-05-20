package department

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

func setupHandlerTest() (*gin.Engine, *mocks.DepartmentServiceMock, *DepartmentHandler) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.DepartmentServiceMock)
	handler := &DepartmentHandler{Service: mockService}
	r := gin.Default()
	return r, mockService, handler
}

func TestCreateDepartmentHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.DepartmentRequest{Name: "Physics"}
	expected := &response.DepartmentResponse{ID: 1, Name: "Physics"}
	mockService.On("CreateDepartment", input).Return(expected, nil)

	r.POST("/departments", handler.CreateDepartment)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateDepartmentHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	input := request.DepartmentRequest{Name: "Updated"}
	expected := &response.DepartmentResponse{ID: 1, Name: "Updated"}
	mockService.On("UpdateDepartment", uint(1), input).Return(expected, nil)

	r.PATCH("/departments/:id", handler.UpdateDepartment)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPatch, "/departments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindDepartmentByIdHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.DepartmentResponse{ID: 2, Name: "Math"}
	mockService.On("FindDepartmentById", uint(2)).Return(expected, nil)

	r.GET("/departments/:id", handler.FindDepartmentById)
	req := httptest.NewRequest(http.MethodGet, "/departments/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestFindAllDepartmentsHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	depts := []*response.DepartmentResponse{{ID: 1, Name: "X"}}
	mockService.On("Count").Return(1, nil)
	mockService.On("FindAllDepartments", 1, 10).Return(depts, nil)

	r.GET("/departments", handler.FindAllDepartments)
	req := httptest.NewRequest(http.MethodGet, "/departments?page=1&per_page=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteDepartmentHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	mockService.On("DeleteDepartmentById", uint(3)).Return(nil)

	r.DELETE("/departments/:id", handler.DeleteDepartmentById)
	req := httptest.NewRequest(http.MethodDelete, "/departments/3", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDepartmentSetTeacherHandler(t *testing.T) {
	r, mockService, handler := setupHandlerTest()
	expected := &response.DepartmentResponse{ID: 1, Name: "Physics"}
	mockService.On("DepartmentSetTeacher", uint(1), uint(2)).Return(expected, nil)

	r.POST("/departments/:departmentId/teachers/:teacherId", handler.DepartmentSetTeacher)
	req := httptest.NewRequest(http.MethodPost, "/departments/1/teachers/2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
