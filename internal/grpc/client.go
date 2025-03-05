package grpc

import (
	"context"
	"fmt"
	authpb "github.com/dankru/proto-definitions/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

type GrpcClient struct {
	conn         *grpc.ClientConn
	tokenService authpb.TokenServiceClient
}

func NewGrpcClient(addr string) *GrpcClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Println("Не удалось установить соединение: %s", err.Error())
	}

	return &GrpcClient{
		tokenService: authpb.NewTokenServiceClient(conn),
		conn:         conn,
	}
}

func (g *GrpcClient) Close() error {
	return g.conn.Close()
}

func (g *GrpcClient) ParseToken(ctx context.Context, token string) (int64, error) {
	message := authpb.TokenRequest{Token: token}

	response, err := g.tokenService.ParseToken(context.Background(), &message)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			// Если ошибка аутентификации (например, токен истек)
			return 0, fmt.Errorf("Token is expired")
		}
		return 0, fmt.Errorf("Ошибка при обработке токена: %v", err)
	}

	return response.Id, err
}

func (g *GrpcClient) GenerateToken(ctx context.Context, userId int64) (string, string, error) {
	message := authpb.UserData{Id: userId}

	response, err := g.tokenService.GenerateToken(context.Background(), &message)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			// Если ошибка аутентификации (например, токен истек)
			return "", "", fmt.Errorf("Token is expired")
		}
		return "", "", fmt.Errorf("Ошибка при обработке токена: %v", err)
	}

	return response.AccessToken, response.RefreshToken, err
}

func (g *GrpcClient) RefreshToken(ctx context.Context, token string) (string, string, error) {
	message := authpb.TokenRequest{Token: token}

	response, err := g.tokenService.RefreshToken(context.Background(), &message)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			return "", "", fmt.Errorf("Token is expired")
		}
		return "", "", fmt.Errorf("Ошибка при обработке токена: %v", err)
	}

	return response.AccessToken, response.RefreshToken, err
}
