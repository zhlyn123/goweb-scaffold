package transaction

import (
	"context"

	"gorm.io/gorm"
)

type Manager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type gormTxKey struct{}

type GormManager struct {
	db *gorm.DB
}

func NewGormManager(db *gorm.DB) *GormManager {
	return &GormManager{db: db}
}

// Run 执行一个事务。
// fn 返回 nil 时提交事务；fn 返回 error 时回滚事务。
func (m *GormManager) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, gormTxKey{}, tx)
		return fn(txCtx)
	})
}

// DBFromContext 返回当前上下文中的事务连接。
// 如果当前上下文没有事务，则返回普通数据库连接。
// 这个函数只给 repository 实现层使用，application 和 domain 层不要直接使用。
func DBFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback.WithContext(ctx)
}
