package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/domain/models"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/lib/jwt"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	userLogout   UserLogout
	tokenTTL     time.Duration
}
type UserSaver interface {
	SaveUser(
		ctx context.Context,
		uuid string,
		email string,
		passHash []byte,
		phone string,
		dateOfBirth string,
		username string,
	) (uuids string, err error)
}
type UserLogout interface {
	LogoutUser(
		ctx context.Context,
		token string,
	) (success bool, err error)
}
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userUuid string) (bool, error)
}

var (
	ErrInvalidUserCredentials = errors.New("invalid credentials")
	ErrInvalidAppUuid         = errors.New("invalid app uuid")
	ErrUserExists             = errors.New("user exists")
)

func New(
	log *slog.Logger, userSaver UserSaver,
	userProvider UserProvider,
	userLogout UserLogout,
	tokenTTl time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		userLogout:   userLogout,
		log:          log,
		tokenTTL:     tokenTTl,
	}
}
func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to login user ")

	user, err := a.userProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Warn("user not found ", err)

			return "", fmt.Errorf("%s:%w", op, ErrInvalidUserCredentials)
		}
		a.log.Error("failed to get user", err)

		return "", fmt.Errorf("%s:%w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials ", err)
		return "", fmt.Errorf("%s:%w", op, ErrInvalidUserCredentials)
	}

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err)

		return "", fmt.Errorf("%s:%w", op, err)
	}

	return token, nil
}
func (a *Auth) RegisterNewUser(ctx context.Context, email string, phone string, dateOfBirth string, password string, username string) (string, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate hash password ", err)
		return "", fmt.Errorf("%s:%w", op, err)
	}

	uid, err := uuid.NewUUID()

	if err != nil {
		log.Error("failed to create uuid")
		return "", fmt.Errorf("%s:%w", op, err)
	}

	newUuids := uid.String()

	uuid, err := a.userSaver.SaveUser(ctx, newUuids, email, passHash, phone, dateOfBirth, username)

	if err != nil {
		fmt.Println(err)

		if errors.Is(err, repository.ErrUserExists) {
			a.log.Warn("user already exists", err)

			return "", fmt.Errorf("%s:%w", op, ErrUserExists)
		}

		log.Error("failed to save user", err)
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user registered")

	UUID := uuid
	return UUID, nil
}
func (a *Auth) IsAdmin(ctx context.Context, userUuid string) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("user_uuid", userUuid),
	)

	log.Info("checking if user is admin")
	isAdmin, err := a.userProvider.IsAdmin(ctx, userUuid)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			log.Warn("app not found", err)

			return false, fmt.Errorf("%s:%w", op, ErrInvalidAppUuid)
		}
		return false, fmt.Errorf("%s:%w", op, err)
	}
	return isAdmin, nil
}
func (a *Auth) Logout(ctx context.Context, token string) (bool, error) {
	const op = "auth.LogOut"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("logging user out")

	success, err := a.userLogout.LogoutUser(ctx, token)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return success, nil

}
