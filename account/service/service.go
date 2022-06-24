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
	accounts *sync.Map
	pb.UnimplementedAccountServiceServer
}

func New() pb.AccountServiceServer {
	return &Service{accounts: &sync.Map{}}
}

func (s *Service) CreateAccount(_ context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {

	now := timestamppb.Now()

	account := &pb.Account{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Birthday:  req.Birthday,
		Gender:    req.Gender,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.accounts.Store(account.Id, account)

	return account, nil
}

func (s *Service) GetAccount(_ context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {

	item, exits := s.accounts.Load(req.Id)
	if !exits {
		return nil, fmt.Errorf("account not found: %s", req.Id)
	}

	return item.(*pb.Account), nil
}

func (s *Service) GetAccounts(req *pb.GetAccountsRequest, stream pb.AccountService_GetAccountsServer) error {

	s.accounts.Range(func(_, value any) bool {
		err := stream.Send(value.(*pb.Account))
		if err != nil {
			log.Println("error on get accounts:", err.Error())
		}
		return true
	})

	return nil
}

func (s *Service) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.Account, error) {
	account, err := s.GetAccount(ctx, &pb.GetAccountRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	s.accounts.Delete(req.Id)

	return account, nil
}
