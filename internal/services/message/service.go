package message

import (
	"context"
	"dialogv2/internal/database/entities"
	"dialogv2/pb"
	"dialogv2/pb/messages"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type ServiceServer struct {
	*messages.UnimplementedMessageServiceServer
	Repository  *Repository
	RedisClient *redis.Client
}

func (s *ServiceServer) StreamMessages(rq *emptypb.Empty, stream messages.MessageService_StreamMessagesServer) error {
	messagesList, err := s.Repository.GetAllMessages()

	if err != nil {
		return err
	}

	for _, message := range *messagesList {
		body := messages.MessageMutateEvent{
			Body: message.Model(),
			Type: pb.MutateEventType_FETCH,
			Id:   message.Uid,
		}
		if err := stream.Send(&body); err != nil {
			return err
		}
	}

	pubsub := s.RedisClient.Subscribe(context.Background(), MutateChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var body messages.MessageMutateEvent
		if err := json.Unmarshal([]byte(msg.Payload), &body); err != nil {
			return err
		}
		if err := stream.Send(&body); err != nil {
			return err
		}
	}

	return err
}

func (s *ServiceServer) GetMessage(ctx context.Context, in *messages.MessageRequest) (*messages.Message, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	message, err := s.Repository.GetMessageById(in.Uid)

	if err != nil {
		return nil, err
	}

	return message.Model(), err
}

func (s *ServiceServer) CreateMessage(ctx context.Context, in *messages.CreateMessageRequest) (*messages.Message, error) {
	message := &entities.Message{Content: in.Content}

	log.Fatal(in.UserId)
	err := s.Repository.CreateMessage(message)

	if err != nil {
		return nil, err
	}

	return message.Model(), err
}

func (s *ServiceServer) UpdateMessage(ctx context.Context, in *messages.UpdateMessageRequest) (*messages.Message, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	message, err := s.Repository.GetMessageById(in.Uid)

	if err != nil {
		return nil, err
	}

	message.Content = in.Content

	err = s.Repository.UpdateMessage(message)

	if err != nil {
		return nil, err
	}

	return message.Model(), err
}

func (s *ServiceServer) DeleteMessage(ctx context.Context, in *messages.MessageRequest) (*pb.GenericResponse, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	message, err := s.Repository.GetMessageById(in.Uid)

	if err != nil {
		return &pb.GenericResponse{Status: pb.RequestStatus_FAILURE}, err
	}

	err = s.Repository.DeleteMessage(message)

	if err != nil {
		return &pb.GenericResponse{Status: pb.RequestStatus_FAILURE}, err
	} else {
		return &pb.GenericResponse{Status: pb.RequestStatus_SUCCESS}, err
	}
}
