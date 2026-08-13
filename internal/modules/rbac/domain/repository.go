package domain

import (
	"context"
	"errors"
)

var (
	ErrRoleNotFound       = errors.New("角色不存在")
	ErrPermissionNotFound = errors.New("权限不存在")
	ErrRoleAlreadyBound   = errors.New("用户已绑定该角色")
	ErrPermissionDenied   = errors.New("权限不足")
)

// Repository 定义 RBAC 模块需要的持久化能力。
// 具体实现放在 infrastructure 层，domain 层不关心数据库或 ORM。
type Repository interface {
	// CreateRole 创建角色。
	CreateRole(ctx context.Context, role *Role) error

	// CreatePermission 创建权限。
	CreatePermission(ctx context.Context, permission *Permission) error

	// GetRoleByID 根据角色 ID 查询角色。
	GetRoleByID(ctx context.Context, id string) (*Role, error)

	// GetRoleByCode 根据角色编码查询角色。
	GetRoleByCode(ctx context.Context, code string) (*Role, error)

	// GetPermissionByID 根据权限 ID 查询权限。
	GetPermissionByID(ctx context.Context, id string) (*Permission, error)

	// GetPermissionByCode 根据权限编码查询权限。
	GetPermissionByCode(ctx context.Context, code string) (*Permission, error)

	// BindRoleToUser 将角色绑定给用户。
	BindRoleToUser(ctx context.Context, userID string, roleID string) error

	// BindPermissionToRole 将权限绑定给角色。
	BindPermissionToRole(ctx context.Context, roleID string, permissionID string) error

	// ListRolesByUserID 查询用户拥有的角色列表。
	ListRolesByUserID(ctx context.Context, userID string) ([]*Role, error)

	// ListPermissionsByUserID 查询用户通过角色获得的权限列表。
	ListPermissionsByUserID(ctx context.Context, userID string) ([]*Permission, error)

	// HasPermission 判断用户是否拥有指定权限。
	HasPermission(ctx context.Context, userID string, permissionCode string) (bool, error)
}
