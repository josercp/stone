package network

import (
	"context"
	"fmt"

	pb "github.com/josercp/stone/pkg/proto"
	"google.golang.org/grpc"
)

//SendPeersList Function...
func SendPeersList(ip string) bool {
	conn, err := grpc.Dial("localhost:7075", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	serviceClient := pb.NewAddServiceClient(conn)

	res, err := serviceClient.Hello(context.Background(), &pb.HelloReq{Msg: string("Hello"), Ip: string("192.168.1.21")})

	if err != nil {
		panic(err)
		return false
	}

	fmt.Println(res.MsgRes)
	return true
}
