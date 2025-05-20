package response

type TeacherResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Courses     []CourseResponse     `json:"courses"`
	Departments []DepartmentResponse `json:"departments"`
}
