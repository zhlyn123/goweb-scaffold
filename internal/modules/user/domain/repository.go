package domain

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrUserEmailExists = errors.New("邮箱已存在")
	ErrUserDisabled    = errors.New("用户已被禁用")
	ErrInvalidPassword = errors.New("邮箱或密码错误")
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
