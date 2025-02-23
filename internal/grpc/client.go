package grpc

import (
	"context"
	authpb "github.com/dankru/proto-definitions/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ParseToken() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(":9000", opts...)
	if err != nil {
		log.Fatalf("Не удалось установить соединение: %s", err.Error())
	}
	defer conn.Close()

	c := authpb.NewTokenServiceClient(conn)

	message := authpb.TokenRequest{Token: "test message"}
	response, err := c.ParseToken(context.Background(), &message)
	if err != nil {
		log.Fatalf("Не удалось отправить сообщение")
	}

	log.Printf(response.Id)
}
