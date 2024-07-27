@echo off
cd "Agent"
go run Agent.go Functions.go Http.go Structs.go Json.go GRPCClient.go
pause