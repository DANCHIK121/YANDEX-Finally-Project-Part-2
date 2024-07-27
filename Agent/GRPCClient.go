package main

import (
	"os"
	"log"
	"fmt"
	"context"
	pb "github.com/DANCHIK121/YANDEX-Finally-Project-Part-2/Agent/proto"
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

	grpcClient := pb.NewTaskServiceClient(conn)

	send, err := grpcClient.SendTask(context.TODO(), &pb.TaskRequest{
		ID: 0,
		Arg1: 1,
		Arg2: 3,
		Operation: "*",
	})
	
	if err != nil {
		log.Println("failed invoking Area: ", err)
	}
	
	update, err := grpcClient.UpdateTaskString(context.TODO(), &pb.TaskRequest{
		ID: 0,
		Arg1: 1,
		Arg2: 3,
		Operation: "*",
	})
	
	if err != nil {
		log.Println("failed invoking Area: ", err)
	}
	
	fmt.Println("Area: ", send.Arg1)
	fmt.Println("Perimeter: ", update.Result)
}