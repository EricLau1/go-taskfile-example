package main

import (
	"flag"
	"go-taskfile-example/exit"
	"go-taskfile-example/httpext"
	"go-taskfile-example/pb"
	"go-taskfile-example/posts/job"
	"go-taskfile-example/posts/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

const defaultGrpcPort = 3004

var (
	port        int
	accountAddr string
)

func init() {
	flag.IntVar(&port, "port", defaultGrpcPort, "posts service port")
	flag.StringVar(&accountAddr, "account_addr", "localhost:3003", "account service address")
	flag.Parse()
}

func main() {
	log.Printf("Starting posts service: Port=%d\n\n", port)

	listener, err := net.Listen("tcp", httpext.ToAddr(port))
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	svc := service.New()
	pb.RegisterPostServiceServer(grpcServer, svc)

	exit.Graceful(func() {
		grpcServer.GracefulStop()

		log.Println("Posts service stopped.")
	})

	accountConn, err := grpc.Dial(accountAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer httpext.HandleClose(accountConn)

	accountClient := pb.NewAccountServiceClient(accountConn)

	go job.Run(accountClient, svc)

	log.Fatalln(grpcServer.Serve(listener))
}
