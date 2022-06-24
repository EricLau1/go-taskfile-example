package job

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"go-taskfile-example/delay"
	"go-taskfile-example/pb"
	"io"
	"log"
)

const MaxPosts = 1000000

func Run(accSvc pb.AccountServiceClient, postsSvc pb.PostServiceServer) {
	ctx := context.Background()

	for i := 0; i < MaxPosts; i++ {

		accounts, err := GetAccounts(ctx, accSvc)
		if err != nil {
			log.Println("error on get accounts:", err.Error())
		} else {
			for _, account := range accounts {
				post, err := postsSvc.CreatePost(ctx, &pb.CreatePostRequest{
					AuthorId: account.Id,
					Title:    faker.Sentence(),
					Body:     faker.Paragraph(),
				})
				if err != nil {
					log.Println("error on create post: ", err.Error())
				} else {
					fmt.Printf("%+v\n", post)
				}
			}
		}

		delay.Random(10)
	}

}

func GetAccounts(ctx context.Context, accSvc pb.AccountServiceClient) ([]*pb.Account, error) {
	res, err := accSvc.GetAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		log.Println("error on get accounts:", err.Error())
	}
	var accounts []*pb.Account
	for {
		var account *pb.Account
		account, err = res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
