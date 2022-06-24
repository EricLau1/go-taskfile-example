package job

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"go-taskfile-example/delay"
	"go-taskfile-example/pb"
	"go-taskfile-example/posts/job"
	"io"
	"log"
)

const MaxComments = 1000000

func Run(accSvc pb.AccountServiceClient, postSvc pb.PostServiceClient, commentSvc pb.CommentServiceServer) {
	ctx := context.Background()

	for i := 0; i < MaxComments; i++ {
		posts, err := GetPosts(ctx, postSvc)
		if err != nil {
			log.Println("error on get posts: ", err.Error())
		} else {
			accounts, err := job.GetAccounts(ctx, accSvc)
			if err != nil {
				log.Println("error on get accounts: ", err.Error())
			} else {
				doComments(ctx, accounts, posts, commentSvc)
				delay.Random(10)
			}
		}
	}
}

func GetPosts(ctx context.Context, postSvc pb.PostServiceClient) ([]*pb.Post, error) {
	res, err := postSvc.GetPosts(ctx, &pb.GetPostsRequest{})
	if err != nil {
		return nil, err
	}
	var posts []*pb.Post
	for {
		post, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func doComments(ctx context.Context, accounts []*pb.Account, posts []*pb.Post, svc pb.CommentServiceServer) {
	for _, post := range posts {
		for _, account := range accounts {
			comment := &pb.CreateCommentRequest{
				AuthorId: account.Id,
				PostId:   post.Id,
				Text:     faker.Word(),
			}
			res, err := svc.CreateComment(ctx, comment)
			if err != nil {
				log.Println("error on create comment: ", err.Error())
			} else {
				fmt.Printf("%+v\n", res)
			}
		}
	}
}
