package entity

type Student struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Email   string
	Courses []Course `gorm:"many2many:course_student"`
}
