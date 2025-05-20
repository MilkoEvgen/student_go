package request

type DepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}
