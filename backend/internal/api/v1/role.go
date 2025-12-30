package v1

import (
	"erb/internal/api/middleware"
	"erb/internal/config"
	"erb/internal/repository"
	"erb/internal/service"
	"erb/pkg/response"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *service.RoleService
	cfg         *config.Config
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(db *gorm.DB, cfg *config.Config) *RoleHandler {
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	return &RoleHandler{
		roleService: roleService,
		cfg:         cfg,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body service.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=model.Role}
// @Failure 400 {object} response.Response
// @Router /api/v1/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param role body service.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=model.Role}
// @Failure 400 {object} response.Response
// @Router /api/v1/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色（软删除）
// @Tags 角色管理
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详细信息
// @Tags 角色管理
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=model.Role}
// @Failure 404 {object} response.Response
// @Router /api/v1/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.roleService.GetRole(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, role)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取角色列表
// @Tags 角色管理
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param name query string false "角色名"
// @Param status query string false "状态"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/v1/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	var req service.ListRoleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	roles, total, err := h.roleService.ListRoles(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      roles,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// RegisterRoutes 注册路由
func (h *RoleHandler) RegisterRoutes(r *gin.RouterGroup, cfg *config.Config) {
	roles := r.Group("/roles")
	roles.Use(middleware.AuthMiddleware(cfg))
	{
		roles.POST("", h.CreateRole)
		roles.GET("", h.ListRoles)
		roles.GET("/:id", h.GetRole)
		roles.PUT("/:id", h.UpdateRole)
		roles.DELETE("/:id", h.DeleteRole)
	}
}



