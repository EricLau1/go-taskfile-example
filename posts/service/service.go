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
	posts *sync.Map
	pb.UnimplementedPostServiceServer
}

func New() pb.PostServiceServer {
	return &Service{posts: &sync.Map{}}
}

func (s *Service) CreatePost(_ context.Context, request *pb.CreatePostRequest) (*pb.Post, error) {
	now := timestamppb.Now()

	post := &pb.Post{
		Id:        uuid.New().String(),
		AuthorId:  request.AuthorId,
		Title:     request.Title,
		Body:      request.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.posts.Store(post.Id, post)

	return post, nil
}

func (s *Service) GetPost(_ context.Context, request *pb.GetPostRequest) (*pb.Post, error) {
	item, exists := s.posts.Load(request.Id)
	if !exists {
		return nil, fmt.Errorf("post not found: %s", request.Id)
	}
	return item.(*pb.Post), nil
}

func (s *Service) GetPosts(_ *pb.GetPostsRequest, stream pb.PostService_GetPostsServer) error {
	s.posts.Range(func(_, value any) bool {

		err := stream.Send(value.(*pb.Post))
		if err != nil {
			log.Println("error on get posts: ", err.Error())
		}

		return true
	})
	return nil
}

func (s *Service) DeletePost(ctx context.Context, request *pb.DeletePostRequest) (*pb.Post, error) {
	post, err := s.GetPost(ctx, &pb.GetPostRequest{Id: request.GetId()})
	if err != nil {
		return nil, err
	}

	s.posts.Delete(post.Id)

	return post, nil
}
