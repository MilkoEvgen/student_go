package app

import (
	"github.com/gin-gonic/gin"
	"student_go/internal/config"
	"student_go/internal/course"
	"student_go/internal/department"
	"student_go/internal/student"
	"student_go/internal/teacher"
	"student_go/pkg/dbcontext"
	"student_go/pkg/migration"
)

func SetupApp() (*gin.Engine, error) {
	if err := config.Load(); err != nil {
		return nil, err
	}

	if err := dbcontext.Connect(); err != nil {
		return nil, err
	}

	migration.ApplyMigrations()

	r := gin.Default()

	studentHandler := student.NewStudentHandler()
	teacherHandler := teacher.NewTeacherHandler()
	courseHandler := course.NewCourseHandler()
	departmentHandler := department.NewDepartmentHandler()

	r.POST("/api/v1/students", studentHandler.CreateStudent)
	r.PATCH("/api/v1/students/:id", studentHandler.UpdateStudent)
	r.GET("/api/v1/students/:id", studentHandler.FindStudentById)
	r.GET("/api/v1/students", studentHandler.FindAllStudents)
	r.GET("/api/v1/students/:id/courses", studentHandler.FindAllCoursesByStudentId)
	r.DELETE("/api/v1/students/:id", studentHandler.DeleteStudentById)
	r.POST("/api/v1/students/:studentId/courses/:courseId", studentHandler.StudentAddCourse)

	r.POST("/api/v1/courses", courseHandler.CreateCourse)
	r.PATCH("/api/v1/courses/:id", courseHandler.UpdateCourse)
	r.GET("/api/v1/courses/:id", courseHandler.FindCourseById)
	r.GET("/api/v1/courses", courseHandler.FindAllCourses)
	r.DELETE("/api/v1/courses/:id", courseHandler.DeleteCourseById)
	r.POST("/api/v1/courses/:courseId/teacher/:teacherId", courseHandler.SetTeacherToCourse)

	r.POST("/api/v1/teachers", teacherHandler.CreateTeacher)
	r.PATCH("/api/v1/teachers/:id", teacherHandler.UpdateTeacher)
	r.GET("/api/v1/teachers/:id", teacherHandler.FindTeacherById)
	r.GET("/api/v1/teachers", teacherHandler.FindAllTeachers)
	r.DELETE("/api/v1/teachers/:id", teacherHandler.DeleteTeacherById)

	r.POST("/api/v1/departments", departmentHandler.CreateDepartment)
	r.PATCH("/api/v1/departments/:id", departmentHandler.UpdateDepartment)
	r.GET("/api/v1/departments/:id", departmentHandler.FindDepartmentById)
	r.GET("/api/v1/departments", departmentHandler.FindAllDepartments)
	r.DELETE("/api/v1/departments/:id", departmentHandler.DeleteDepartmentById)
	r.POST("/api/v1/departments/:departmentId/teacher/:teacherId", departmentHandler.DepartmentSetTeacher)

	return r, nil
}
