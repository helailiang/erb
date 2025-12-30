package service

import (
	"erb/internal/model"
	"erb/internal/repository"
	"errors"
	"fmt"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo *repository.RoleRepository
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	PermissionIDs []uint `json:"permission_ids"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	PermissionIDs []uint `json:"permission_ids"`
}

// ListRoleRequest 角色列表请求
type ListRoleRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Name     string `form:"name"`
	Status   string `form:"status"`
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *CreateRoleRequest) (*model.Role, error) {
	// 检查角色名是否已存在
	existingRole, _ := s.roleRepo.GetByName(req.Name)
	if existingRole != nil {
		return nil, errors.New("角色名已存在")
	}

	status := req.Status
	if status == "" {
		status = "enabled"
	}

	role := &model.Role{
		Name:        req.Name,
		Description:  req.Description,
		Status:       status,
		IsSystem:     false,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	// 分配权限
	if len(req.PermissionIDs) > 0 {
		if err := s.roleRepo.AssignPermissions(role.ID, req.PermissionIDs); err != nil {
			return nil, fmt.Errorf("分配权限失败: %w", err)
		}
	}

	// 重新获取角色信息（包含权限）
	role, err := s.roleRepo.GetByID(role.ID)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(roleID uint, req *UpdateRoleRequest) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 系统角色不允许修改名称
	if role.IsSystem && req.Name != "" && req.Name != role.Name {
		return nil, errors.New("系统角色不允许修改名称")
	}

	// 更新字段
	if req.Name != "" {
		// 检查新名称是否已存在
		if req.Name != role.Name {
			existingRole, _ := s.roleRepo.GetByName(req.Name)
			if existingRole != nil {
				return nil, errors.New("角色名已存在")
			}
		}
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != "" {
		role.Status = req.Status
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, fmt.Errorf("更新角色失败: %w", err)
	}

	// 更新权限
	if req.PermissionIDs != nil {
		if err := s.roleRepo.AssignPermissions(roleID, req.PermissionIDs); err != nil {
			return nil, fmt.Errorf("更新权限失败: %w", err)
		}
	}

	// 重新获取角色信息
	role, err = s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(roleID uint) error {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 系统角色不允许删除
	if role.IsSystem {
		return errors.New("系统角色不允许删除")
	}

	return s.roleRepo.Delete(roleID)
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(roleID uint) (*model.Role, error) {
	return s.roleRepo.GetByID(roleID)
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(req *ListRoleRequest) ([]model.Role, int64, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.Name != "" {
		conditions["name LIKE ?"] = fmt.Sprintf("%%%s%%", req.Name)
	}
	if req.Status != "" {
		conditions["status = ?"] = req.Status
	}

	return s.roleRepo.List(req.Page, req.PageSize, conditions)
}



