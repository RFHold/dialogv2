package messages

import (
	"context"
	"dialogv2/internal/database/entities"
	"dialogv2/pb"
	"dialogv2/pb/messages"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceServer struct {
	*messages.UnimplementedServiceServer
	Repository  *Repository
	RedisClient *redis.Client
}

func (s *ServiceServer) Stream(rq *emptypb.Empty, stream messages.Service_StreamServer) error {
	messagesList, err := s.Repository.GetAll()

	if err != nil {
		return err
	}

	for _, message := range *messagesList {
		body := messages.MutateEvent{
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
		var body messages.MutateEvent
		if err := json.Unmarshal([]byte(msg.Payload), &body); err != nil {
			return err
		}
		if err := stream.Send(&body); err != nil {
			return err
		}
	}

	return err
}

func (s *ServiceServer) Get(ctx context.Context, in *pb.GenericRequest) (*messages.Message, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	message, err := s.Repository.GetById(in.Uid)

	if err != nil {
		return nil, err
	}

	return message.Model(), err
}

func (s *ServiceServer) Create(ctx context.Context, in *messages.CreateRequest) (*messages.Message, error) {
	message := &entities.Message{Content: in.Content}

	err := s.Repository.Create(message)

	if err != nil {
		return nil, err
	}

	return message.Model(), err
}

func (s *ServiceServer) Update(ctx context.Context, in *messages.UpdateRequest) (*messages.Message, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	message, err := s.Repository.GetById(in.Uid)

	if err != nil {
		return nil, err
	}

	message.Content = in.Content

	err = s.Repository.Update(message)

	if err != nil {
		return nil, err
	}

	return message.Model(), err
}

func (s *ServiceServer) Delete(ctx context.Context, in *pb.GenericRequest) (*pb.GenericResponse, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	message, err := s.Repository.GetById(in.Uid)

	if err != nil {
		return &pb.GenericResponse{Status: pb.RequestStatus_FAILURE}, err
	}

	err = s.Repository.Delete(message)

	if err != nil {
		return &pb.GenericResponse{Status: pb.RequestStatus_FAILURE}, err
	} else {
		return &pb.GenericResponse{Status: pb.RequestStatus_SUCCESS}, err
	}
}
