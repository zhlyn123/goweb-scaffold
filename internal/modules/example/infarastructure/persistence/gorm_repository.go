package persistence

import (
	"context"
	"errors"
	"fmt"
	"goweb-scaffold/internal/modules/example/domain"
	"goweb-scaffold/internal/platform/transaction"
	"time"

	"gorm.io/gorm"
)

type exampleItemModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at"`
}

func (exampleItemModel) TableName() string {
	return "example_items"
}

type gormRepository struct {
	db *gorm.DB
}

func newGormRepository(db *gorm.DB) *gormRepository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, item *domain.ExampleItem) error {
	db := transaction.DBFromContext(ctx, r.db)

	model := toModel(item)
	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("创建示例数据失败: %w", err)
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id string) (*domain.ExampleItem, error) {
	db := transaction.DBFromContext(ctx, r.db)

	var model exampleItemModel
	err := db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrExampleItemNotFound
	}

	if err != nil {
		return nil, err
	}

	return toEntity(&model), nil
}

func toModel(item *domain.ExampleItem) exampleItemModel {
	return exampleItemModel{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreateAt,
		UpdatedAt:   item.UpdateAt,
		DeletedAt:   item.DeletedAt,
	}
}

func toEntity(model *exampleItemModel) *domain.ExampleItem {
	return &domain.ExampleItem{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreateAt:    model.CreatedAt,
		UpdateAt:    model.UpdatedAt,
		DeletedAt:   model.DeletedAt,
	}
}
