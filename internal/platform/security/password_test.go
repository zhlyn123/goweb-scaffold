package security

import "testing"

func TestPasswordHasher(t *testing.T) {
	hasher := NewPasswordHasherWithCost(4)

	hash, err := hasher.HashPassword("password123")
	if err != nil {
		t.Fatalf("生成密码哈希失败: %v", err)
	}

	if hash == "password123" {
		t.Fatalf("密码哈希值和明文密码相同")
	}

	if err := hasher.CheckPassword("password123", hash); err != nil {
		t.Fatalf("正确密码应该校验通过: %v", err)
	}

	if err := hasher.CheckPassword("wrong-password", hash); err != ErrInvalidPassword {
		t.Fatalf("错误密码应该返回 ErrInvalidPassword，实际返回: %v", err)
	}

}
