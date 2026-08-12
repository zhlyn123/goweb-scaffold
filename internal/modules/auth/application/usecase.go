package application

import (
	"context"
	"errors"
	"fmt"
	userdomain "goweb-scaffold/internal/modules/user/domain"
	platformsecurity "goweb-scaffold/internal/platform/security"
	"time"
)

type PasswordHasher interface {
	HashPassword(plainPassword string) (string, error)
	CheckPassword(plainPassword string, passwrodHash string) error
}

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, userID string) (string, time.Time, error)
}

type IDGenerator interface {
	NewID() string
}

type Usecase struct {
	userRepo     userdomain.Repository
	passwordHash PasswordHasher
	tokenIssuer  TokenIssuer
	idGenerator  IDGenerator
}

func NewUsecase(
	userRepo userdomain.Repository,
	passwordHash PasswordHasher,
	tokenIssuer TokenIssuer,
	idGenerator IDGenerator,
) *Usecase {
	return &Usecase{
		userRepo:     userRepo,
		passwordHash: passwordHash,
		tokenIssuer:  tokenIssuer,
		idGenerator:  idGenerator,
	}
}

type RegisterCommand struct {
	Email    string
	Password string
	NickName string
}

type RegisterResult struct {
	UserID string
	Email  string
}

func (u *Usecase) Register(ctx context.Context, cmd RegisterCommand) (*RegisterResult, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, userdomain.ErrUserNotFound) {
		return nil, fmt.Errorf("检查邮箱是否已注册失败: %w", err)
	}

	if existingUser != nil {
		return nil, userdomain.ErrUserEmailExists
	}

	passwordHash, err := u.passwordHash.HashPassword(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("加密用户密码失败: %w", err)
	}

	now := time.Now()

	user := &userdomain.User{
		ID:           u.idGenerator.NewID(),
		Email:        cmd.Email,
		PasswordHash: passwordHash,
		NickName:     cmd.NickName,
		Status:       userdomain.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &RegisterResult{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
	ExpiresAt   time.Time
	TokenType   string
}

func (u *Usecase) Login(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, cmd.Email)
	if errors.Is(err, userdomain.ErrUserNotFound) {
		return nil, userdomain.ErrInvalidPassword
	}

	if err != nil {
		return nil, fmt.Errorf("查询登录用户失败: %w", err)
	}

	if user.Status != userdomain.UserStatusActive {
		return nil, userdomain.ErrUserDisabled
	}

	if err := u.passwordHash.CheckPassword(cmd.Password, user.PasswordHash); err != nil {
		if errors.Is(err, platformsecurity.ErrInvalidPassword) {
			return nil, userdomain.ErrInvalidPassword
		}

		return nil, fmt.Errorf("校验用户密码失败: %w", err)
	}

	accessToken, expiresAt, err := u.tokenIssuer.IssueAccessToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("签发访问令牌失败: %w", err)
	}

	return &LoginResult{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}, nil
}
