package node

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	pb "github.com/josercp/stone/pkg/proto"

	badger "github.com/dgraph-io/badger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	dbPath = "github.com/josercp/stone/badger/db"
)

type IPS struct {
	IP []string
}

func (ip IPS) encodeIPS() []byte {
	data, err := json.Marshal(ip)
	if err != nil {
		panic(err)
	}

	return data
}

func decodeIPS(data []byte) (IPS, error) {
	var i IPS
	err := json.Unmarshal(data, &i)
	return i, err
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func SendPeersList(ip string) bool {
	conn, err := grpc.Dial("localhost:7075", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	serviceClient := pb.NewAddServiceClient(conn)

	res, err := serviceClient.Hello(context.Background(), &pb.HelloReq{Msg: string("Hello"), Ip: string("192.168.1.21")})

	if err != nil {
		panic(err)
	}

	fmt.Println(res.MsgRes)
}

type server struct {
	pb.UnimplementedAddServiceServer
}

func main() {

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
	msg, ipReq := req.GetMsg(), req.GetIp()
	fmt.Println("Receiving ", msg, " message from ", ipReq)

	result := "Receiving " + msg + " response from " + ipReq

	exist := false

	//DB CONFIG
	options := badger.DefaultOptions(dbPath)
	options.ValueDir = dbPath
	options.Logger = nil

	//OPEN BD
	db, errOpen := badger.Open(options)
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	defer db.Close()

	var ips = make([]string, 1)

	//GET VALUES
	errView := db.View(func(txn *badger.Txn) error {
		item, errGet := txn.Get([]byte("known_nodes"))
		if errGet != nil {
			return errGet
			//log.Fatal(errGet)
		}

		var valCopy []byte
		errValue := item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		if errValue != nil {
			log.Fatal(errValue)
		}

		valCopy, errValue = item.ValueCopy(nil)
		if errValue != nil {
			log.Fatal(errValue)
		}

		ip, errDec := decodeIPS(valCopy)
		if errDec != nil {
			log.Fatal(errDec)
		}

		searchRes := Contains(ip.IP, ipReq)
		if searchRes {
			exist = true
			return nil
			//Known node
		}
		ips = append(ip.IP, ipReq)

		//fmt.Printf("%+v\n", ip.IP)
		//fmt.Printf("Adding peer: %s\n", valCopy)

		return nil
	})
	if errView != nil {
		ips[0] = ipReq
		errUpd := db.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte("known_nodes"), IPS{
				IP: ips,
			}.encodeIPS())
			errSet := txn.SetEntry(e)
			return errSet
		})
		if errUpd != nil {
			log.Fatal(errUpd)
		}
		fmt.Printf("%+v\n", ips)
		//log.Fatal(errView)
	} else {
		if exist == false {
			//UPDATE PEERS
			errUpd := db.Update(func(txn *badger.Txn) error {
				e := badger.NewEntry([]byte("known_nodes"), IPS{
					IP: ips,
				}.encodeIPS())
				errSet := txn.SetEntry(e)
				return errSet
			})
			if errUpd != nil {
				log.Fatal(errUpd)
			}
			fmt.Printf("%+v\n", ips)
		} else {
			fmt.Printf("Known node")
		}
	}
	return &pb.HelloRes{MsgRes: result}, nil
}
