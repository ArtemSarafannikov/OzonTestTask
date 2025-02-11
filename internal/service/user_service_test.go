package service

import (
	"context"
	"errors"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/testutils"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestUserService_Register_Success(t *testing.T) {
	// Патчим функции из utils.
	oldHashPassword := utils.HashPassword
	oldGenerateJWT := utils.GenerateJWT
	defer func() {
		utils.HashPassword = oldHashPassword
		utils.GenerateJWT = oldGenerateJWT
	}()

	// HashPassword возвращает детерминированный результат.
	utils.HashPassword = func(password string) (string, error) {
		return "hashed_" + password, nil
	}
	// GenerateJWT возвращает токен в формате "token_<userID>"
	utils.GenerateJWT = func(userID string) (string, error) {
		return "token_" + userID, nil
	}

	// Создаём FakeRepository, который будет присваивать ID новому пользователю.
	fakeRepo := &testutils.FakeRepository{
		CreateUserFn: func(ctx context.Context, user *models.User) (*models.User, error) {
			user.ID = "u1"
			// Можно установить и дату создания, если нужно.
			user.CreatedAt = time.Now()
			return user, nil
		},
	}

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)

	token, user, err := svc.Register(context.Background(), "testuser", "password123")
	require.NoError(t, err)
	require.NotNil(t, user)
	// Проверяем, что пароль захешировался правильно.
	require.Equal(t, "hashed_password123", user.Password)
	require.Equal(t, "testuser", user.Username)
	// Токен генерируется на основе ID, который FakeRepository установил как "u1"
	require.Equal(t, "token_u1", token)
}

// TestUserService_Register_HashError проверяет, что при ошибке хеширования возвращается ошибка.
func TestUserService_Register_HashError(t *testing.T) {
	oldHashPassword := utils.HashPassword
	defer func() { utils.HashPassword = oldHashPassword }()

	// Симулируем ошибку хеширования.
	utils.HashPassword = func(password string) (string, error) {
		return "", errors.New("hash error")
	}

	fakeRepo := &testutils.FakeRepository{}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)

	token, _, err := svc.Register(context.Background(), "testuser", "password123")
	require.Error(t, err)
	require.EqualError(t, err, "hash error")
	// Токен пустой, а user может быть не nil (в зависимости от реализации), но нам важно, что произошла ошибка.
	require.Empty(t, token)
}

// TestUserService_Register_CreateUserError проверяет ситуацию, когда репозиторий не может создать пользователя.
func TestUserService_Register_CreateUserError(t *testing.T) {
	oldHashPassword := utils.HashPassword
	oldGenerateJWT := utils.GenerateJWT
	defer func() {
		utils.HashPassword = oldHashPassword
		utils.GenerateJWT = oldGenerateJWT
	}()

	utils.HashPassword = func(password string) (string, error) {
		return "hashed_" + password, nil
	}
	// Переопределим GenerateJWT, хотя он не будет вызван, если создание пользователя завершится ошибкой.
	utils.GenerateJWT = func(userID string) (string, error) {
		return "token_" + userID, nil
	}

	fakeRepo := &testutils.FakeRepository{
		CreateUserFn: func(ctx context.Context, user *models.User) (*models.User, error) {
			return nil, errors.New("create user failed")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)

	token, user, err := svc.Register(context.Background(), "testuser", "password123")
	require.Error(t, err)
	require.EqualError(t, err, "create user failed")
	require.Empty(t, token)
	require.Nil(t, user)
}

// TestUserService_Login_Success проверяет успешную аутентификацию пользователя.
func TestUserService_Login_Success(t *testing.T) {
	oldCheckPasswordHash := utils.CheckPasswordHash
	oldGenerateJWT := utils.GenerateJWT
	defer func() {
		utils.CheckPasswordHash = oldCheckPasswordHash
		utils.GenerateJWT = oldGenerateJWT
	}()

	// CheckPasswordHash возвращает true, если hash соответствует "hashed_" + password.
	utils.CheckPasswordHash = func(password, hash string) bool {
		return hash == "hashed_"+password
	}
	utils.GenerateJWT = func(userID string) (string, error) {
		return "token_" + userID, nil
	}

	fakeRepo := &testutils.FakeRepository{
		GetUserByLoginFn: func(ctx context.Context, login string) (*models.User, error) {
			// Симулируем, что найден пользователь с хешированным паролем.
			return &models.User{
				ID:       "u1",
				Username: login,
				Password: "hashed_password123",
			}, nil
		},
		// Для фиксации времени активности (асинхронно) можно просто оставить пустую реализацию.
		FixLastActivityFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)
	token, user, err := svc.Login(context.Background(), "testuser", "password123")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "u1", user.ID)
	require.Equal(t, "token_u1", token)
}

// TestUserService_Login_InvalidCredentials_WrongPassword проверяет, что неверный пароль приводит к ошибке.
func TestUserService_Login_InvalidCredentials_WrongPassword(t *testing.T) {
	oldCheckPasswordHash := utils.CheckPasswordHash
	defer func() { utils.CheckPasswordHash = oldCheckPasswordHash }()

	// Для любой комбинации вернём false.
	utils.CheckPasswordHash = func(password, hash string) bool {
		return false
	}

	fakeRepo := &testutils.FakeRepository{
		GetUserByLoginFn: func(ctx context.Context, login string) (*models.User, error) {
			return &models.User{
				ID:       "u1",
				Username: login,
				Password: "hashed_password123",
			}, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)
	token, user, err := svc.Login(context.Background(), "testuser", "wrongpassword")
	require.Error(t, err)
	require.Equal(t, cstErrors.InvalidCredentials, err)
	require.Empty(t, token)
	require.Nil(t, user)
}

// TestUserService_Login_GetUserByLoginError проверяет, что если репозиторий не может найти пользователя, возвращается ошибка.
func TestUserService_Login_GetUserByLoginError(t *testing.T) {
	fakeRepo := &testutils.FakeRepository{
		GetUserByLoginFn: func(ctx context.Context, login string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)
	token, user, err := svc.Login(context.Background(), "testuser", "password123")
	// При ошибке поиска пользователь считается недействительным.
	require.Error(t, err)
	require.Equal(t, cstErrors.InvalidCredentials, err)
	require.Empty(t, token)
	require.Nil(t, user)
}

// TestUserService_GetUserByID проверяет получение пользователя по ID.
func TestUserService_GetUserByID(t *testing.T) {
	fakeRepo := &testutils.FakeRepository{
		GetUserByIDFn: func(ctx context.Context, id string) (*models.User, error) {
			if id == "u1" {
				return &models.User{
					ID:       "u1",
					Username: "testuser",
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewUserService(log, fakeRepo)
	user, err := svc.GetUserByID(context.Background(), "u1")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "u1", user.ID)
}
