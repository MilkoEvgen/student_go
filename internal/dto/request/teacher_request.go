package request

type TeacherRequest struct {
	Name string `json:"name" binding:"required"`
}
