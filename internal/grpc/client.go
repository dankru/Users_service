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
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{}
}

func (g *GrpcClient) ParseToken(ctx context.Context, token string) (int64, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient("172.17.0.1:9000", opts...)
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
		// Для других ошибок просто передаем их
		return 0, fmt.Errorf("Ошибка при обработке токена: %v", err)
		// Для других ошибок просто передаем их
	}

	return response.Id, err
}
