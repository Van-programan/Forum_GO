package controller

import (
	"context"
	"fmt"

	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/proto/forumservice"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ForumGRPCController struct {
	forumservice.UnimplementedForumServiceServer
	forumUC usecase.Forum
}

func NewForumGRPCController(forumUC usecase.Forum) *ForumGRPCController {
	return &ForumGRPCController{forumUC: forumUC}
}

func (c *ForumGRPCController) CreateTopic(ctx context.Context, req *forumservice.CreateTopicRequest) (*forumservice.Topic, error) {
	topic, err := c.forumUC.CreateTopic(ctx, req.Title, req.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %v", err)
	}

	return &forumservice.Topic{
		Id:        topic.ID,
		Title:     topic.Title,
		AuthorId:  topic.AuthorID,
		CreatedAt: timestamppb.New(topic.CreatedAt),
	}, nil
}

func (c *ForumGRPCController) GetTopic(ctx context.Context, req *forumservice.GetTopicRequest) (*forumservice.Topic, error) {
	topic, err := c.forumUC.GetTopicByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %v", err)
	}

	return &forumservice.Topic{
		Id:        topic.ID,
		Title:     topic.Title,
		AuthorId:  topic.AuthorID,
		CreatedAt: timestamppb.New(topic.CreatedAt),
	}, nil
}

func (c *ForumGRPCController) GetTopics(ctx context.Context, req *forumservice.GetTopicsRequest) (*forumservice.GetTopicsResponse, error) {
	topics, err := c.forumUC.GetTopics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get topics: %v", err)
	}

	protoTopics := make([]*forumservice.Topic, len(topics))
	for i, topic := range topics {
		protoTopics[i] = &forumservice.Topic{
			Id:        topic.ID,
			Title:     topic.Title,
			AuthorId:  topic.AuthorID,
			CreatedAt: timestamppb.New(topic.CreatedAt),
		}
	}

	return &forumservice.GetTopicsResponse{
		Topics: protoTopics,
	}, nil
}

func (c *ForumGRPCController) CreateMessage(ctx context.Context, req *forumservice.CreateMessageRequest) (*forumservice.Message, error) {
	message, err := c.forumUC.CreateMessage(ctx, req.TopicId, req.UserId, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	return &forumservice.Message{
		Id:        message.ID,
		TopicId:   message.TopicID,
		UserId:    message.UserID,
		Content:   message.Content,
		CreatedAt: timestamppb.New(message.CreatedAt),
	}, nil
}

func (c *ForumGRPCController) GetMessages(ctx context.Context, req *forumservice.GetMessagesRequest) (*forumservice.GetMessagesResponse, error) {
	messages, err := c.forumUC.GetMessages(ctx, req.TopicId)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}

	protoMessages := make([]*forumservice.Message, len(messages))
	for i, msg := range messages {
		protoMessages[i] = &forumservice.Message{
			Id:        msg.ID,
			TopicId:   msg.TopicID,
			UserId:    msg.UserID,
			Content:   msg.Content,
			CreatedAt: timestamppb.New(msg.CreatedAt),
		}
	}

	return &forumservice.GetMessagesResponse{
		Messages: protoMessages,
	}, nil
}

func (c *ForumGRPCController) DeleteTopic(ctx context.Context, req *forumservice.DeleteTopicRequest) (*emptypb.Empty, error) {
	if err := c.forumUC.DeleteTopic(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete topic: %v", err)
	}
	return &emptypb.Empty{}, nil
}
