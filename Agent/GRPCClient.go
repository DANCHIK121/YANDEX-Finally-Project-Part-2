package main

import (
	"os"
	"log"
	"fmt"
	pb "github.com/DANCHIK121/YANDEX-Finally-Project-Part-2/Server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GRPCmain() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	defer conn.Close()
}