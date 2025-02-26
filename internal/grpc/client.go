package grpc

import (
	"context"
	authpb "github.com/dankru/proto-definitions/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type GrpcClient struct {
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{}
}

func (g *GrpcClient) ParseToken() (int64, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient("172.17.0.1:9000", opts...)
	if err != nil {
		log.Fatalf("Не удалось установить соединение: %s", err.Error())
	}
	defer conn.Close()

	c := authpb.NewTokenServiceClient(conn)

	message := authpb.TokenRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDA2NDM5NDgsImlhdCI6MTc0MDU4OTk0OCwic3ViIjoiMyJ9.nbiNOZiBAKds-Ta1JNN19iyiRm7wIILFNgAAal9qXJ8"}

	response, err := c.ParseToken(context.Background(), &message)
	if err != nil {
		log.Fatalf("Не удалось отправить сообщение: %s", err.Error())
	}

	return response.Id, err
}
