package department

import (
	"fmt"
	"go.uber.org/zap"
	"student_go/internal/dto/request"
	"student_go/internal/dto/response"
	"student_go/internal/entity"
	"student_go/internal/teacher"
	"student_go/pkg/dbcontext"
	"student_go/pkg/log"
)

type Service interface {
	CreateDepartment(input request.DepartmentRequest) (*response.DepartmentResponse, error)
	UpdateDepartment(id uint, input request.DepartmentRequest) (*response.DepartmentResponse, error)
	FindDepartmentById(id uint) (*response.DepartmentResponse, error)
	FindAllDepartments(page, limit int) ([]*response.DepartmentResponse, error)
	DeleteDepartmentById(id uint) error
	DepartmentSetTeacher(departmentId uint, teacherId uint) (*response.DepartmentResponse, error)
	Count() (int, error)
}

type service struct {
	departmentRepository Repository
	teacherRepository    teacher.Repository
}

func NewDepartmentService(
	departmentRepository Repository,
	teacherRepository teacher.Repository,
) Service {
	return &service{
		departmentRepository: departmentRepository,
		teacherRepository:    teacherRepository,
	}
}

func (s *service) CreateDepartment(input request.DepartmentRequest) (*response.DepartmentResponse, error) {
	log.Log.Info("CreateDepartment (service) called", zap.String("name", input.Name))

	dept := entity.Department{
		Name: input.Name,
	}
	savedDept, err := s.departmentRepository.Save(&dept)
	if err != nil {
		return nil, err
	}

	resp := &response.DepartmentResponse{
		ID:   savedDept.ID,
		Name: savedDept.Name,
	}
	return resp, nil
}

func (s *service) UpdateDepartment(id uint, input request.DepartmentRequest) (*response.DepartmentResponse, error) {
	log.Log.Info("UpdateDepartment (service) called", zap.Uint("id", id), zap.String("name", input.Name))

	dept := entity.Department{
		ID:   id,
		Name: input.Name,
	}
	updatedDept, err := s.departmentRepository.Update(&dept)
	if err != nil {
		return nil, err
	}

	var headOfDepartment *response.TeacherResponse
	if updatedDept.HeadOfDepartment != nil {
		headOfDepartment = &response.TeacherResponse{
			ID:   updatedDept.HeadOfDepartment.ID,
			Name: updatedDept.HeadOfDepartment.Name,
		}
	}

	departmentResp := &response.DepartmentResponse{
		ID:               updatedDept.ID,
		Name:             updatedDept.Name,
		HeadOfDepartment: headOfDepartment,
	}
	return departmentResp, nil
}

func (s *service) FindDepartmentById(id uint) (*response.DepartmentResponse, error) {
	log.Log.Info("FindDepartmentById (service) called", zap.Uint("id", id))

	dept, err := s.departmentRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	var headOfDepartment *response.TeacherResponse
	if dept.HeadOfDepartment != nil {
		headOfDepartment = &response.TeacherResponse{
			ID:   dept.HeadOfDepartment.ID,
			Name: dept.HeadOfDepartment.Name,
		}
	}

	departmentResp := &response.DepartmentResponse{
		ID:               dept.ID,
		Name:             dept.Name,
		HeadOfDepartment: headOfDepartment,
	}
	return departmentResp, nil
}

func (s *service) FindAllDepartments(page, limit int) ([]*response.DepartmentResponse, error) {
	log.Log.Info("FindAllDepartments (service) called", zap.Int("page", page), zap.Int("limit", limit))

	depts, err := s.departmentRepository.FindAll(page, limit)
	if err != nil {
		return nil, err
	}

	var deptResponses []*response.DepartmentResponse
	for _, dept := range depts {
		var headOfDepartment *response.TeacherResponse
		if dept.HeadOfDepartment != nil {
			headOfDepartment = &response.TeacherResponse{
				ID:   dept.HeadOfDepartment.ID,
				Name: dept.HeadOfDepartment.Name,
			}
		}

		departmentResp := &response.DepartmentResponse{
			ID:               dept.ID,
			Name:             dept.Name,
			HeadOfDepartment: headOfDepartment,
		}
		deptResponses = append(deptResponses, departmentResp)
	}

	return deptResponses, nil
}

func (s *service) DeleteDepartmentById(id uint) error {
	log.Log.Info("DeleteDepartmentById (service) called", zap.Uint("id", id))
	return s.departmentRepository.DeleteById(id)
}

func (s *service) DepartmentSetTeacher(departmentId uint, teacherId uint) (*response.DepartmentResponse, error) {
	log.Log.Info("DepartmentSetTeacher (service) called",
		zap.Uint("department_id", departmentId),
		zap.Uint("teacher_id", teacherId),
	)

	exists, err := s.departmentRepository.ExistsById(departmentId)
	if err != nil || !exists {
		return nil, fmt.Errorf("department not found")
	}

	exists, err = s.teacherRepository.ExistsById(teacherId)
	if err != nil || !exists {
		return nil, fmt.Errorf("teacher not found")
	}

	err = dbcontext.DB.Model(&entity.Department{}).
		Where("id = ?", departmentId).
		Update("head_of_department_id", teacherId).Error
	if err != nil {
		return nil, fmt.Errorf("failed to assign teacher to department: %w", err)
	}

	return s.FindDepartmentById(departmentId)
}

func (s *service) Count() (int, error) {
	return s.departmentRepository.Count()
}
