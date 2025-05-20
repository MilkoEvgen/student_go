package department

import (
	"student_go/internal/entity"
	"student_go/pkg/dbcontext"
)

type Repository interface {
	ExistsById(id uint) (bool, error)
	Save(department *entity.Department) (*entity.Department, error)
	Update(department *entity.Department) (*entity.Department, error)
	FindById(id uint) (*entity.Department, error)
	FindAll(page, limit int) ([]entity.Department, error)
	DeleteById(id uint) error
	Count() (int, error)
}

type repository struct{}

func NewDepartmentRepository() Repository {
	return &repository{}
}

func (r *repository) ExistsById(id uint) (bool, error) {
	var exists bool
	err := dbcontext.DB.
		Model(&entity.Department{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists, err
}

func (r *repository) Save(department *entity.Department) (*entity.Department, error) {
	err := dbcontext.DB.Create(department).Error
	return department, err
}

func (r *repository) Update(department *entity.Department) (*entity.Department, error) {
	err := dbcontext.DB.Model(&entity.Department{}).
		Where("id = ?", department.ID).
		Updates(map[string]interface{}{
			"name": department.Name,
		}).Error

	if err != nil {
		return nil, err
	}

	var updatedDepartment *entity.Department

	err = dbcontext.DB.
		Preload("HeadOfDepartment").
		First(&updatedDepartment, department.ID).Error

	if err != nil {
		return nil, err
	}

	return updatedDepartment, err
}

func (r *repository) FindById(id uint) (*entity.Department, error) {
	var department entity.Department
	result := dbcontext.DB.
		Preload("HeadOfDepartment").
		First(&department, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &department, nil
}

func (r *repository) FindAll(page, limit int) ([]entity.Department, error) {
	var departments []entity.Department

	offset := (page - 1) * limit

	result := dbcontext.DB.
		Preload("HeadOfDepartment").
		Offset(offset).
		Find(&departments)

	if result.Error != nil {
		return nil, result.Error
	}

	return departments, nil
}

func (r *repository) DeleteById(id uint) error {
	result := dbcontext.DB.Delete(&entity.Department{}, id)

	return result.Error
}

func (r *repository) Count() (int, error) {
	var count int64
	err := dbcontext.DB.Model(&entity.Department{}).Count(&count).Error
	return int(count), err
}
