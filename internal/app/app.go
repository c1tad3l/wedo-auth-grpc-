package app

import (
	"context"
	grpcapp "github.com/c1tad3l/wedo-auth-grpc-/internal/app/grpc"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/repository/postgresql"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/service/auth"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

type UserLogout struct {
	LogoutTime time.Time
}

func (u *UserLogout) LogoutUser(ctx context.Context, token string) (bool, error) {
	u.LogoutTime = time.Now()
	return true, nil
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgresql.New(storagePath)
	if err != nil {
		panic(err)
	}
	userLogout := &UserLogout{}

	authservice := auth.New(log, storage, storage, userLogout, tokenTTL)
	grpcApp := grpcapp.New(log, authservice, grpcPort)

	return &App{

		GRPCSrv: grpcApp,
	}
}
