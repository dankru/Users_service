package grpc

import (
	"context"
	"fmt"
	authpb "github.com/dankru/proto-definitions/pkg/auth"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

type GrpcClient struct {
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{}
}

func (g *GrpcClient) ParseToken(ctx context.Context, token string) (int64, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(viper.GetString("authServer.host")+viper.GetString("authServer.port"), opts...)
	if err != nil {
		log.Fatalf("Не удалось установить соединение: %s", err.Error())
	}
	defer conn.Close()

	c := authpb.NewTokenServiceClient(conn)

	message := authpb.TokenRequest{Token: token}

	response, err := c.ParseToken(context.Background(), &message)
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
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(viper.GetString("authServer.host")+viper.GetString("authServer.port"), opts...)
	if err != nil {
		log.Fatalf("Не удалось установить соединение: %s", err.Error())
	}
	defer conn.Close()

	c := authpb.NewTokenServiceClient(conn)

	message := authpb.UserData{Id: userId}

	response, err := c.GenerateToken(context.Background(), &message)
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
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(viper.GetString("authServer.host")+viper.GetString("authServer.port"), opts...)
	if err != nil {
		log.Fatalf("Не удалось установить соединение: %s", err.Error())
	}
	defer conn.Close()

	c := authpb.NewTokenServiceClient(conn)

	message := authpb.TokenRequest{Token: token}

	response, err := c.RefreshToken(context.Background(), &message)
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
