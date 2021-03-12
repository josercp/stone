package grpc

import (
	"context"
	"fmt"
	"net"

	netReq "github.com/josercp/stone/pkg/network"
	pb "github.com/josercp/stone/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAddServiceServer
}

func GrpcServer() {
	listener, err := net.Listen("tcp", ":7075")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	pb.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)
	fmt.Println("Server started")

	if e := srv.Serve(listener); e != nil {
		panic(err)
	}
}

func (s *server) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloRes, error) {
	result := netReq.Hello(ctx, req)
	return &pb.HelloRes{MsgRes: result}, nil
}
