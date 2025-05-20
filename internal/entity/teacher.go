package entity

type Teacher struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Courses     []Course     `gorm:"foreignKey:TeacherID"`
	Departments []Department `gorm:"foreignKey:HeadOfDepartmentID"`
}
