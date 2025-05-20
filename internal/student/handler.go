package student

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"student_go/internal/course"
	"student_go/internal/dto/request"
	"student_go/pkg/log"
	"student_go/pkg/pagination"
)

type StudentHandler struct {
	Service Service
}

func NewStudentHandler() *StudentHandler {
	return &StudentHandler{
		Service: NewStudentService(NewStudentRepository(), course.NewCourseRepository()),
	}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req request.StudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in CreateStudent", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("CreateStudent called",
		zap.String("name", req.Name),
		zap.String("email", req.Email),
	)

	studentResp, err := h.Service.CreateStudent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save student"})
		return
	}

	c.JSON(http.StatusCreated, studentResp)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	var req request.StudentRequest

	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid ID in UpdateStudent", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	id := uint(parsedID)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in UpdateStudent", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("UpdateStudent called",
		zap.String("name", req.Name),
		zap.String("email", req.Email),
	)

	studentResp, err := h.Service.UpdateStudent(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
		return
	}

	c.JSON(http.StatusOK, studentResp)
}

func (h *StudentHandler) FindStudentById(c *gin.Context) {
	idParam := c.Param("id")

	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid ID in FindStudentById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	log.Log.Info("FindStudentById called", zap.String("id", idParam))

	id := uint(parsedID)
	studentResp, err := h.Service.FindStudentById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
	}

	c.JSON(http.StatusOK, studentResp)
}

func (h *StudentHandler) FindAllStudents(c *gin.Context) {
	count, err := h.Service.Count()
	if err != nil {
		log.Log.Error("Failed to count students", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count students"})
		return
	}

	pages := pagination.NewFromRequest(c.Request, count)

	log.Log.Info("FindAllStudents called",
		zap.Int("page", pages.Page),
		zap.Int("per_page", pages.PerPage),
		zap.Int("total_count", pages.TotalCount),
	)

	studentResp, err := h.Service.FindAllStudent(pages.Page, pages.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get students"})
		return
	}

	pages.Items = studentResp
	c.JSON(http.StatusOK, pages)
}

func (h *StudentHandler) FindAllCoursesByStudentId(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid ID in FindAllCoursesByStudentId", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	log.Log.Info("FindStudentById called", zap.String("id", idParam))

	id := uint(parsedID)

	studentResp, err := h.Service.FindStudentById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all courses"})
		return
	}

	c.JSON(http.StatusOK, studentResp.Courses)
}

func (h *StudentHandler) DeleteStudentById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid ID in DeleteStudentById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	log.Log.Info("DeleteStudentById called", zap.String("id", idParam))

	id := uint(parsedID)

	err = h.Service.DeleteStudentById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *StudentHandler) StudentAddCourse(c *gin.Context) {
	studentIdParam := c.Param("studentId")
	parsedStudentID, err := strconv.ParseUint(studentIdParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid student ID",
			zap.String("student_id", studentIdParam),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	courseIdParam := c.Param("courseId")
	parsedCourseID, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid course ID",
			zap.String("course_id", courseIdParam),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}

	studentId := uint(parsedStudentID)
	courseId := uint(parsedCourseID)

	log.Log.Info("StudentAddCourse called",
		zap.Uint("student_id", studentId),
		zap.Uint("course_id", courseId),
	)

	studentResp, err := h.Service.AddCourseToStudent(studentId, courseId)
	if err != nil {
		if err.Error() == "student not found" || err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, studentResp)
}
