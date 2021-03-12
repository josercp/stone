package network

import (
	"context"
	"fmt"
	"log"

	pb "github.com/josercp/stone/pkg/proto"
	utils "github.com/josercp/stone/pkg/utils"
)

//Hello Function...
func Hello(ctx context.Context, req *pb.HelloReq) string {
	msg, ipReq := req.GetMsg(), req.GetIp()
	fmt.Println("Receiving ", msg, " message from ", ipReq)

	result := "Receiving " + msg + " response from " + ipReq
	exist := false
	var ips = make([]string, 1)

	knowNodes := utils.GetKnowNodes("known_nodes")
	searchRes := utils.Contains(knowNodes.IP, ipReq)

	if searchRes {
		exist = true
		fmt.Printf("Known node")
		return "nil"
		//Known node
	}
	ips = append(knowNodes.IP, ipReq)
	ips[0] = ipReq
	fmt.Printf("%+v\n", ips)

	if exist == false {
		//UPDATE PEERS
		knowNodesUpd := utils.SetKnowNodes(ips)
		if knowNodesUpd == false {
			log.Fatal("Error")
		}
		fmt.Printf("%+v\n", ips)
	}

	return result
}
