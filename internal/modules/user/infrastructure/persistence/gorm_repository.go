package persistence

import (
	"context"
	"errors"
	"fmt"
	"goweb-scaffold/internal/modules/user/domain"
	"goweb-scaffold/internal/platform/transaction"
	"time"

	"gorm.io/gorm"
)

type userModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	NickName     string    `gorm:"column:nickname"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at"`
}

func (userModel) TableName() string {
	return "users"
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *gormRepository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateUser(ctx context.Context, user *domain.User) error {
	db := transaction.DBFromContext(ctx, r.db)

	model := toModel(user)
	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

func (r *gormRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model userModel
	err := db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return toEntity(&model), nil
}

func (r *gormRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model userModel
	err := db.
		Where("email = ? AND deleted_at IS NULL", email).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("根据邮箱查询用户失败: %w", err)
	}

	return toEntity(&model), nil
}

func toModel(user *domain.User) userModel {
	return userModel{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		NickName:     user.NickName,
		Status:       string(user.Status),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		DeletedAt:    user.DeletedAt,
	}
}

func toEntity(model *userModel) *domain.User {
	return &domain.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		NickName:     model.NickName,
		Status:       domain.UserStatus(model.Status),
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		DeletedAt:    model.DeletedAt,
	}
}
