package main

import (
	"context"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/app"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/config"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/lib/logger/handlers/slogpretty"
	authV1 "github.com/c1tad3l/wedo-auth-grpc-/pkg/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := SetupLogger(cfg.Env)

	log.Info("starting application",
		slog.Any("cfg", cfg),
	)

	application := app.New(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	//gateway

	conn, err := grpc.NewClient(
		"0.0.0.0:45044",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Error("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()

	err = authV1.RegisterAuthHandler(context.Background(), gwmux, conn)

	if err != nil {
		log.Error("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	err = gwServer.ListenAndServe()
	if err != nil {
		log.Error("Failed to listen")
		return
	}

	//Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped ")

}
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log

}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
