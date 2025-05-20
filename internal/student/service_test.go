package student

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

func newTestStudentService() (Service, *mocks2.StudentRepository, *mocks2.CourseRepository) {
	mockStudentRepo := new(mocks2.StudentRepository)
	mockCourseRepo := new(mocks2.CourseRepository)

	svc := NewStudentService(mockStudentRepo, mockCourseRepo)

	return svc, mockStudentRepo, mockCourseRepo
}

func TestCreateStudent(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	input := request.StudentRequest{
		Name:  "Bob",
		Email: "bob@example.com",
	}

	savedStudent := &entity.Student{
		ID:    2,
		Name:  input.Name,
		Email: input.Email,
	}

	mockStudentRepo.On("Save", mock.AnythingOfType("*entity.Student")).Return(savedStudent, nil)

	result, err := studentSvc.CreateStudent(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(2), result.ID)
	assert.Equal(t, "Bob", result.Name)
	assert.Equal(t, "bob@example.com", result.Email)

	mockStudentRepo.AssertExpectations(t)
}

func TestCreateStudent_Error(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	input := request.StudentRequest{
		Name:  "Charlie",
		Email: "charlie@example.com",
	}

	expectedErr := errors.New("db error")

	mockStudentRepo.On("Save", mock.AnythingOfType("*entity.Student")).Return(nil, expectedErr)

	result, err := studentSvc.CreateStudent(input)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")

	mockStudentRepo.AssertExpectations(t)
}

func TestUpdateStudent(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	input := request.StudentRequest{
		Name:  "Eve",
		Email: "eve@example.com",
	}

	updatedStudent := &entity.Student{
		ID:    3,
		Name:  input.Name,
		Email: input.Email,
		Courses: []entity.Course{
			{
				ID:    20,
				Title: "Physics",
				Teacher: &entity.Teacher{
					ID:   200,
					Name: "Dr. Newton",
				},
			},
		},
	}

	mockStudentRepo.On("Update", mock.AnythingOfType("*entity.Student")).Return(updatedStudent, nil)

	result, err := studentSvc.UpdateStudent(3, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(3), result.ID)
	assert.Equal(t, "Eve", result.Name)
	assert.Equal(t, "eve@example.com", result.Email)
	assert.Len(t, result.Courses, 1)
	assert.Equal(t, "Physics", result.Courses[0].Title)
	assert.Equal(t, "Dr. Newton", result.Courses[0].Teacher.Name)

	mockStudentRepo.AssertExpectations(t)
}

func TestUpdateStudent_Error(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	input := request.StudentRequest{
		Name:  "Error",
		Email: "error@example.com",
	}

	expectedErr := errors.New("update failed")

	mockStudentRepo.On("Update", mock.AnythingOfType("*entity.Student")).Return(nil, expectedErr)

	result, err := studentSvc.UpdateStudent(5, input)

	assert.Nil(t, result)
	assert.EqualError(t, err, "update failed")

	mockStudentRepo.AssertExpectations(t)
}

func TestFindStudentById(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	mockStudent := &entity.Student{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Courses: []entity.Course{
			{
				ID:    10,
				Title: "Math",
				Teacher: &entity.Teacher{
					ID:   100,
					Name: "Dr. Smith",
				},
			},
		},
	}

	mockStudentRepo.On("FindById", uint(1)).Return(mockStudent, nil)

	result, err := studentSvc.FindStudentById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Alice", result.Name)
	assert.Equal(t, "alice@example.com", result.Email)
	assert.Len(t, result.Courses, 1)
	assert.Equal(t, "Math", result.Courses[0].Title)
	assert.Equal(t, "Dr. Smith", result.Courses[0].Teacher.Name)

	mockStudentRepo.AssertExpectations(t)
}

func TestFindStudentById_Error(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	expectedErr := errors.New("student not found")

	mockStudentRepo.On("FindById", uint(999)).Return(nil, expectedErr)

	result, err := studentSvc.FindStudentById(999)

	assert.Nil(t, result)
	assert.EqualError(t, err, "student not found")

	mockStudentRepo.AssertExpectations(t)
}

func TestFindAllStudent(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	mockStudents := []entity.Student{
		{
			ID:    1,
			Name:  "Alice",
			Email: "alice@example.com",
			Courses: []entity.Course{
				{
					ID:    10,
					Title: "Biology",
					Teacher: &entity.Teacher{
						ID:   100,
						Name: "Dr. Darwin",
					},
				},
			},
		},
	}

	mockStudentRepo.On("FindAll", 1, 5).Return(mockStudents, nil)

	result, err := studentSvc.FindAllStudent(1, 5)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Alice", result[0].Name)
	assert.Equal(t, "Biology", result[0].Courses[0].Title)
	assert.Equal(t, "Dr. Darwin", result[0].Courses[0].Teacher.Name)

	mockStudentRepo.AssertExpectations(t)
}

func TestFindAllStudent_Error(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	expectedErr := errors.New("find all error")

	mockStudentRepo.On("FindAll", 1, 5).Return(nil, expectedErr)

	result, err := studentSvc.FindAllStudent(1, 5)

	assert.Nil(t, result)
	assert.EqualError(t, err, "find all error")

	mockStudentRepo.AssertExpectations(t)
}

func TestDeleteStudentById(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	mockStudentRepo.On("DeleteById", uint(1)).Return(nil)

	err := studentSvc.DeleteStudentById(1)

	assert.NoError(t, err)
	mockStudentRepo.AssertExpectations(t)
}

func TestDeleteStudentById_Error(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	expectedErr := errors.New("delete failed")
	mockStudentRepo.On("DeleteById", uint(999)).Return(expectedErr)

	err := studentSvc.DeleteStudentById(999)

	assert.EqualError(t, err, "delete failed")
	mockStudentRepo.AssertExpectations(t)
}

func TestAddCourseToStudent_StudentNotFound(t *testing.T) {
	studentSvc, mockStudentRepo, _ := newTestStudentService()

	mockStudentRepo.On("ExistsById", uint(1)).Return(false, nil)

	result, err := studentSvc.AddCourseToStudent(1, 10)

	assert.Nil(t, result)
	assert.EqualError(t, err, "student not found")

	mockStudentRepo.AssertExpectations(t)
}

func TestAddCourseToStudent_CourseNotFound(t *testing.T) {
	studentSvc, mockStudentRepo, mockCourseRepo := newTestStudentService()

	mockStudentRepo.On("ExistsById", uint(1)).Return(true, nil)
	mockCourseRepo.On("ExistsById", uint(10)).Return(false, nil)

	result, err := studentSvc.AddCourseToStudent(1, 10)

	assert.Nil(t, result)
	assert.EqualError(t, err, "course not found")

	mockStudentRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}
