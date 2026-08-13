package persistence

import (
	"context"
	"errors"
	"fmt"
	"goweb-scaffold/internal/modules/rbac/domain"
	"goweb-scaffold/internal/platform/transaction"
	"time"

	"gorm.io/gorm"
)

// roleModel 对应 roles 表。
type roleModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at"`
}

func (roleModel) TableName() string {
	return "roles"
}

// permissionModel 对应 permissions 表。
type permissionModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at"`
}

func (permissionModel) TableName() string {
	return "permissions"
}

// userRoleModel 对应 user_roles 表。
type userRoleModel struct {
	UserID    string    `gorm:"column:user_id;primaryKey"`
	RoleID    string    `gorm:"column:role_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (userRoleModel) TableName() string {
	return "user_roles"
}

// rolePermissionModel 对应 role_permissions 表。
type rolePermissionModel struct {
	RoleID       string    `gorm:"column:role_id;primaryKey"`
	PermissionID string    `gorm:"column:permission_id;primaryKey"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (rolePermissionModel) TableName() string {
	return "role_permissions"
}

// gormRepository 使用 GORM 实现 RBAC 仓储接口。
type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository 创建 RBAC 的 GORM 仓储实现。
func NewGormRepository(db *gorm.DB) *gormRepository {
	return &gormRepository{db: db}
}

// CreateRole 创建角色。
func (r *gormRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	db := transaction.DBFromContext(ctx, r.db)

	model := toRoleModel(role)
	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}

	return nil
}

// CreatePermission 创建权限。
func (r *gormRepository) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	db := transaction.DBFromContext(ctx, r.db)

	model := toPermissionModel(permission)
	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("创建权限失败: %w", err)
	}

	return nil
}

// GetRoleByID 根据角色 ID 查询角色。
func (r *gormRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model roleModel
	err := db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRoleNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("根据 ID 查询角色失败: %w", err)
	}

	return toRoleEntity(&model), nil
}

// GetRoleByCode 根据角色编码查询角色。
func (r *gormRepository) GetRoleByCode(ctx context.Context, code string) (*domain.Role, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model roleModel
	err := db.
		Where("code = ? AND deleted_at IS NULL", code).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRoleNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("根据编码查询角色失败: %w", err)
	}

	return toRoleEntity(&model), nil
}

// GetPermissionByID 根据权限 ID 查询权限。
func (r *gormRepository) GetPermissionByID(ctx context.Context, id string) (*domain.Permission, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model permissionModel
	err := db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPermissionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("根据 ID 查询权限失败: %w", err)
	}

	return toPermissionEntity(&model), nil
}

// GetPermissionByCode 根据权限编码查询权限。
func (r *gormRepository) GetPermissionByCode(ctx context.Context, code string) (*domain.Permission, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model permissionModel
	err := db.
		Where("code = ? AND deleted_at IS NULL", code).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPermissionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("根据编码查询权限失败: %w", err)
	}

	return toPermissionEntity(&model), nil
}

// BindRoleToUser 将角色绑定给用户。
func (r *gormRepository) BindRoleToUser(ctx context.Context, userID string, roleID string) error {
	db := transaction.DBFromContext(ctx, r.db)

	model := userRoleModel{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("绑定用户角色失败: %w", err)
	}

	return nil
}

// BindPermissionToRole 将权限绑定给角色。
func (r *gormRepository) BindPermissionToRole(ctx context.Context, roleID string, permissionID string) error {
	db := transaction.DBFromContext(ctx, r.db)

	model := rolePermissionModel{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
	}

	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("绑定角色权限失败: %w", err)
	}

	return nil
}

// ListRolesByUserID 查询用户拥有的角色列表。
func (r *gormRepository) ListRolesByUserID(ctx context.Context, userID string) ([]*domain.Role, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var models []roleModel
	err := db.
		Table("roles").
		Select("roles.*").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.deleted_at IS NULL", userID).
		Find(&models).
		Error
	if err != nil {
		return nil, fmt.Errorf("查询用户角色列表失败: %w", err)
	}

	roles := make([]*domain.Role, 0, len(models))
	for i := range models {
		roles = append(roles, toRoleEntity(&models[i]))
	}

	return roles, nil
}

// ListPermissionsByUserID 查询用户通过角色获得的权限列表。
func (r *gormRepository) ListPermissionsByUserID(ctx context.Context, userID string) ([]*domain.Permission, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var models []permissionModel
	err := db.
		Table("permissions").
		Select("DISTINCT permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL", userID).
		Find(&models).
		Error
	if err != nil {
		return nil, fmt.Errorf("查询用户权限列表失败: %w", err)
	}

	permissions := make([]*domain.Permission, 0, len(models))
	for i := range models {
		permissions = append(permissions, toPermissionEntity(&models[i]))
	}

	return permissions, nil
}

// HasPermission 判断用户是否拥有指定权限。
func (r *gormRepository) HasPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var count int64
	err := db.
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permissions.code = ? AND permissions.deleted_at IS NULL", permissionCode).
		Count(&count).
		Error
	if err != nil {
		return false, fmt.Errorf("检查用户权限失败: %w", err)
	}

	return count > 0, nil
}

// toRoleModel 将领域角色转换为数据库模型。
func toRoleModel(role *domain.Role) roleModel {
	return roleModel{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		DeletedAt:   role.DeletedAt,
	}
}

// toRoleEntity 将数据库角色模型转换为领域实体。
func toRoleEntity(model *roleModel) *domain.Role {
	return &domain.Role{
		ID:          model.ID,
		Code:        model.Code,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   model.DeletedAt,
	}
}

// toPermissionModel 将领域权限转换为数据库模型。
func toPermissionModel(permission *domain.Permission) permissionModel {
	return permissionModel{
		ID:          permission.ID,
		Code:        permission.Code,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
		DeletedAt:   permission.DeletedAt,
	}
}

// toPermissionEntity 将数据库权限模型转换为领域实体。
func toPermissionEntity(model *permissionModel) *domain.Permission {
	return &domain.Permission{
		ID:          model.ID,
		Code:        model.Code,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   model.DeletedAt,
	}
}
