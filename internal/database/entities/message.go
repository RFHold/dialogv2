package entities

import (
	"dialogv2/pb/messages"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Base
	Content string
}

func (m *Message) Model() *messages.Message {
	return &messages.Message{
		Uid:       m.Uid,
		Content:   m.Content,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
		DeletedAt: timestamppb.New(m.DeletedAt.Time),
	}
}
