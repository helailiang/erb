package service

import (
	"erb/internal/model"
	"erb/internal/repository"
	"erb/pkg/utils"
	"errors"
	"fmt"
	"time"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	RealName   string    `json:"real_name" binding:"required"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	EmployeeNo string    `json:"employee_no"`
	JoinDate   *time.Time `json:"join_date"`
	ManagerID  *uint     `json:"manager_id"`
	RoleIDs    []uint    `json:"role_ids"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	RealName   string    `json:"real_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	EmployeeNo string    `json:"employee_no"`
	JoinDate   *time.Time `json:"join_date"`
	ManagerID  *uint     `json:"manager_id"`
	Status     string    `json:"status"`
	Avatar     string    `json:"avatar"`
}

// ListUserRequest 用户列表请求
type ListUserRequest struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	Username   string `form:"username"`
	RealName   string `form:"real_name"`
	Department string `form:"department"`
	Status     string `form:"status"`
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 验证密码强度
	if len(req.Password) < 8 {
		return nil, errors.New("密码长度至少8位")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	now := time.Now()
	user := &model.User{
		Username:          req.Username,
		Password:          hashedPassword,
		RealName:          req.RealName,
		Email:             req.Email,
		Phone:             req.Phone,
		Department:        req.Department,
		Position:          req.Position,
		EmployeeNo:        req.EmployeeNo,
		JoinDate:          req.JoinDate,
		ManagerID:         req.ManagerID,
		Status:            "normal", // 默认正常状态
		PasswordChangedAt: &now,
	}

	// 创建用户
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return nil, fmt.Errorf("分配角色失败: %w", err)
		}
	}

	// 重新获取用户信息（包含角色）
	user, err = s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	// 隐藏密码
	user.Password = ""

	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(userID uint, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Department != "" {
		user.Department = req.Department
	}
	if req.Position != "" {
		user.Position = req.Position
	}
	if req.EmployeeNo != "" {
		user.EmployeeNo = req.EmployeeNo
	}
	if req.JoinDate != nil {
		user.JoinDate = req.JoinDate
	}
	if req.ManagerID != nil {
		user.ManagerID = req.ManagerID
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 重新获取用户信息
	user, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 隐藏密码
	user.Password = ""

	return user, nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}

// GetUser 获取用户详情
func (s *UserService) GetUser(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 隐藏密码
	user.Password = ""

	return user, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(req *ListUserRequest) ([]model.User, int64, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.Username != "" {
		conditions["username LIKE ?"] = fmt.Sprintf("%%%s%%", req.Username)
	}
	if req.RealName != "" {
		conditions["real_name LIKE ?"] = fmt.Sprintf("%%%s%%", req.RealName)
	}
	if req.Department != "" {
		conditions["department = ?"] = req.Department
	}
	if req.Status != "" {
		conditions["status = ?"] = req.Status
	}

	users, total, err := s.userRepo.List(req.Page, req.PageSize, conditions)
	if err != nil {
		return nil, 0, err
	}

	// 隐藏所有用户的密码
	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(userID uint, newPassword string) error {
	// 验证密码强度
	if len(newPassword) < 8 {
		return errors.New("密码长度至少8位")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	now := time.Now()
	user.Password = hashedPassword
	user.PasswordChangedAt = &now

	return s.userRepo.Update(user)
}



