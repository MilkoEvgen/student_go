package entity

type Department struct {
	ID                 uint `gorm:"primaryKey"`
	Name               string
	HeadOfDepartmentID *uint
	HeadOfDepartment   *Teacher `gorm:"foreignKey:HeadOfDepartmentID"`
}
