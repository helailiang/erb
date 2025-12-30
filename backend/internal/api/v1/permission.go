package v1

import (
	"erb/internal/api/middleware"
	"erb/internal/config"
	"erb/internal/repository"
	"erb/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionRepo *repository.PermissionRepository
	cfg            *config.Config
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(db *gorm.DB, cfg *config.Config) *PermissionHandler {
	permissionRepo := repository.NewPermissionRepository(db)
	return &PermissionHandler{
		permissionRepo: permissionRepo,
		cfg:            cfg,
	}
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 获取所有权限列表
// @Tags 权限管理
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.Permission}
// @Router /api/v1/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.permissionRepo.GetAll()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, permissions)
}

// RegisterRoutes 注册路由
func (h *PermissionHandler) RegisterRoutes(r *gin.RouterGroup, cfg *config.Config) {
	permissions := r.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(cfg))
	{
		permissions.GET("", h.ListPermissions)
	}
}


