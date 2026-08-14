package application

import (
	"context"
	"errors"
	"testing"
	"time"

	userdomain "goweb-scaffold/internal/modules/user/domain"
	platformsecurity "goweb-scaffold/internal/platform/security"
)

type memoryUserRepository struct {
	usersByEmail map[string]*userdomain.User
	usersByID    map[string]*userdomain.User
	createErr    error
}

func newMemoryUserRepository() *memoryUserRepository {
	return &memoryUserRepository{
		usersByEmail: map[string]*userdomain.User{},
		usersByID:    map[string]*userdomain.User{},
	}
}

func (r *memoryUserRepository) CreateUser(_ context.Context, user *userdomain.User) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.usersByEmail[user.Email] = user
	r.usersByID[user.ID] = user
	return nil
}

func (r *memoryUserRepository) GetUserByID(_ context.Context, id string) (*userdomain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, userdomain.ErrUserNotFound
	}
	return user, nil
}

func (r *memoryUserRepository) GetUserByEmail(_ context.Context, email string) (*userdomain.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, userdomain.ErrUserNotFound
	}
	return user, nil
}

type testPasswordHasher struct {
	hashErr  error
	checkErr error
}

func (h testPasswordHasher) HashPassword(plainPassword string) (string, error) {
	if h.hashErr != nil {
		return "", h.hashErr
	}
	return "hash:" + plainPassword, nil
}

func (h testPasswordHasher) CheckPassword(plainPassword string, passwordHash string) error {
	if h.checkErr != nil {
		return h.checkErr
	}
	if passwordHash != "hash:"+plainPassword {
		return platformsecurity.ErrInvalidPassword
	}
	return nil
}

type testTokenIssuer struct {
	token string
	err   error
}

func (i testTokenIssuer) IssueAccessToken(context.Context, string) (string, time.Time, error) {
	if i.err != nil {
		return "", time.Time{}, i.err
	}
	return i.token, time.Unix(3600, 0), nil
}

type testIDGenerator struct {
	id string
}

func (g testIDGenerator) NewID() string {
	return g.id
}

type testTransactionManager struct {
	called bool
}

func (m *testTransactionManager) Run(ctx context.Context, fn func(context.Context) error) error {
	m.called = true
	return fn(ctx)
}

type testRoleBinder struct {
	roleCode string
	err      error
}

func (b *testRoleBinder) BindRoleToUser(_ context.Context, _ string, roleCode string) error {
	b.roleCode = roleCode
	return b.err
}

func TestRegisterCreatesActiveUserAndBindsDefaultRole(t *testing.T) {
	repo := newMemoryUserRepository()
	tx := &testTransactionManager{}
	binder := &testRoleBinder{}
	usecase := NewUsecase(
		repo,
		testPasswordHasher{},
		testTokenIssuer{token: "token"},
		testIDGenerator{id: "user-1"},
		tx,
		binder,
		"user",
	)

	result, err := usecase.Register(context.Background(), RegisterCommand{
		Email:    "user@example.com",
		Password: "password123",
		NickName: "Demo",
	})
	if err != nil {
		t.Fatalf("Register error = %v, want nil", err)
	}
	if result.UserID != "user-1" || result.Email != "user@example.com" {
		t.Fatalf("Register result = %#v", result)
	}
	if !tx.called {
		t.Fatal("expected transaction manager to be used")
	}
	if binder.roleCode != "user" {
		t.Fatalf("default role = %q, want user", binder.roleCode)
	}
	created := repo.usersByID["user-1"]
	if created == nil {
		t.Fatal("expected user to be persisted")
	}
	if created.PasswordHash == "password123" {
		t.Fatal("expected stored password to be hashed")
	}
	if created.Status != userdomain.UserStatusActive {
		t.Fatalf("status = %s, want %s", created.Status, userdomain.UserStatusActive)
	}
}

func TestRegisterRejectsExistingEmail(t *testing.T) {
	repo := newMemoryUserRepository()
	repo.usersByEmail["user@example.com"] = &userdomain.User{ID: "existing", Email: "user@example.com"}
	usecase := NewUsecase(repo, testPasswordHasher{}, testTokenIssuer{}, testIDGenerator{id: "user-1"}, nil, nil, "")

	_, err := usecase.Register(context.Background(), RegisterCommand{
		Email:    "user@example.com",
		Password: "password123",
	})
	if !errors.Is(err, userdomain.ErrUserEmailExists) {
		t.Fatalf("Register error = %v, want %v", err, userdomain.ErrUserEmailExists)
	}
}

func TestRegisterReturnsRoleBindingError(t *testing.T) {
	repo := newMemoryUserRepository()
	binderErr := errors.New("bind failed")
	usecase := NewUsecase(
		repo,
		testPasswordHasher{},
		testTokenIssuer{},
		testIDGenerator{id: "user-1"},
		nil,
		&testRoleBinder{err: binderErr},
		"user",
	)

	_, err := usecase.Register(context.Background(), RegisterCommand{
		Email:    "user@example.com",
		Password: "password123",
	})
	if !errors.Is(err, binderErr) {
		t.Fatalf("Register error = %v, want wrapped %v", err, binderErr)
	}
}

func TestLoginIssuesAccessToken(t *testing.T) {
	repo := newMemoryUserRepository()
	repo.usersByEmail["user@example.com"] = &userdomain.User{
		ID:           "user-1",
		Email:        "user@example.com",
		PasswordHash: "hash:password123",
		Status:       userdomain.UserStatusActive,
	}
	usecase := NewUsecase(repo, testPasswordHasher{}, testTokenIssuer{token: "access-token"}, testIDGenerator{}, nil, nil, "")

	result, err := usecase.Login(context.Background(), LoginCommand{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login error = %v, want nil", err)
	}
	if result.AccessToken != "access-token" || result.TokenType != "Bearer" {
		t.Fatalf("Login result = %#v", result)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	repo := newMemoryUserRepository()
	repo.usersByEmail["user@example.com"] = &userdomain.User{
		ID:           "user-1",
		Email:        "user@example.com",
		PasswordHash: "hash:password123",
		Status:       userdomain.UserStatusActive,
	}
	usecase := NewUsecase(repo, testPasswordHasher{}, testTokenIssuer{}, testIDGenerator{}, nil, nil, "")

	_, err := usecase.Login(context.Background(), LoginCommand{
		Email:    "user@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, userdomain.ErrInvalidPassword) {
		t.Fatalf("Login error = %v, want %v", err, userdomain.ErrInvalidPassword)
	}
}

func TestLoginRejectsDisabledUser(t *testing.T) {
	repo := newMemoryUserRepository()
	repo.usersByEmail["user@example.com"] = &userdomain.User{
		ID:           "user-1",
		Email:        "user@example.com",
		PasswordHash: "hash:password123",
		Status:       userdomain.UserStatusDisabled,
	}
	usecase := NewUsecase(repo, testPasswordHasher{}, testTokenIssuer{}, testIDGenerator{}, nil, nil, "")

	_, err := usecase.Login(context.Background(), LoginCommand{
		Email:    "user@example.com",
		Password: "password123",
	})
	if !errors.Is(err, userdomain.ErrUserDisabled) {
		t.Fatalf("Login error = %v, want %v", err, userdomain.ErrUserDisabled)
	}
}
