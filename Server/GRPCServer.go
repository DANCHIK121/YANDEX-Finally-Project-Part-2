package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"net/http"
	pb "github.com/DANCHIK121/YANDEX-Finally-Project-Part-2/Server/proto"
	"google.golang.org/grpc"
)

func GRPCmain(JWTToken string, w http.ResponseWriter, TaskRes TaskResult) {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}
	
	log.Println("tcp listener started at port: ", port)

	grpcServer := grpc.NewServer()
	TaskServiceServer := NewServer(JWTToken, TaskRes, w)

	pb.RegisterTaskServiceServer(grpcServer, TaskServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}