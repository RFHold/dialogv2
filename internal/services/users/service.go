package users

import (
	"context"
	"dialogv2/internal/database/entities"
	"dialogv2/pb"
	"dialogv2/pb/users"
	"fmt"
)

type ServiceServer struct {
	*users.UnimplementedServiceServer
	Repository *Repository
}

func (s *ServiceServer) Get(ctx context.Context, in *pb.GenericRequest) (*users.User, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	user, err := s.Repository.GetById(in.Uid)

	if err != nil {
		return nil, err
	}

	return user.Model(), err
}

func (s *ServiceServer) Create(ctx context.Context, in *users.CreateRequest) (*users.User, error) {
	user := &entities.User{FullName: in.User.FullName, Email: in.User.Email, Phone: in.User.Phone}

	err := s.Repository.Create(user)

	if err != nil {
		return nil, err
	}

	return user.Model(), err
}

func (s *ServiceServer) Update(ctx context.Context, in *users.UpdateRequest) (*users.User, error) {
	if in.User.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	user, err := s.Repository.GetById(in.User.Uid)

	if err != nil {
		return nil, err
	}

	user.FullName = in.User.FullName
	user.Email = in.User.Email
	user.Phone = in.User.Phone

	err = s.Repository.Update(user)

	if err != nil {
		return nil, err
	}

	return user.Model(), err
}

func (s *ServiceServer) Delete(ctx context.Context, in *pb.GenericRequest) (*pb.GenericResponse, error) {
	if in.Uid == "" {
		return nil, fmt.Errorf("uid must be specified")
	}

	user, err := s.Repository.GetById(in.Uid)

	if err != nil {
		return &pb.GenericResponse{Status: pb.RequestStatus_FAILURE}, err
	}

	err = s.Repository.Delete(user)

	if err != nil {
		return &pb.GenericResponse{Status: pb.RequestStatus_FAILURE}, err
	} else {
		return &pb.GenericResponse{Status: pb.RequestStatus_SUCCESS}, err
	}
}
