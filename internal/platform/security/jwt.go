package security

import (
	"context"
	"errors"
	"fmt"
	"goweb-scaffold/internal/platform/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenMissingUserID = errors.New("token 缺少用户 ID")
	ErrTokenInvalid       = errors.New("token 无效")
)

type AccessTokenCliaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret         []byte
	issuer         string
	accessTokenTTL time.Duration
}

func NewTokenService(cfg *config.JWTConfig) *TokenService {
	return &TokenService{
		secret:         []byte(cfg.Secret),
		issuer:         cfg.Issuer,
		accessTokenTTL: cfg.AccessTokenTTL,
	}
}

func (s *TokenService) IssueAccessToken(ctx context.Context, userID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTokenTTL)

	claims := AccessTokenCliaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("签名 access token 失败: %w", err)
	}

	return signedToken, expiresAt, nil
}
