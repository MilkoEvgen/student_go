package course

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

func newTestCourseService() (Service, *mocks2.CourseRepository, *mocks2.TeacherRepository) {
	mockCourseRepo := new(mocks2.CourseRepository)
	mockTeacherRepo := new(mocks2.TeacherRepository)

	svc := NewCourseService(mockCourseRepo, mockTeacherRepo)
	return svc, mockCourseRepo, mockTeacherRepo
}

func TestCreateCourse(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	input := request.CourseRequest{Title: "Math"}
	saved := &entity.Course{ID: 1, Title: "Math"}

	mockCourseRepo.On("Save", mock.AnythingOfType("*entity.Course")).Return(saved, nil)

	result, err := svc.CreateCourse(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Math", result.Title)

	mockCourseRepo.AssertExpectations(t)
}

func TestCreateCourse_Error(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	input := request.CourseRequest{Title: "Physics"}
	mockCourseRepo.On("Save", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.CreateCourse(input)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")

	mockCourseRepo.AssertExpectations(t)
}

func TestFindCourseById(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourse := &entity.Course{
		ID:    1,
		Title: "Math",
		Teacher: &entity.Teacher{
			ID:   2,
			Name: "Dr. Euler",
		},
		Students: []entity.Student{
			{ID: 3, Name: "Bob", Email: "bob@example.com"},
		},
	}

	mockCourseRepo.On("FindById", uint(1)).Return(mockCourse, nil)

	result, err := svc.FindCourseById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Math", result.Title)
	assert.Equal(t, "Dr. Euler", result.Teacher.Name)
	assert.Len(t, result.Students, 1)
	assert.Equal(t, "Bob", result.Students[0].Name)

	mockCourseRepo.AssertExpectations(t)
}

func TestFindCourseById_Error(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourseRepo.On("FindById", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.FindCourseById(999)

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")

	mockCourseRepo.AssertExpectations(t)
}

func TestFindAllCourse(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourses := []entity.Course{
		{
			ID:    1,
			Title: "Biology",
			Teacher: &entity.Teacher{
				ID:   10,
				Name: "Dr. Darwin",
			},
			Students: []entity.Student{
				{ID: 11, Name: "Alice", Email: "alice@example.com"},
			},
		},
	}

	mockCourseRepo.On("FindAll", 1, 5).Return(mockCourses, nil)

	result, err := svc.FindAllCourse(1, 5)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Biology", result[0].Title)
	assert.Equal(t, "Dr. Darwin", result[0].Teacher.Name)

	mockCourseRepo.AssertExpectations(t)
}

func TestFindAllCourse_Error(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourseRepo.On("FindAll", 1, 5).Return(nil, errors.New("db error"))

	result, err := svc.FindAllCourse(1, 5)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")

	mockCourseRepo.AssertExpectations(t)
}

func TestUpdateCourse(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	input := request.CourseRequest{Title: "Updated"}
	mockUpdated := &entity.Course{
		ID:    5,
		Title: "Updated",
		Teacher: &entity.Teacher{
			ID:   1,
			Name: "Dr. Taylor",
		},
		Students: []entity.Student{
			{ID: 3, Name: "Eve", Email: "eve@example.com"},
		},
	}

	mockCourseRepo.On("Update", mock.Anything).Return(mockUpdated, nil)

	result, err := svc.UpdateCourse(5, input)

	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Title)
	assert.Equal(t, "Dr. Taylor", result.Teacher.Name)
	assert.Len(t, result.Students, 1)
	assert.Equal(t, "Eve", result.Students[0].Name)

	mockCourseRepo.AssertExpectations(t)
}

func TestUpdateCourse_Error(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourseRepo.On("Update", mock.Anything).Return(nil, errors.New("update error"))

	result, err := svc.UpdateCourse(1, request.CourseRequest{Title: "X"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "update error")

	mockCourseRepo.AssertExpectations(t)
}

func TestDeleteCourseById(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourseRepo.On("DeleteById", uint(1)).Return(nil)

	err := svc.DeleteCourseById(1)

	assert.NoError(t, err)
	mockCourseRepo.AssertExpectations(t)
}

func TestDeleteCourseById_Error(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourseRepo.On("DeleteById", uint(2)).Return(errors.New("delete error"))

	err := svc.DeleteCourseById(2)

	assert.EqualError(t, err, "delete error")
	mockCourseRepo.AssertExpectations(t)
}

func TestSetTeacherToCourse_CourseNotFound(t *testing.T) {
	svc, mockCourseRepo, _ := newTestCourseService()

	mockCourseRepo.On("ExistsById", uint(1)).Return(false, nil)

	result, err := svc.SetTeacherToCourse(1, 10)

	assert.Nil(t, result)
	assert.EqualError(t, err, "course not found")

	mockCourseRepo.AssertExpectations(t)
}

func TestSetTeacherToCourse_TeacherNotFound(t *testing.T) {
	svc, mockCourseRepo, mockTeacherRepo := newTestCourseService()

	mockCourseRepo.On("ExistsById", uint(1)).Return(true, nil)
	mockTeacherRepo.On("ExistsById", uint(2)).Return(false, nil)

	result, err := svc.SetTeacherToCourse(1, 2)

	assert.Nil(t, result)
	assert.EqualError(t, err, "teacher not found")

	mockCourseRepo.AssertExpectations(t)
	mockTeacherRepo.AssertExpectations(t)
}
