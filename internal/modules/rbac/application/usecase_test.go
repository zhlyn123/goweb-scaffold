package application

import (
	"context"
	"errors"
	"goweb-scaffold/internal/modules/rbac/domain"
	"testing"
)

type fakeIDGenerator struct{}

func (g fakeIDGenerator) NewID() string {
	return "generated-id"
}

type fakeRepository struct {
	hasPermission bool
}

func (r *fakeRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	return nil
}

func (r *fakeRepository) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	return nil
}

func (r *fakeRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	return nil, domain.ErrRoleNotFound
}

func (r *fakeRepository) GetRoleByCode(ctx context.Context, code string) (*domain.Role, error) {
	return &domain.Role{ID: "role-1", Code: code}, nil
}

func (r *fakeRepository) GetPermissionByID(ctx context.Context, id string) (*domain.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}

func (r *fakeRepository) GetPermissionByCode(ctx context.Context, code string) (*domain.Permission, error) {
	return &domain.Permission{ID: "permission-1", Code: code}, nil
}

func (r *fakeRepository) BindRoleToUser(ctx context.Context, userID string, roleID string) error {
	return nil
}

func (r *fakeRepository) BindPermissionToRole(ctx context.Context, roleID string, permissionID string) error {
	return nil
}

func (r *fakeRepository) ListRolesByUserID(ctx context.Context, userID string) ([]*domain.Role, error) {
	return nil, nil
}

func (r *fakeRepository) ListPermissionsByUserID(ctx context.Context, userID string) ([]*domain.Permission, error) {
	return nil, nil
}

func (r *fakeRepository) HasPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	return r.hasPermission, nil
}

func TestCheckPermissionAllowsUserWithPermission(t *testing.T) {
	usecase := NewUsecase(&fakeRepository{hasPermission: true}, fakeIDGenerator{})

	err := usecase.CheckPermission(context.Background(), CheckPermissionCommand{
		UserID:         "user-1",
		PermissionCode: "user:read",
	})
	if err != nil {
		t.Fatalf("CheckPermission error = %v, want nil", err)
	}
}

func TestCheckPermissionDeniesUserWithoutPermission(t *testing.T) {
	usecase := NewUsecase(&fakeRepository{hasPermission: false}, fakeIDGenerator{})

	err := usecase.CheckPermission(context.Background(), CheckPermissionCommand{
		UserID:         "user-1",
		PermissionCode: "user:read",
	})
	if !errors.Is(err, domain.ErrPermissionDenied) {
		t.Fatalf("CheckPermission error = %v, want %v", err, domain.ErrPermissionDenied)
	}
}
