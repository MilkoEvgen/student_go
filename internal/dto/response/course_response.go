package response

type CourseResponse struct {
	ID       uint              `json:"id"`
	Title    string            `json:"title"`
	Teacher  *TeacherResponse  `json:"teacher"`
	Students []StudentResponse `json:"students"`
}
