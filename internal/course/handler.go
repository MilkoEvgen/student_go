package course

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"student_go/internal/dto/request"
	"student_go/internal/teacher"
	"student_go/pkg/log"
	"student_go/pkg/pagination"
)

type Handler struct {
	Service Service
}

func NewCourseHandler() *Handler {
	return &Handler{
		Service: NewCourseService(NewCourseRepository(), teacher.NewTeacherRepository()),
	}
}

func (h *Handler) CreateCourse(c *gin.Context) {
	var req request.CourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in CreateCourse", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("CreateCourse called", zap.String("title", req.Title))

	courseResp, err := h.Service.CreateCourse(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save course"})
		return
	}

	c.JSON(http.StatusCreated, courseResp)
}

func (h *Handler) UpdateCourse(c *gin.Context) {
	var req request.CourseRequest

	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid course ID in UpdateCourse", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}
	id := uint(parsedID)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in UpdateCourse", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("UpdateCourse called", zap.Uint("id", id), zap.String("title", req.Title))

	courseResp, err := h.Service.UpdateCourse(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update course"})
		return
	}

	c.JSON(http.StatusOK, courseResp)
}

func (h *Handler) FindCourseById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid course ID in FindCourseById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}
	id := uint(parsedID)

	log.Log.Info("FindCourseById called", zap.Uint("id", id))

	courseResp, err := h.Service.FindCourseById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, courseResp)
}

func (h *Handler) FindAllCourses(c *gin.Context) {
	count, err := h.Service.Count()
	if err != nil {
		log.Log.Error("Failed to count courses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count courses"})
		return
	}

	pages := pagination.NewFromRequest(c.Request, count)

	log.Log.Info("FindAllCourses called",
		zap.Int("page", pages.Page),
		zap.Int("per_page", pages.PerPage),
		zap.Int("total_count", pages.TotalCount),
	)

	courses, err := h.Service.FindAllCourse(pages.Page, pages.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all courses"})
		return
	}

	pages.Items = courses
	c.JSON(http.StatusOK, pages)
}

func (h *Handler) DeleteCourseById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid course ID in DeleteCourseById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}
	id := uint(parsedID)

	log.Log.Info("DeleteCourseById called", zap.Uint("id", id))

	err = h.Service.DeleteCourseById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) SetTeacherToCourse(c *gin.Context) {
	courseIdParam := c.Param("courseId")
	parsedCourseID, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid course ID in SetTeacherToCourse", zap.String("course_id", courseIdParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}

	teacherIdParam := c.Param("teacherId")
	parsedTeacherID, err := strconv.ParseUint(teacherIdParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid teacher ID in SetTeacherToCourse", zap.String("teacher_id", teacherIdParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher ID"})
		return
	}

	courseId := uint(parsedCourseID)
	teacherId := uint(parsedTeacherID)

	log.Log.Info("SetTeacherToCourse called",
		zap.Uint("course_id", courseId),
		zap.Uint("teacher_id", teacherId),
	)

	courseResp, err := h.Service.SetTeacherToCourse(courseId, teacherId)
	if err != nil {
		if err.Error() == "course not found" || err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, courseResp)
}
