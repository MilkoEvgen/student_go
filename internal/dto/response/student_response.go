package response

type StudentResponse struct {
	ID      uint             `json:"id"`
	Name    string           `json:"name"`
	Email   string           `json:"email"`
	Courses []CourseResponse `json:"courses"`
}
