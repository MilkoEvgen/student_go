package student

import (
	"student_go/internal/entity"
	"student_go/pkg/dbcontext"
)

type Repository interface {
	ExistsById(id uint) (bool, error)
	Save(student *entity.Student) (*entity.Student, error)
	Update(student *entity.Student) (*entity.Student, error)
	FindById(id uint) (*entity.Student, error)
	FindAll(page, limit int) ([]entity.Student, error)
	DeleteById(id uint) error
	Count() (int, error)
}

type repository struct{}

func NewStudentRepository() Repository {
	return &repository{}
}

func (r *repository) ExistsById(id uint) (bool, error) {
	var exists bool
	err := dbcontext.DB.
		Model(&entity.Student{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists, err
}

func (r *repository) Save(student *entity.Student) (*entity.Student, error) {
	err := dbcontext.DB.Create(student).Error
	return student, err
}

func (r *repository) Update(student *entity.Student) (*entity.Student, error) {
	err := dbcontext.DB.Save(&student).Error
	if err != nil {
		return nil, err
	}

	var updatedStudent *entity.Student

	err = dbcontext.DB.
		Preload("Courses").
		Preload("Courses.Teacher").
		First(&updatedStudent, student.ID).Error

	if err != nil {
		return nil, err
	}
	return updatedStudent, err
}

func (r *repository) FindById(id uint) (*entity.Student, error) {
	var student entity.Student
	result := dbcontext.DB.
		Preload("Courses").
		Preload("Courses.Teacher").
		First(&student, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}

func (r *repository) FindAll(page, limit int) ([]entity.Student, error) {
	var students []entity.Student

	offset := (page - 1) * limit

	result := dbcontext.DB.
		Preload("Courses").
		Preload("Courses.Teacher").
		Limit(limit).
		Offset(offset).
		Find(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	return students, nil
}

func (r *repository) DeleteById(id uint) error {
	result := dbcontext.DB.Delete(&entity.Student{}, id)

	return result.Error
}

func (r *repository) Count() (int, error) {
	var count int64
	err := dbcontext.DB.Model(&entity.Student{}).Count(&count).Error
	return int(count), err
}
