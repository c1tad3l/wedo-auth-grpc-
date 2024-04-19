package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/domain/models"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/repository"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}
func (s *Storage) SaveUser(ctx context.Context, uuid string, email string, passHash string, phone string, dateOfBirth string, username string) (uuids string, err error) {
	const op = "repository.postgres.SaveUser"

	sqlQuery := `INSERT INTO users(uuid ,email,phone,dateOfBirth, pass_hash,username) VALUES($1,$2,$3,$4,$5,$6)`

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	res, err := s.db.ExecContext(ctx, sqlQuery, uuid, email, phone, dateOfBirth, passHash, username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(""+
		"result", res)

	return uuid, nil
}
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgres.User"

	var user models.User

	err := s.db.QueryRowContext(ctx, "SELECT uuid,email,pass_hash,phone,dateOfBirth,username FROM users WHERE email =$1", email).Scan(
		&user.Uuid,
		&user.Email,
		&user.PassHash,
		&user.Phone,
		&user.DateOfBirth,
		&user.Username,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
func (s *Storage) IsAdmin(ctx context.Context, userUuid string) (bool, error) {
	const op = "repository.postgres.IsAdmin"

	var isAdmin bool

	err := s.db.QueryRowContext(ctx, "SELECT is_admin FROM users WHERE uuid = $1", userUuid).Scan(&isAdmin)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
