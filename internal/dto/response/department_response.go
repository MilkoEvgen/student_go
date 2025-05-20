package response

type DepartmentResponse struct {
	ID               uint             `json:"id"`
	Name             string           `json:"name"`
	HeadOfDepartment *TeacherResponse `json:"headOfDepartment"`
}
