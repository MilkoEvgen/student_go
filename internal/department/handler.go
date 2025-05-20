package department

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

type DepartmentHandler struct {
	Service Service
}

func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{
		Service: NewDepartmentService(NewDepartmentRepository(), teacher.NewTeacherRepository()),
	}
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req request.DepartmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in CreateDepartment", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("CreateDepartment called", zap.String("name", req.Name))

	deptResp, err := h.Service.CreateDepartment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save department"})
		return
	}

	c.JSON(http.StatusCreated, deptResp)
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	var req request.DepartmentRequest

	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid department ID in UpdateDepartment", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department ID"})
		return
	}
	id := uint(parsedID)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Warn("Invalid request in UpdateDepartment", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Log.Info("UpdateDepartment called", zap.Uint("id", id), zap.String("name", req.Name))

	deptResp, err := h.Service.UpdateDepartment(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update department"})
		return
	}

	c.JSON(http.StatusOK, deptResp)
}

func (h *DepartmentHandler) FindDepartmentById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid department ID in FindDepartmentById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department ID"})
		return
	}
	id := uint(parsedID)

	log.Log.Info("FindDepartmentById called", zap.Uint("id", id))

	deptResp, err := h.Service.FindDepartmentById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, deptResp)
}

func (h *DepartmentHandler) FindAllDepartments(c *gin.Context) {
	count, err := h.Service.Count()
	if err != nil {
		log.Log.Error("Failed to count departments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count departments"})
		return
	}

	pages := pagination.NewFromRequest(c.Request, count)

	log.Log.Info("FindAllDepartments called",
		zap.Int("page", pages.Page),
		zap.Int("per_page", pages.PerPage),
		zap.Int("total_count", pages.TotalCount),
	)

	depts, err := h.Service.FindAllDepartments(pages.Page, pages.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get departments"})
		return
	}

	pages.Items = depts
	c.JSON(http.StatusOK, pages)
}

func (h *DepartmentHandler) DeleteDepartmentById(c *gin.Context) {
	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid department ID in DeleteDepartmentById", zap.String("id", idParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department ID"})
		return
	}
	id := uint(parsedID)

	log.Log.Info("DeleteDepartmentById called", zap.Uint("id", id))

	err = h.Service.DeleteDepartmentById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DepartmentHandler) DepartmentSetTeacher(c *gin.Context) {
	departmentIdParam := c.Param("departmentId")
	parsedDepartmentID, err := strconv.ParseUint(departmentIdParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid department ID in DepartmentSetTeacher", zap.String("department_id", departmentIdParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department ID"})
		return
	}

	teacherIdParam := c.Param("teacherId")
	parsedTeacherID, err := strconv.ParseUint(teacherIdParam, 10, 32)
	if err != nil {
		log.Log.Warn("Invalid teacher ID in DepartmentSetTeacher", zap.String("teacher_id", teacherIdParam), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher ID"})
		return
	}

	departmentId := uint(parsedDepartmentID)
	teacherId := uint(parsedTeacherID)

	log.Log.Info("DepartmentSetTeacher called",
		zap.Uint("department_id", departmentId),
		zap.Uint("teacher_id", teacherId),
	)

	departmentResp, err := h.Service.DepartmentSetTeacher(departmentId, teacherId)
	if err != nil {
		if err.Error() == "department not found" || err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, departmentResp)
}
