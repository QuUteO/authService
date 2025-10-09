package sqlite

import (
	models "Auth/internal/domain/model"
	"Auth/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, status.Error(codes.Internal, op)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) UserSave(ctx context.Context, email string, passHash []byte,
) (uid int64, err error) {
	const op = "service.userSave"

	smtp, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := smtp.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr *sqlite3.Error

		if !errors.As(err, &sqliteErr) && status.Code(err) == codes.AlreadyExists {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// получаем id пользователя
	uid, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}

func (s *Storage) User(ctx context.Context, email string) (user models.User, err error) {
	const op = "service.user"

	smtp, err := s.db.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := smtp.QueryRowContext(ctx, email)

	var User models.User
	err = row.Scan(&User.ID, &User.Email, &User.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
	}

	return User, nil
}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	const op = "service.isAdmin"

	smtp, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := smtp.QueryRowContext(ctx, id)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, AppID int64) (app models.App, err error) {
	const op = "service.app"

	smtp, err := s.db.Prepare("SELECT * FROM app WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
	}

	res, err := smtp.QueryContext(ctx, AppID)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
	}

	var App models.App
	err = res.Scan(&App.ID, &App.Name, &App.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
	}

	return App, nil
}
