package teacher

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"student_go/internal/dto/request"
	"student_go/pkg/log"
	"student_go/pkg/pagination"
)

type TeacherHandler struct {
	Service Service
}

func NewTeacherHandler() *TeacherHandler {
	return &TeacherHandler{
		Service: NewTeacherService(NewTeacherRepository()),
	}
}

func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var req request.TeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in CreateTeacher", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("CreateTeacher called", zap.String("name", req.Name))

	teacherResp, err := h.Service.CreateTeacher(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save teacher"})
		return
	}
	c.JSON(http.StatusCreated, teacherResp)
}

func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	var req request.TeacherRequest

	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid teacher ID in UpdateTeacher", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher ID"})
		return
	}
	id := uint(parsedID)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in UpdateTeacher", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("UpdateTeacher called", zap.Uint("id", id), zap.String("name", req.Name))

	teacherResp, err := h.Service.UpdateTeacher(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update teacher"})
		return
	}

	c.JSON(http.StatusOK, teacherResp)
}

func (h *TeacherHandler) FindTeacherById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid teacher ID in FindTeacherById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher ID"})
		return
	}
	id := uint(parsedID)

	log.Log.Info("FindTeacherById called", zap.Uint("id", id))

	teacherResp, err := h.Service.FindTeacherById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, teacherResp)
}

func (h *TeacherHandler) FindAllTeachers(c *gin.Context) {
	count, err := h.Service.Count()
	if err != nil {
		log.Log.Error("Failed to count teachers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count teachers"})
		return
	}

	pages := pagination.NewFromRequest(c.Request, count)

	log.Log.Info("FindAllTeachers called",
		zap.Int("page", pages.Page),
		zap.Int("per_page", pages.PerPage),
		zap.Int("total_count", pages.TotalCount),
	)

	teachers, err := h.Service.FindAllTeachers(pages.Page, pages.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all teachers"})
		return
	}
	pages.Items = teachers
	c.JSON(http.StatusOK, pages)
}

func (h *TeacherHandler) DeleteTeacherById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid teacher ID in DeleteTeacherById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher ID"})
		return
	}
	id := uint(parsedID)

	log.Log.Info("DeleteTeacherById called", zap.Uint("id", id))

	err = h.Service.DeleteTeacherById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
