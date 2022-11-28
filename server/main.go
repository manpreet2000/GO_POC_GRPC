package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/manpreet/grpc/transaction"
)

type server struct {
	pb.UnimplementedTransactionServiceServer
}

// considering 3 nodes
var transactionsCount = make(map[int32]int)
var transactionsHash = make(map[int32]string)
var ports = [3]string{":8085", ":3000", ":50051"}
var currentPort = flag.String("port", ":8085", "the port to connect to")

// var currentPort = ports[2]
var wg sync.WaitGroup

func (s *server) SendTransaction(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	id := req.GetTransactionId()
	hash := req.GetTransactionHash()
	if transactionsCount[id] > 1 {
		return &pb.TransactionResponse{IsValid: false}, nil
	}
	transactionsCount[id]++
	transactionsHash[id] = hash
	for _, port := range ports {
		wg.Add(1)
		go broadCast(id, hash, port)
	}
	if transactionsCount[id] > len(ports)/2 {
		fmt.Printf("Saving file for port %s and transaction count is %d", *currentPort, transactionsCount[id])
		wg.Add(1)
		go saveFile(id, hash)
	}
	fmt.Printf("Transactions count of id %v is %v \n", id, transactionsCount[id])
	return &pb.TransactionResponse{IsValid: true}, nil
}
func saveFile(id int32, hash string) {

	defer wg.Done()

	input, err := ioutil.ReadFile(*currentPort + ".txt")
	if err != nil {
		f, err := os.Create(*currentPort + ".txt")
		if err != nil {
			log.Fatal("Error in creating file ", err)
		}
		defer f.Close()

		_, err2 := f.WriteString(fmt.Sprintf("%d -> %v \n", id, hash))
		if err2 != nil {
			fmt.Printf("error writing string: %v", err)
		}
		return
	}

	input = append(input, fmt.Sprintf("%d -> %v \n", id, hash)...)
	err = ioutil.WriteFile(*currentPort+".txt", input, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
func broadCast(id int32, hash string, port string) {
	defer wg.Done()
	fmt.Println("Port ", port)
	if port == *currentPort {
		fmt.Println("same port returning ", port)
		return
	}
	addr := "localhost" + port
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendTransaction(ctx, &pb.TransactionRequest{TransactionId: id, TransactionHash: hash})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r.IsValid)
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *currentPort)
	if err != nil {
		log.Fatalln("Error ", err)
	}
	s := grpc.NewServer()
	pb.RegisterTransactionServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
