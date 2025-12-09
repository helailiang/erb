package v1

import (
	"erb/internal/api/middleware"
	"erb/internal/config"
	"erb/internal/repository"
	"erb/internal/service"
	"erb/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ip := c.ClientIP()
	resp, err := h.authService.Login(&req, ip)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的信息
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=model.User}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.authService.GetUserInfo(userID.(uint))
	if err != nil {
		response.Error(c, 404, err.Error())
		return
	}

	response.Success(c, user)
}

// RegisterRoutes 注册路由
func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, cfg *config.Config) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.GET("/me", middleware.AuthMiddleware(cfg), h.GetUserInfo)
	}
}
