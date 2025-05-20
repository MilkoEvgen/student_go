package course

import (
	"student_go/internal/entity"
	"student_go/pkg/dbcontext"
)

type Repository interface {
	ExistsById(id uint) (bool, error)
	Save(course *entity.Course) (*entity.Course, error)
	Update(course *entity.Course) (*entity.Course, error)
	FindById(id uint) (*entity.Course, error)
	FindAll(page, limit int) ([]entity.Course, error)
	DeleteById(id uint) error
	Count() (int, error)
}

type repository struct{}

func NewCourseRepository() Repository {
	return &repository{}
}

func (r *repository) ExistsById(id uint) (bool, error) {
	var exists bool
	err := dbcontext.DB.
		Model(&entity.Course{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists, err
}

func (r *repository) Save(course *entity.Course) (*entity.Course, error) {
	err := dbcontext.DB.Create(course).Error
	return course, err
}

func (r *repository) Update(course *entity.Course) (*entity.Course, error) {
	err := dbcontext.DB.Model(&entity.Course{}).
		Where("id = ?", course.ID).
		Updates(map[string]interface{}{
			"title": course.Title,
		}).Error

	if err != nil {
		return nil, err
	}

	var updated *entity.Course

	err = dbcontext.DB.
		Preload("Students").
		Preload("Teacher").
		First(&updated, course.ID).Error

	if err != nil {
		return nil, err
	}

	return updated, err
}

func (r *repository) FindById(id uint) (*entity.Course, error) {
	var course entity.Course
	result := dbcontext.DB.
		Preload("Students").
		Preload("Teacher").
		First(&course, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

func (r *repository) FindAll(page, limit int) ([]entity.Course, error) {
	var courses []entity.Course

	offset := (page - 1) * limit

	result := dbcontext.DB.
		Preload("Students").
		Preload("Teacher").
		Offset(offset).
		Find(&courses)

	if result.Error != nil {
		return nil, result.Error
	}

	return courses, nil
}

func (r *repository) DeleteById(id uint) error {
	result := dbcontext.DB.Delete(&entity.Course{}, id)

	return result.Error
}

func (r *repository) Count() (int, error) {
	var count int64
	err := dbcontext.DB.Model(&entity.Course{}).Count(&count).Error
	return int(count), err
}
