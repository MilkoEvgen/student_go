package course

import (
	"fmt"
	"go.uber.org/zap"
	"student_go/internal/dto/request"
	response3 "student_go/internal/dto/response"
	"student_go/internal/entity"
	"student_go/internal/teacher"
	"student_go/pkg/dbcontext"
	"student_go/pkg/log"
)

type Service interface {
	CreateCourse(input request.CourseRequest) (*response3.CourseResponse, error)
	UpdateCourse(id uint, input request.CourseRequest) (*response3.CourseResponse, error)
	FindCourseById(id uint) (*response3.CourseResponse, error)
	FindAllCourse(page, limit int) ([]*response3.CourseResponse, error)
	DeleteCourseById(id uint) error
	SetTeacherToCourse(courseId uint, teacherId uint) (*response3.CourseResponse, error)
	Count() (int, error)
}

type service struct {
	courseRepository  Repository
	teacherRepository teacher.Repository
}

func NewCourseService(
	courseRepository Repository,
	teacherRepository teacher.Repository) Service {
	return &service{
		courseRepository:  courseRepository,
		teacherRepository: teacherRepository,
	}
}

func (s *service) CreateCourse(input request.CourseRequest) (*response3.CourseResponse, error) {
	log.Log.Info("CreateCourse (service) called", zap.String("title", input.Title))

	course := entity.Course{
		Title: input.Title,
	}
	savedCourse, err := s.courseRepository.Save(&course)
	if err != nil {
		return nil, err
	}

	resp := &response3.CourseResponse{
		ID:    savedCourse.ID,
		Title: savedCourse.Title,
	}
	return resp, nil
}

func (s *service) UpdateCourse(id uint, input request.CourseRequest) (*response3.CourseResponse, error) {
	log.Log.Info("UpdateCourse (service) called",
		zap.Uint("id", id),
		zap.String("title", input.Title),
	)

	course := entity.Course{
		ID:    id,
		Title: input.Title,
	}
	updatedCourse, err := s.courseRepository.Update(&course)
	if err != nil {
		return nil, err
	}

	var teacherResp *response3.TeacherResponse
	if updatedCourse.Teacher != nil {
		teacherResp = &response3.TeacherResponse{
			ID:   updatedCourse.Teacher.ID,
			Name: updatedCourse.Teacher.Name,
		}
	}

	studentsResp := make([]response3.StudentResponse, 0, len(updatedCourse.Students))
	for _, student := range updatedCourse.Students {
		studentsResp = append(studentsResp, response3.StudentResponse{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
		})
	}

	courseResp := &response3.CourseResponse{
		ID:       course.ID,
		Title:    course.Title,
		Teacher:  teacherResp,
		Students: studentsResp,
	}
	return courseResp, nil
}

func (s *service) FindCourseById(id uint) (*response3.CourseResponse, error) {
	log.Log.Info("FindCourseById (service) called", zap.Uint("id", id))

	course, err := s.courseRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	var teacherResp *response3.TeacherResponse
	if course.Teacher != nil {
		teacherResp = &response3.TeacherResponse{
			ID:   course.Teacher.ID,
			Name: course.Teacher.Name,
		}
	}

	studentsResp := make([]response3.StudentResponse, 0, len(course.Students))
	for _, student := range course.Students {
		studentsResp = append(studentsResp, response3.StudentResponse{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
		})
	}

	courseResp := &response3.CourseResponse{
		ID:       course.ID,
		Title:    course.Title,
		Teacher:  teacherResp,
		Students: studentsResp,
	}
	return courseResp, nil
}

func (s *service) FindAllCourse(page, limit int) ([]*response3.CourseResponse, error) {
	log.Log.Info("FindAllCourse (service) called", zap.Int("page", page), zap.Int("limit", limit))

	courses, err := s.courseRepository.FindAll(page, limit)
	if err != nil {
		return nil, err
	}

	var courseResponses []*response3.CourseResponse
	for _, course := range courses {
		var teacherResp *response3.TeacherResponse
		if course.Teacher != nil {
			teacherResp = &response3.TeacherResponse{
				ID:   course.Teacher.ID,
				Name: course.Teacher.Name,
			}
		}

		studentsResp := make([]response3.StudentResponse, 0, len(course.Students))
		for _, student := range course.Students {
			studentsResp = append(studentsResp, response3.StudentResponse{
				ID:    student.ID,
				Name:  student.Name,
				Email: student.Email,
			})
		}

		resp := &response3.CourseResponse{
			ID:       course.ID,
			Title:    course.Title,
			Teacher:  teacherResp,
			Students: studentsResp,
		}
		courseResponses = append(courseResponses, resp)
	}

	return courseResponses, nil
}

func (s *service) DeleteCourseById(id uint) error {
	log.Log.Info("DeleteCourseById (service) called", zap.Uint("id", id))
	return s.courseRepository.DeleteById(id)
}

func (s *service) SetTeacherToCourse(courseId uint, teacherId uint) (*response3.CourseResponse, error) {
	log.Log.Info("SetTeacherToCourse (service) called",
		zap.Uint("course_id", courseId),
		zap.Uint("teacher_id", teacherId),
	)

	exists, err := s.courseRepository.ExistsById(courseId)
	if err != nil || !exists {
		return nil, fmt.Errorf("course not found")
	}

	exists, err = s.teacherRepository.ExistsById(teacherId)
	if err != nil || !exists {
		return nil, fmt.Errorf("teacher not found")
	}

	err = dbcontext.DB.Model(&entity.Course{}).Where("id = ?", courseId).Update("teacher_id", teacherId).Error
	if err != nil {
		return nil, fmt.Errorf("failed to assign teacher to course: %w", err)
	}

	return s.FindCourseById(courseId)
}

func (s *service) Count() (int, error) {
	return s.courseRepository.Count()
}
