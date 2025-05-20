package teacher

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"student_go/internal/dto/request"
	"student_go/internal/entity"
	"student_go/internal/mocks"
	"student_go/pkg/log"
	"testing"
)

func init() {
	logger, _ := zap.NewDevelopment()
	log.Log = logger
}

func newTestTeacherService() (Service, *mocks.TeacherRepository) {
	mockRepo := new(mocks.TeacherRepository)
	svc := NewTeacherService(mockRepo)
	return svc, mockRepo
}

func TestCreateTeacher(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	input := request.TeacherRequest{Name: "Newton"}
	saved := &entity.Teacher{ID: 1, Name: "Newton"}

	mockRepo.On("Save", mock.AnythingOfType("*entity.Teacher")).Return(saved, nil)

	result, err := svc.CreateTeacher(input)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Newton", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestCreateTeacher_Error(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	mockRepo.On("Save", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.CreateTeacher(request.TeacherRequest{Name: "Fail"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")

	mockRepo.AssertExpectations(t)
}

func TestUpdateTeacher(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	input := request.TeacherRequest{Name: "Updated"}
	updated := &entity.Teacher{
		ID:   2,
		Name: "Updated",
		Courses: []entity.Course{
			{ID: 1, Title: "Math"},
		},
		Departments: []entity.Department{
			{ID: 2, Name: "Science"},
		},
	}

	mockRepo.On("Update", mock.Anything).Return(updated, nil)

	result, err := svc.UpdateTeacher(2, input)

	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Len(t, result.Courses, 1)
	assert.Equal(t, "Math", result.Courses[0].Title)
	assert.Len(t, result.Departments, 1)
	assert.Equal(t, "Science", result.Departments[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTeacher_Error(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	mockRepo.On("Update", mock.Anything).Return(nil, errors.New("update error"))

	result, err := svc.UpdateTeacher(1, request.TeacherRequest{Name: "Fail"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "update error")

	mockRepo.AssertExpectations(t)
}

func TestFindTeacherById(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	teacher := &entity.Teacher{
		ID:   3,
		Name: "Gauss",
		Courses: []entity.Course{
			{ID: 11, Title: "Algebra"},
		},
		Departments: []entity.Department{
			{ID: 22, Name: "Math"},
		},
	}

	mockRepo.On("FindById", uint(3)).Return(teacher, nil)

	result, err := svc.FindTeacherById(3)

	assert.NoError(t, err)
	assert.Equal(t, "Gauss", result.Name)
	assert.Equal(t, "Algebra", result.Courses[0].Title)
	assert.Equal(t, "Math", result.Departments[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestFindTeacherById_Error(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.FindTeacherById(999)

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")

	mockRepo.AssertExpectations(t)
}

func TestFindAllTeachers(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	teachers := []entity.Teacher{
		{
			ID:   1,
			Name: "Tesla",
			Courses: []entity.Course{
				{ID: 5, Title: "Physics"},
			},
			Departments: []entity.Department{
				{ID: 8, Name: "Engineering"},
			},
		},
	}

	mockRepo.On("FindAll", 1, 5).Return(teachers, nil)

	result, err := svc.FindAllTeachers(1, 5)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Tesla", result[0].Name)
	assert.Equal(t, "Physics", result[0].Courses[0].Title)
	assert.Equal(t, "Engineering", result[0].Departments[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestFindAllTeachers_Error(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	mockRepo.On("FindAll", 1, 5).Return(nil, errors.New("fail"))

	result, err := svc.FindAllTeachers(1, 5)

	assert.Nil(t, result)
	assert.EqualError(t, err, "fail")

	mockRepo.AssertExpectations(t)
}

func TestDeleteTeacherById(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	mockRepo.On("DeleteById", uint(1)).Return(nil)

	err := svc.DeleteTeacherById(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTeacherById_Error(t *testing.T) {
	svc, mockRepo := newTestTeacherService()

	mockRepo.On("DeleteById", uint(2)).Return(errors.New("delete error"))

	err := svc.DeleteTeacherById(2)

	assert.EqualError(t, err, "delete error")
	mockRepo.AssertExpectations(t)
}
