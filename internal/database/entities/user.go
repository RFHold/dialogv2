package entities

import (
	"dialogv2/pb/users"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	Base
	FullName string
	Email    string
	Phone    string
}

func (u *User) Model() *users.User {
	return &users.User{
		Uid:       u.Uid,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
		DeletedAt: timestamppb.New(u.DeletedAt.Time),
	}
}
