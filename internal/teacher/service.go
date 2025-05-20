package teacher

import (
	"go.uber.org/zap"
	"student_go/internal/dto/request"
	"student_go/internal/dto/response"
	"student_go/internal/entity"
	"student_go/pkg/log"
)

type Service interface {
	CreateTeacher(input request.TeacherRequest) (*response.TeacherResponse, error)
	UpdateTeacher(id uint, input request.TeacherRequest) (*response.TeacherResponse, error)
	FindTeacherById(id uint) (*response.TeacherResponse, error)
	FindAllTeachers(page, limit int) ([]*response.TeacherResponse, error)
	DeleteTeacherById(id uint) error
	Count() (int, error)
}

type service struct {
	repo Repository
}

func NewTeacherService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTeacher(input request.TeacherRequest) (*response.TeacherResponse, error) {
	log.Log.Info("CreateTeacher (service) called", zap.String("name", input.Name))

	teacher := entity.Teacher{
		Name: input.Name,
	}
	savedTeacher, err := s.repo.Save(&teacher)
	if err != nil {
		return nil, err
	}

	resp := &response.TeacherResponse{
		ID:   savedTeacher.ID,
		Name: savedTeacher.Name,
	}
	return resp, nil
}

func (s *service) UpdateTeacher(id uint, input request.TeacherRequest) (*response.TeacherResponse, error) {
	log.Log.Info("UpdateTeacher (service) called", zap.Uint("id", id), zap.String("name", input.Name))

	teacher := entity.Teacher{
		ID:   id,
		Name: input.Name,
	}
	updatedTeacher, err := s.repo.Update(&teacher)
	if err != nil {
		return nil, err
	}

	var coursesResp []response.CourseResponse
	for _, course := range updatedTeacher.Courses {
		courseResp := response.CourseResponse{
			ID:    course.ID,
			Title: course.Title,
		}
		coursesResp = append(coursesResp, courseResp)
	}

	var departmentsResp []response.DepartmentResponse
	for _, department := range updatedTeacher.Departments {
		departmentResp := response.DepartmentResponse{
			ID:   department.ID,
			Name: department.Name,
		}
		departmentsResp = append(departmentsResp, departmentResp)
	}

	teacherResp := &response.TeacherResponse{
		ID:          teacher.ID,
		Name:        teacher.Name,
		Courses:     coursesResp,
		Departments: departmentsResp,
	}
	return teacherResp, nil
}

func (s *service) FindTeacherById(id uint) (*response.TeacherResponse, error) {
	log.Log.Info("FindTeacherById (service) called", zap.Uint("id", id))

	teacher, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	var coursesResp []response.CourseResponse
	for _, course := range teacher.Courses {
		courseResp := response.CourseResponse{
			ID:    course.ID,
			Title: course.Title,
		}
		coursesResp = append(coursesResp, courseResp)
	}

	var departmentsResp []response.DepartmentResponse
	for _, department := range teacher.Departments {
		departmentResp := response.DepartmentResponse{
			ID:   department.ID,
			Name: department.Name,
		}
		departmentsResp = append(departmentsResp, departmentResp)
	}

	teacherResp := &response.TeacherResponse{
		ID:          teacher.ID,
		Name:        teacher.Name,
		Courses:     coursesResp,
		Departments: departmentsResp,
	}
	return teacherResp, nil
}

func (s *service) FindAllTeachers(page, limit int) ([]*response.TeacherResponse, error) {
	log.Log.Info("FindAllTeachers (service) called", zap.Int("page", page), zap.Int("limit", limit))

	teachers, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, err
	}

	var teacherResponses []*response.TeacherResponse
	for _, teacher := range teachers {
		var coursesResp []response.CourseResponse
		for _, course := range teacher.Courses {
			courseResp := response.CourseResponse{
				ID:    course.ID,
				Title: course.Title,
			}
			coursesResp = append(coursesResp, courseResp)
		}

		var departmentsResp []response.DepartmentResponse
		for _, department := range teacher.Departments {
			departmentResp := response.DepartmentResponse{
				ID:   department.ID,
				Name: department.Name,
			}
			departmentsResp = append(departmentsResp, departmentResp)
		}

		teacherResp := &response.TeacherResponse{
			ID:          teacher.ID,
			Name:        teacher.Name,
			Courses:     coursesResp,
			Departments: departmentsResp,
		}
		teacherResponses = append(teacherResponses, teacherResp)
	}

	return teacherResponses, nil
}

func (s *service) DeleteTeacherById(id uint) error {
	log.Log.Info("DeleteTeacherById (service) called", zap.Uint("id", id))
	return s.repo.DeleteById(id)
}

func (s *service) Count() (int, error) {
	return s.repo.Count()
}
