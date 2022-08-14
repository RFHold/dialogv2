package dataservice

import (
	"context"
	"dialogv2/internal/config"
	"dialogv2/pb"
	"dialogv2/pb/messages"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
)

type MessageClient struct {
	client messages.ServiceClient
}

func DialMessageClient(cfg *config.ServiceConfig) *MessageClient {
	var conn *grpc.ClientConn
	messageServer := fmt.Sprintf("localhost:%s", cfg.MessagesServicePort)
	conn, err := grpc.Dial(messageServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	c := messages.NewServiceClient(conn)

	return &MessageClient{client: c}
}

func (c MessageClient) GetMessage(uid string) (*messages.Message, error) {
	return c.client.Get(context.Background(), &pb.GenericRequest{Uid: uid})
}

func (c MessageClient) StreamMessages(cb func(event *messages.MutateEvent) error) error {
	stream, err := c.client.Stream(context.Background(), &emptypb.Empty{})

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

func (c MessageClient) CreateMessage(input *messages.CreateRequest) (*messages.Message, error) {
	return c.client.Create(context.Background(), input)
}

func (c MessageClient) UpdateMessage(input *messages.UpdateRequest) (*messages.Message, error) {
	return c.client.Update(context.Background(), input)
}

func (c MessageClient) DeleteMessage(uid string) error {
	_, err := c.client.Delete(context.Background(), &pb.GenericRequest{Uid: uid})

	return err
}
