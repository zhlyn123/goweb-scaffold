package transaction

import (
	"context"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDBFromContextWithoutTransaction(t *testing.T) {
	ctx := context.Background()

	db, err := gorm.Open(
		postgres.Open("host=localhost port=5432 user=postgres password=postgres dbname=goweb_scaffold sslmode=disable"),
		&gorm.Config{
			// 单元测试只验证 GORM 对象行为，不真的连接数据库。
			DisableAutomaticPing: true,
		},
	)
	if err != nil {
		t.Fatalf("创建 GORM 对象失败: %v", err)
	}

	got := DBFromContext(ctx, db)
	if got == nil {
		t.Fatal("DBFromContext() = nil, want *gorm.DB")
	}

	if got.Statement.Context != ctx {
		t.Fatalf("DBFromContext() context = %v, want %v", got.Statement.Context, ctx)
	}
}

func TestDBFromContextWithTransaction(t *testing.T) {
	fallback := &gorm.DB{}
	tx := &gorm.DB{}

	ctx := context.WithValue(context.Background(), gormTxKey{}, tx)

	got := DBFromContext(ctx, fallback)
	if got != tx {
		t.Fatalf("DBFromContext() = %p, want transaction %p", got, tx)
	}
}