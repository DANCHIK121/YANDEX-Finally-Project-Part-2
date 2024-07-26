@echo off
cd "Server"
go run StartServer.go DataBase.go Functions.go Json.go Structs.go GRPCServer.go GRPC.go 
pause