package department

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"student_go/internal/dto/request"
	"student_go/internal/entity"
	mocks2 "student_go/internal/mocks"
	"student_go/pkg/log"
	"testing"
)

func init() {
	logger, _ := zap.NewDevelopment()
	log.Log = logger
}

func newTestDepartmentService() (Service, *mocks2.DepartmentRepository, *mocks2.TeacherRepository) {
	mockDeptRepo := new(mocks2.DepartmentRepository)
	mockTeacherRepo := new(mocks2.TeacherRepository)
	svc := NewDepartmentService(mockDeptRepo, mockTeacherRepo)
	return svc, mockDeptRepo, mockTeacherRepo
}

func TestCreateDepartment(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	input := request.DepartmentRequest{Name: "Math"}
	saved := &entity.Department{ID: 1, Name: "Math"}

	mockRepo.On("Save", mock.AnythingOfType("*entity.Department")).Return(saved, nil)

	result, err := svc.CreateDepartment(input)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Math", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestCreateDepartment_Error(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	input := request.DepartmentRequest{Name: "Physics"}
	mockRepo.On("Save", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.CreateDepartment(input)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")

	mockRepo.AssertExpectations(t)
}

func TestUpdateDepartment(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	input := request.DepartmentRequest{Name: "Updated"}
	updated := &entity.Department{
		ID:   2,
		Name: "Updated",
		HeadOfDepartment: &entity.Teacher{
			ID:   10,
			Name: "Dr. Admin",
		},
	}

	mockRepo.On("Update", mock.Anything).Return(updated, nil)

	result, err := svc.UpdateDepartment(2, input)

	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, "Dr. Admin", result.HeadOfDepartment.Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDepartment_Error(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockRepo.On("Update", mock.Anything).Return(nil, errors.New("update error"))

	result, err := svc.UpdateDepartment(3, request.DepartmentRequest{Name: "X"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "update error")

	mockRepo.AssertExpectations(t)
}

func TestFindDepartmentById(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockDept := &entity.Department{
		ID:   4,
		Name: "History",
		HeadOfDepartment: &entity.Teacher{
			ID:   20,
			Name: "Dr. Past",
		},
	}

	mockRepo.On("FindById", uint(4)).Return(mockDept, nil)

	result, err := svc.FindDepartmentById(4)

	assert.NoError(t, err)
	assert.Equal(t, "History", result.Name)
	assert.Equal(t, "Dr. Past", result.HeadOfDepartment.Name)

	mockRepo.AssertExpectations(t)
}

func TestFindDepartmentById_Error(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockRepo.On("FindById", uint(404)).Return(nil, errors.New("not found"))

	result, err := svc.FindDepartmentById(404)

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")

	mockRepo.AssertExpectations(t)
}

func TestFindAllDepartments(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockDepts := []entity.Department{
		{
			ID:   1,
			Name: "Chemistry",
			HeadOfDepartment: &entity.Teacher{
				ID:   7,
				Name: "Dr. Lab",
			},
		},
	}

	mockRepo.On("FindAll", 1, 5).Return(mockDepts, nil)

	result, err := svc.FindAllDepartments(1, 5)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Chemistry", result[0].Name)
	assert.Equal(t, "Dr. Lab", result[0].HeadOfDepartment.Name)

	mockRepo.AssertExpectations(t)
}

func TestFindAllDepartments_Error(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockRepo.On("FindAll", 1, 5).Return(nil, errors.New("db error"))

	result, err := svc.FindAllDepartments(1, 5)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")

	mockRepo.AssertExpectations(t)
}

func TestDeleteDepartmentById(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockRepo.On("DeleteById", uint(1)).Return(nil)

	err := svc.DeleteDepartmentById(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteDepartmentById_Error(t *testing.T) {
	svc, mockRepo, _ := newTestDepartmentService()

	mockRepo.On("DeleteById", uint(999)).Return(errors.New("delete error"))

	err := svc.DeleteDepartmentById(999)

	assert.EqualError(t, err, "delete error")
	mockRepo.AssertExpectations(t)
}

func TestDepartmentSetTeacher_DepartmentNotFound(t *testing.T) {
	svc, mockDeptRepo, _ := newTestDepartmentService()

	mockDeptRepo.On("ExistsById", uint(1)).Return(false, nil)

	result, err := svc.DepartmentSetTeacher(1, 10)

	assert.Nil(t, result)
	assert.EqualError(t, err, "department not found")

	mockDeptRepo.AssertExpectations(t)
}

func TestDepartmentSetTeacher_TeacherNotFound(t *testing.T) {
	svc, mockDeptRepo, mockTeacherRepo := newTestDepartmentService()

	mockDeptRepo.On("ExistsById", uint(1)).Return(true, nil)
	mockTeacherRepo.On("ExistsById", uint(10)).Return(false, nil)

	result, err := svc.DepartmentSetTeacher(1, 10)

	assert.Nil(t, result)
	assert.EqualError(t, err, "teacher not found")

	mockDeptRepo.AssertExpectations(t)
	mockTeacherRepo.AssertExpectations(t)
}
