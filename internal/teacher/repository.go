package teacher

import (
	"student_go/internal/entity"
	"student_go/pkg/dbcontext"
)

type Repository interface {
	ExistsById(id uint) (bool, error)
	Save(teacher *entity.Teacher) (*entity.Teacher, error)
	Update(teacher *entity.Teacher) (*entity.Teacher, error)
	FindById(id uint) (*entity.Teacher, error)
	FindAll(page, limit int) ([]entity.Teacher, error)
	DeleteById(id uint) error
	Count() (int, error)
}

type repository struct{}

func NewTeacherRepository() Repository {
	return &repository{}
}

func (r *repository) ExistsById(id uint) (bool, error) {
	var exists bool
	err := dbcontext.DB.
		Model(&entity.Teacher{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists, err
}

func (r *repository) Save(teacher *entity.Teacher) (*entity.Teacher, error) {
	err := dbcontext.DB.Create(teacher).Error
	return teacher, err
}

func (r *repository) Update(teacher *entity.Teacher) (*entity.Teacher, error) {
	err := dbcontext.DB.Save(&teacher).Error

	if err != nil {
		return nil, err
	}

	var updatedTeacher *entity.Teacher

	err = dbcontext.DB.
		Preload("Courses").
		Preload("Departments").
		First(&updatedTeacher, teacher.ID).Error

	return updatedTeacher, err
}

func (r *repository) FindById(id uint) (*entity.Teacher, error) {
	var teacher entity.Teacher
	result := dbcontext.DB.
		Preload("Courses").
		Preload("Departments").
		First(&teacher, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &teacher, nil
}

func (r *repository) FindAll(page, limit int) ([]entity.Teacher, error) {
	var teachers []entity.Teacher

	offset := (page - 1) * limit

	result := dbcontext.DB.
		Preload("Courses").
		Preload("Departments").
		Offset(offset).
		Find(&teachers)

	if result.Error != nil {
		return nil, result.Error
	}

	return teachers, nil
}

func (r *repository) DeleteById(id uint) error {
	result := dbcontext.DB.Delete(&entity.Teacher{}, id)

	return result.Error
}

func (r *repository) Count() (int, error) {
	var count int64
	err := dbcontext.DB.Model(&entity.Teacher{}).Count(&count).Error
	return int(count), err
}
