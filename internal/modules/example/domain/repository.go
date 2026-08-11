package domain

import (
	"context"
	"errors"
)

var ErrExampleItemNotFound = errors.New("示例数据不存在")

type repository interface {
	Create(ctx context.Context, item *ExampleItem) error
	FindByID(ctx context.Context, id string) (*ExampleItem, error)
}
