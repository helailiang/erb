package repository

import (
	"erb/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").Preload("Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepository) List(page, pageSize int, conditions map[string]interface{}) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	
	// 应用筛选条件
	for key, value := range conditions {
		if key == "username LIKE ?" {
			query = query.Where("username LIKE ?", value)
		} else if key == "real_name LIKE ?" {
			query = query.Where("real_name LIKE ?", value)
		} else if key == "department = ?" {
			query = query.Where("department = ?", value)
		} else if key == "status = ?" {
			query = query.Where("status = ?", value)
		} else {
			query = query.Where(key, value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AssignRoles 分配角色
func (r *UserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	var roles []model.Role
	if err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}
	
	user, err := r.GetByID(userID)
	if err != nil {
		return err
	}

	return r.db.Model(user).Association("Roles").Replace(roles)
}

// GetUserPermissions 获取用户所有权限（包括角色权限和直接分配的权限）
func (r *UserRepository) GetUserPermissions(userID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	
	// 获取用户直接分配的权限
	var userPermissions []model.Permission
	if err := r.db.Model(&model.User{ID: userID}).Association("Permissions").Find(&userPermissions); err != nil {
		return nil, err
	}
	
	// 获取用户角色
	user, err := r.GetByID(userID)
	if err != nil {
		return nil, err
	}
	
	// 获取角色权限
	var rolePermissions []model.Permission
	for _, role := range user.Roles {
		var perms []model.Permission
		if err := r.db.Model(&role).Association("Permissions").Find(&perms); err != nil {
			continue
		}
		rolePermissions = append(rolePermissions, perms...)
	}
	
	// 合并权限并去重
	permissionMap := make(map[uint]model.Permission)
	for _, perm := range userPermissions {
		permissionMap[perm.ID] = perm
	}
	for _, perm := range rolePermissions {
		permissionMap[perm.ID] = perm
	}
	
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}
	
	return permissions, nil
}

