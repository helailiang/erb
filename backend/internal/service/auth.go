package service

import (
	"erb/internal/config"
	"erb/internal/model"
	"erb/internal/repository"
	"erb/pkg/jwt"
	"erb/pkg/utils"
	"errors"
	"time"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string      `json:"token"`
	User     *model.User `json:"user"`
	ExpireAt time.Time   `json:"expire_at"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest, ip string) (*LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查账号状态
	if user.Status != "normal" {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ip
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新登录信息失败")
	}

	// 生成Token
	token, err := jwt.GenerateToken(user.ID, user.Username, s.cfg.JWT.Secret, s.cfg.JWT.ExpireTime)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	expireAt := time.Now().Add(time.Duration(s.cfg.JWT.ExpireTime) * time.Hour)

	// 隐藏密码
	user.Password = ""

	return &LoginResponse{
		Token:    token,
		User:     user,
		ExpireAt: expireAt,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 隐藏密码
	user.Password = ""

	return user, nil
}

