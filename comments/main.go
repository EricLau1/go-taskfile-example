package main

import (
	"flag"
	"go-taskfile-example/comments/job"
	"go-taskfile-example/comments/service"
	"go-taskfile-example/exit"
	"go-taskfile-example/httpext"
	"go-taskfile-example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

const defaultGrpcPort = 3005

var (
	port        int
	accountAddr string
	postsAddr   string
)

func init() {
	flag.IntVar(&port, "port", defaultGrpcPort, "comments service port")
	flag.StringVar(&accountAddr, "account_addr", "localhost:3003", "account service address")
	flag.StringVar(&postsAddr, "posts_addr", "localhost:3004", "posts service address")
	flag.Parse()
}

func main() {
	log.Printf("Starting comments service: Port=%d\n\n", port)

	listener, err := net.Listen("tcp", httpext.ToAddr(port))
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	svc := service.New()
	pb.RegisterCommentServiceServer(grpcServer, svc)

	exit.Graceful(func() {
		grpcServer.GracefulStop()

		log.Println("Comments service stopped.")
	})

	accountConn, err := grpc.Dial(accountAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer httpext.HandleClose(accountConn)

	postsConn, err := grpc.Dial(postsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer httpext.HandleClose(postsConn)

	accountClient := pb.NewAccountServiceClient(accountConn)
	postsClient := pb.NewPostServiceClient(postsConn)

	go job.Run(accountClient, postsClient, svc)

	log.Fatalln(grpcServer.Serve(listener))
}
