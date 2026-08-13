package domain

import "time"

// Role 表示系统角色，例如 admin、user。
type Role struct {
	ID          string
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

// Permission 表示系统权限点，例如 user:read、user:manage。
type Permission struct {
	ID          string
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

// UserRole 表示用户和角色的绑定关系。
type UserRole struct {
	UserID    string
	RoleID    string
	CreatedAt time.Time
}

// RolePermission 表示角色和权限的绑定关系。
type RolePermission struct {
	RoleID       string
	PermissionID string
	CreatedAt    time.Time
}