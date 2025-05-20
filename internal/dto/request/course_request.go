package request

type CourseRequest struct {
	Title string `json:"title" binding:"required"`
}
