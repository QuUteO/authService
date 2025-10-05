package auth

import (
	models "Auth/internal/domain/model"
	jwtToken "Auth/internal/lib/jwt"
	"Auth/internal/lib/sl"
	"Auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"time"
)

type Auth struct {
	log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	AppProvider  AppProvider
	tokenTTl     time.Duration
}

// UserSaver интерфейс, сохраняющий user в бд
type UserSaver interface {
	UserSave(
		ctx context.Context,
		email string,
		PassHash []byte,
	) (uid int64, err error)
}

// UserProvider берет из бд данные о пользователе
type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

// AppProvider берет из бд данные о приложении, куда хочет зайти пользователь
type AppProvider interface {
	App(ctx context.Context, AppID int64) (app models.App, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserIsExists       = errors.New("user not found")
	ErrUserNotFound       = errors.New("user not found")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
		AppProvider:  appProvider,
		tokenTTl:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int64) (string, error) {

	const op = "Auth.service.auth.Login"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("email", email),
		slog.String("appID", strconv.FormatInt(appID, 10)),
	)

	log.Info("logging in")

	user, err := a.UserProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		log.Error("user not found", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid password", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.AppProvider.App(ctx, appID)
	if err != nil {
		log.Error("app not found", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged")

	token, err := jwtToken.NewToken(user, app, a.tokenTTl)
	if err != nil {
		log.Error("failed to create token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "service.auth.Register"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("email", email),
	)

	log.Info("Registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.UserSaver.UserSave(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserIsExists)
		}
		log.Error("Failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, UserID string) (bool, error) {
	const op = "service.auth.IsAdmin"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("email", UserID),
	)

	log.Info("Checking user")

	isAdmin, err := a.IsAdmin(ctx, UserID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return isAdmin, nil
}
