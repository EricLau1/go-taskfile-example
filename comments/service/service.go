package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-taskfile-example/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sync"
)

type Service struct {
	comments *sync.Map
	pb.UnimplementedCommentServiceServer
}

func New() pb.CommentServiceServer {
	return &Service{comments: &sync.Map{}}
}

func (s *Service) CreateComment(_ context.Context, request *pb.CreateCommentRequest) (*pb.Comment, error) {

	now := timestamppb.Now()

	comment := &pb.Comment{
		Id:        uuid.New().String(),
		AuthorId:  request.AuthorId,
		PostId:    request.PostId,
		Text:      request.Text,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.comments.Store(comment.Id, comment)

	return comment, nil
}

func (s *Service) GetComment(_ context.Context, request *pb.GetCommentRequest) (*pb.Comment, error) {
	item, exits := s.comments.Load(request.Id)
	if !exits {
		return nil, fmt.Errorf("comment not found: %s", request.Id)
	}
	return item.(*pb.Comment), nil
}

func (s *Service) GetComments(_ *pb.GetCommentsRequest, stream pb.CommentService_GetCommentsServer) error {
	s.comments.Range(func(_, value any) bool {
		err := stream.Send(value.(*pb.Comment))
		if err != nil {
			log.Println("error on get comments:", err.Error())
		}
		return true
	})
	return nil
}

func (s *Service) DeleteComment(ctx context.Context, request *pb.DeleteCommentRequest) (*pb.Comment, error) {
	comment, err := s.GetComment(ctx, &pb.GetCommentRequest{Id: request.GetId()})
	if err != nil {
		return nil, err
	}
	s.comments.Delete(comment.Id)
	return comment, nil
}
