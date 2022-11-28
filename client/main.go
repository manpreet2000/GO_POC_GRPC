package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/manpreet/grpc/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port = flag.String("port", ":8085", "port")
var id = flag.Int("id", 1, "id")
var hash = flag.String("hash", "hehe", "hash")

func main() {
	flag.Parse()
	fmt.Println("PORT ", *port)
	conn, err := grpc.Dial("localhost"+*port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shaHash := sha256.Sum256([]byte(*hash))
	encryptedHash := hex.EncodeToString(shaHash[:])
	r, err := c.SendTransaction(ctx, &pb.TransactionRequest{TransactionId: int32(*id), TransactionHash: encryptedHash})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("is valid: %v", r.IsValid)
}
