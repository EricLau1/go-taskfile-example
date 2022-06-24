package main

import (
	"flag"
	"go-taskfile-example/account/job"
	"go-taskfile-example/account/service"
	"go-taskfile-example/exit"
	"go-taskfile-example/httpext"
	"go-taskfile-example/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

const defaultGrpcPort = 3003

var port int

func init() {
	flag.IntVar(&port, "port", defaultGrpcPort, "account service port")
	flag.Parse()
}

func main() {
	log.Printf("Starting account service: Port=%d\n\n", port)

	listener, err := net.Listen("tcp", httpext.ToAddr(port))
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	svc := service.New()
	pb.RegisterAccountServiceServer(grpcServer, svc)

	go job.Run(svc)

	exit.Graceful(func() {
		grpcServer.GracefulStop()

		log.Println("Account service stopped.")
	})

	log.Fatalln(grpcServer.Serve(listener))
}
