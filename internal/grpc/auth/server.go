package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/repository"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/service/auth"
	authV1 "github.com/c1tad3l/wedo-auth-grpc-/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"regexp"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
	RegisterNewUser(ctx context.Context, email string, password string, phone string, dateOfBirth string, username string) (uerUuid string, err error)
	IsAdmin(ctx context.Context, userUuid string) (bool, error)
	Logout(ctx context.Context, token string) (bool, error)
}
type serverApi struct {
	authV1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authV1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

const (
	emptyValue = ""
)

func (s *serverApi) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidUserCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	if err != nil {
		fmt.Println("gRPC call failed:", err)
	}

	md := metadata.Pairs("Authorization", accessToken)
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	md = metadata.Pairs("RefreshToken", refreshToken)
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (s *serverApi) Register(ctx context.Context,
	req *authV1.RegisterRequest) (*authV1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userUuid, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword(), req.GetPhone(), req.GetDateOfBirth(), req.GetUsername())
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authV1.RegisterResponse{
		UserUuid: userUuid,
	}, nil
}
func (s *serverApi) IsAdmin(ctx context.Context, req *authV1.IsAdminRequest) (*authV1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserUuid())

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authV1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
func (s *serverApi) Logout(ctx context.Context, req *authV1.LogoutRequest) (*authV1.LogoutResponse, error) {
	if err := validateToken(req); err != nil {
		return nil, err
	}

	success, err := s.auth.Logout(ctx, req.GetAccessToken())
	var header metadata.MD
	header.Delete("Authorization")

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "не авторизован")
	}
	return &authV1.LogoutResponse{
		Success: success,
	}, nil
}
func validateLogin(req *authV1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil

}

func validateRegister(req *authV1.RegisterRequest) error {

	if req.GetEmail() == "" || !checkingEmailReg(req.GetEmail()) {
		return status.Error(codes.InvalidArgument, "email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	return nil

}

func validateIsAdmin(req *authV1.IsAdminRequest) error {
	userUuid := req.GetUserUuid()
	if userUuid == emptyValue {
		return status.Error(codes.InvalidArgument, "user_uuid is required ")
	}
	return nil
}

func checkingEmailReg(email string) bool {

	matched, _ := regexp.MatchString(`([A-Za-z0-9_\-.])+@([A-Za-z0-9_\-.])+\.([A-Za-z]{2,4})`, email)

	if !matched {

		return false
	}
	return true
}

func validateToken(req *authV1.LogoutRequest) error {
	token := req.GetAccessToken()
	if token == emptyValue {
		return status.Error(codes.NotFound, "no token found")
	}
	return nil
}
