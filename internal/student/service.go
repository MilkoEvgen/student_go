package student

import (
	"fmt"
	"go.uber.org/zap"
	"student_go/internal/course"
	"student_go/internal/dto/request"
	response3 "student_go/internal/dto/response"
	"student_go/internal/entity"
	"student_go/pkg/dbcontext"
	"student_go/pkg/log"
)

type Service interface {
	CreateStudent(input request.StudentRequest) (*response3.StudentResponse, error)
	UpdateStudent(id uint, input request.StudentRequest) (*response3.StudentResponse, error)
	FindStudentById(id uint) (*response3.StudentResponse, error)
	FindAllStudent(page, limit int) ([]*response3.StudentResponse, error)
	DeleteStudentById(id uint) error
	AddCourseToStudent(studentId uint, courseId uint) (*response3.StudentResponse, error)
	Count() (int, error)
}

type service struct {
	studentRepository Repository
	courseRepository  course.Repository
}

func NewStudentService(
	studentRepository Repository,
	courseRepository course.Repository) Service {
	return &service{
		studentRepository: studentRepository,
		courseRepository:  courseRepository,
	}
}

func (s *service) CreateStudent(input request.StudentRequest) (*response3.StudentResponse, error) {
	log.Log.Info("CreateStudent (service) called", zap.String("name", input.Name), zap.String("email", input.Email))

	student := entity.Student{
		Name:  input.Name,
		Email: input.Email,
	}
	savedStudent, err := s.studentRepository.Save(&student)
	if err != nil {
		return nil, err
	}

	resp := &response3.StudentResponse{
		ID:    savedStudent.ID,
		Name:  savedStudent.Name,
		Email: savedStudent.Email,
	}
	return resp, nil
}

func (s *service) UpdateStudent(id uint, input request.StudentRequest) (*response3.StudentResponse, error) {
	log.Log.Info("UpdateStudent (service) called", zap.Uint("id", id), zap.String("name", input.Name), zap.String("email", input.Email))

	student := entity.Student{
		ID:    id,
		Name:  input.Name,
		Email: input.Email,
	}
	updatedStudent, err := s.studentRepository.Update(&student)
	if err != nil {
		return nil, err
	}

	coursesResp := make([]response3.CourseResponse, 0, len(student.Courses))
	for _, course := range updatedStudent.Courses {
		var teacherResp *response3.TeacherResponse
		if course.Teacher != nil {
			teacherResp = &response3.TeacherResponse{
				ID:   course.Teacher.ID,
				Name: course.Teacher.Name,
			}
		}

		courseResp := response3.CourseResponse{
			ID:      course.ID,
			Title:   course.Title,
			Teacher: teacherResp,
		}
		coursesResp = append(coursesResp, courseResp)
	}

	studentResp := &response3.StudentResponse{
		ID:      student.ID,
		Name:    student.Name,
		Email:   student.Email,
		Courses: coursesResp,
	}
	return studentResp, nil
}

func (s *service) FindStudentById(id uint) (*response3.StudentResponse, error) {
	log.Log.Info("FindStudentById (service) called", zap.Uint("id", id))

	student, err := s.studentRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	coursesResp := make([]response3.CourseResponse, 0, len(student.Courses))
	for _, course := range student.Courses {
		var teacherResp *response3.TeacherResponse
		if course.Teacher != nil {
			teacherResp = &response3.TeacherResponse{
				ID:   course.Teacher.ID,
				Name: course.Teacher.Name,
			}
		}

		courseResp := response3.CourseResponse{
			ID:      course.ID,
			Title:   course.Title,
			Teacher: teacherResp,
		}
		coursesResp = append(coursesResp, courseResp)
	}

	studentResp := &response3.StudentResponse{
		ID:      student.ID,
		Name:    student.Name,
		Email:   student.Email,
		Courses: coursesResp,
	}
	return studentResp, nil
}

func (s *service) FindAllStudent(page, limit int) ([]*response3.StudentResponse, error) {
	log.Log.Info("FindAllStudent (service) called", zap.Int("page", page), zap.Int("limit", limit))

	students, err := s.studentRepository.FindAll(page, limit)
	if err != nil {
		return nil, err
	}

	var studentResponses []*response3.StudentResponse
	for _, student := range students {
		coursesResp := make([]response3.CourseResponse, 0, len(student.Courses))
		for _, course := range student.Courses {
			var teacherResp *response3.TeacherResponse
			if course.Teacher != nil {
				teacherResp = &response3.TeacherResponse{
					ID:   course.Teacher.ID,
					Name: course.Teacher.Name,
				}
			}

			courseResp := response3.CourseResponse{
				ID:      course.ID,
				Title:   course.Title,
				Teacher: teacherResp,
			}
			coursesResp = append(coursesResp, courseResp)
		}

		studentResp := &response3.StudentResponse{
			ID:      student.ID,
			Name:    student.Name,
			Email:   student.Email,
			Courses: coursesResp,
		}
		studentResponses = append(studentResponses, studentResp)
	}

	return studentResponses, nil
}

func (s *service) DeleteStudentById(id uint) error {
	log.Log.Info("DeleteStudentById (service) called", zap.Uint("id", id))
	return s.studentRepository.DeleteById(id)
}

func (s *service) AddCourseToStudent(studentId uint, courseId uint) (*response3.StudentResponse, error) {
	log.Log.Info("AddCourseToStudent (service) called", zap.Uint("student_id", studentId), zap.Uint("course_id", courseId))

	exists, err := s.studentRepository.ExistsById(studentId)
	if err != nil || !exists {
		return nil, fmt.Errorf("student not found")
	}

	exists, err = s.courseRepository.ExistsById(courseId)
	if err != nil || !exists {
		return nil, fmt.Errorf("course not found")
	}

	student := entity.Student{ID: studentId}
	course := entity.Course{ID: courseId}

	err = dbcontext.DB.Model(&student).Association("Courses").Append(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to add course to student: %w", err)
	}

	return s.FindStudentById(studentId)
}

func (s *service) Count() (int, error) {
	return s.studentRepository.Count()
}
