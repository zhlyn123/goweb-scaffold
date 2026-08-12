package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("密码不正确")

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

func NewPasswordHasherWithCost(cost int) *PasswordHasher {
	return &PasswordHasher{
		cost: cost,
	}
}

// HashPassword 将明文密码转换成 bcrypt 哈希值。
// 注册用户时只能保存这个哈希值，不能保存明文密码。
func (h *PasswordHasher) HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), h.cost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}
	return string(hash), nil
}

// CheckPassword 校验明文密码和数据库里的密码哈希是否匹配。
// 登录时使用这个方法判断用户输入的密码是否正确。
func (h *PasswordHasher) CheckPassword(plainPassword string, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidPassword
	}

	if err != nil {
		return fmt.Errorf("校验密码失败: %w", err)
	}

	return nil
}
