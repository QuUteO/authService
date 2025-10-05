package auth

// пакет, где описываются ручки gRPC

import (
	pb "Auth/app/gen/go/sso"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int32) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

// ServerAPI grpc-ручки
type ServerAPI struct {
	pb.UnimplementedAuthServer
	auth Auth
}

// RegisterAuthServer регистрируем обработчик ручек
func RegisterAuthServer(gRPC *grpc.Server, auth Auth) {
	pb.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

const (
	ZeroValue = 0
)

// Login ручка для логирования пользователя в системе
func (sso *ServerAPI) Login(ctx context.Context, req *pb.LoginResponse) (*pb.LoginRequest, error) {
	if err := loginCheck(req); err != nil {
		return nil, err
	}

	token, err := sso.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed internal")
	}

	return &pb.LoginRequest{
		Token: token,
	}, nil
}

// Register ручка для регистрация пользователя
func (sso *ServerAPI) Register(ctx context.Context, req *pb.RegisterResponse) (*pb.RegisterRequest, error) {
	if err := registerCheck(req); err != nil {
		return nil, err
	}

	userID, err := sso.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed internal")
	}

	return &pb.RegisterRequest{
		UserId: userID,
	}, nil
}

func (sso *ServerAPI) IsAdmin(ctx context.Context, req *pb.AdminResponse) (*pb.AdminRequest, error) {
	if err := isAdminCheck(req); err != nil {
		return nil, err
	}

	isAdmin, err := sso.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed internal")
	}

	return &pb.AdminRequest{
		IsAdmin: isAdmin,
	}, nil
}

func loginCheck(req *pb.LoginResponse) error {
	if req.GetEmail() == "" {
		return status.Error(codes.FailedPrecondition, "Email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.FailedPrecondition, "Password is required")
	}

	if req.GetAppId() == ZeroValue {
		return status.Error(codes.FailedPrecondition, "AppId is required")
	}

	return nil
}

func registerCheck(req *pb.RegisterResponse) error {
	if req.GetEmail() == "" {
		return status.Error(codes.FailedPrecondition, "Email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.FailedPrecondition, "Password is required")
	}

	return nil
}

func isAdminCheck(req *pb.AdminResponse) error {
	if req.GetUserId() == ZeroValue {
		return status.Error(codes.FailedPrecondition, "UserId is required")
	}

	return nil
}
