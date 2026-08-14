package application

import (
	"context"
	"fmt"
	"goweb-scaffold/internal/modules/rbac/domain"
	"time"
)

// IDGenerator 定义 ID 生成能力。
// 这里依赖接口而不是具体 UUID 实现，方便测试时替换。
type IDGenerator interface {
	NewID() string
}

// Usecase 封装 RBAC 模块的应用层用例。
type Usecase struct {
	rbacRepo    domain.Repository
	idGenerator IDGenerator
}

// NewUsecase 创建 RBAC 应用服务。
func NewUsecase(
	rbacRepo domain.Repository,
	idGenerator IDGenerator,
) *Usecase {
	return &Usecase{
		rbacRepo:    rbacRepo,
		idGenerator: idGenerator,
	}
}

// CreateRoleCommand 表示创建角色的输入参数。
type CreateRoleCommand struct {
	Code        string
	Name        string
	Description string
}

// CreatePermissionCommand 表示创建权限的输入参数。
type CreatePermissionCommand struct {
	Code        string
	Name        string
	Description string
}

// BindRoleToUserCommand 表示给用户绑定角色的输入参数。
type BindRoleToUserCommand struct {
	UserID   string
	RoleCode string
}

// GrantPermissionToRoleCommand 表示给角色绑定权限的输入参数。
type GrantPermissionToRoleCommand struct {
	RoleCode       string
	PermissionCode string
}

// CheckPermissionCommand 表示检查用户权限的输入参数。
type CheckPermissionCommand struct {
	UserID         string
	PermissionCode string
}

// CreateRole 创建角色。
func (u *Usecase) CreateRole(ctx context.Context, cmd CreateRoleCommand) (*domain.Role, error) {
	now := time.Now()

	role := &domain.Role{
		ID:          u.idGenerator.NewID(),
		Code:        cmd.Code,
		Name:        cmd.Name,
		Description: cmd.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.rbacRepo.CreateRole(ctx, role); err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	return role, nil
}

// CreatePermission 创建权限。
func (u *Usecase) CreatePermission(ctx context.Context, cmd CreatePermissionCommand) (*domain.Permission, error) {
	now := time.Now()

	permission := &domain.Permission{
		ID:          u.idGenerator.NewID(),
		Code:        cmd.Code,
		Name:        cmd.Name,
		Description: cmd.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.rbacRepo.CreatePermission(ctx, permission); err != nil {
		return nil, fmt.Errorf("创建权限失败: %w", err)
	}

	return permission, nil
}

// BindRoleToUser 根据角色编码给用户绑定角色。
// 外部调用方不需要知道 role_id，只需要传稳定的 role code。
func (u *Usecase) BindRoleToUser(ctx context.Context, cmd BindRoleToUserCommand) error {
	role, err := u.rbacRepo.GetRoleByCode(ctx, cmd.RoleCode)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}

	if err := u.rbacRepo.BindRoleToUser(ctx, cmd.UserID, role.ID); err != nil {
		return fmt.Errorf("绑定用户角色失败: %w", err)
	}

	return nil
}

// GrantPermissionToRole 根据角色编码和权限编码建立授权关系。
// 外部调用方不需要知道数据库 ID，只需要使用稳定 code。
func (u *Usecase) GrantPermissionToRole(ctx context.Context, cmd GrantPermissionToRoleCommand) error {
	role, err := u.rbacRepo.GetRoleByCode(ctx, cmd.RoleCode)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}

	permission, err := u.rbacRepo.GetPermissionByCode(ctx, cmd.PermissionCode)
	if err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}

	if err := u.rbacRepo.BindPermissionToRole(ctx, role.ID, permission.ID); err != nil {
		return fmt.Errorf("绑定角色权限失败: %w", err)
	}

	return nil
}

// ListUserPermissions 查询用户最终拥有的权限列表。
func (u *Usecase) ListUserPermissions(ctx context.Context, userID string) ([]*domain.Permission, error) {
	permissions, err := u.rbacRepo.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户权限列表失败: %w", err)
	}

	return permissions, nil
}

// CheckPermission 检查用户是否拥有指定权限。
// 如果没有权限，返回 domain.ErrPermissionDenied。
func (u *Usecase) CheckPermission(ctx context.Context, cmd CheckPermissionCommand) error {
	ok, err := u.rbacRepo.HasPermission(ctx, cmd.UserID, cmd.PermissionCode)
	if err != nil {
		return fmt.Errorf("检查用户权限失败: %w", err)
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}