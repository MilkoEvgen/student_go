package entity

type Course struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	TeacherID *uint
	Students  []Student `gorm:"many2many:course_student"`
	Teacher   *Teacher  `gorm:"foreignKey:TeacherID"`
}
