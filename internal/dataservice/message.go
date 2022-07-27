package dataservice

import (
	"context"
	"dialog/internal/config"
	"dialog/pb/messages"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
)

type MessageClient struct {
	client messages.MessageServiceClient
}

func DialMessageClient(cfg *config.ServiceConfig) *MessageClient {
	var conn *grpc.ClientConn
	messageServer := fmt.Sprintf("localhost:%s", cfg.MessagesServicePort)
	conn, err := grpc.Dial(messageServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	c := messages.NewMessageServiceClient(conn)

	return &MessageClient{client: c}
}

func (c MessageClient) GetMessage(uid string) (*messages.Message, error) {
	return c.client.GetMessage(context.Background(), &messages.MessageRequest{Uid: uid})
}

func (c MessageClient) StreamMessages(cb func(event *messages.MessageMutateEvent) error) error {
	stream, err := c.client.StreamMessages(context.Background(), &emptypb.Empty{})

	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := cb(message); err != nil {
			return err
		}
	}

	return err
}

func (c MessageClient) CreateMessage(input *messages.CreateMessageRequest) (*messages.Message, error) {
	return c.client.CreateMessage(context.Background(), input)
}

func (c MessageClient) UpdateMessage(input *messages.UpdateMessageRequest) (*messages.Message, error) {
	return c.client.UpdateMessage(context.Background(), input)
}

func (c MessageClient) DeleteMessage(uid string) error {
	_, err := c.client.DeleteMessage(context.Background(), &messages.MessageRequest{Uid: uid})

	return err
}
