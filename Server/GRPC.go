package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	pb "github.com/DANCHIK121/YANDEX-Finally-Project-Part-2/Server/proto"
)

func NewServer(JWTToken string, TaskRes TaskResult, w http.ResponseWriter) *Server {
	return &Server{JWTToken: JWTToken, TaskRes: TaskRes, w: w}
}

func (s *Server) SendTask(
	ctx context.Context, 
	in *pb.TaskRequest,
) (*pb.SendTaskMessage, error) {

	log.Println("invoked GetTask: ", in)
	result, index := SearchingTokenInJWTStore(s.JWTToken, JWTTokensStore)
	if result {
		if len(GetExpressionsList(JWTTokensStore[index].UserId)) <= 0 {
			msg := "Something went wrong"
			return &pb.SendTaskMessage{}, errors.New(msg)
		}

		decodedExpression := PostFixDecoding(GetExpressionsList(JWTTokensStore[index].UserId)[len(GetExpressionsList(JWTTokensStore[index].UserId))-2])
		number, err := strconv.Atoi(decodedExpression[0])
		if err != nil {
			msg := "Something went wrong"
			return &pb.SendTaskMessage{}, errors.New(msg)
		}

		second_number, err := strconv.Atoi(decodedExpression[len(decodedExpression)-1])
		if err != nil {
			msg := "Something went wrong"
			return &pb.SendTaskMessage{}, errors.New(msg)
		}

		return &pb.SendTaskMessage{
			ID: float32(GetExpressionsList(JWTTokensStore[index].UserId)[len(GetExpressionsList(JWTTokensStore[index].UserId))-1].ID),
			Arg1: float32(number),
			Arg2: float32(second_number),
			Operation: decodedExpression[1],
		}, nil
	} else {
		// Sending result
		log.Println("\n\nJWT token is not founded") // Succesfuly data
		return &pb.SendTaskMessage{}, errors.New("JWT token is not founded")
	}
}

func (s *Server) UpdateTaskString(
	ctx context.Context, 
	in *pb.TaskRequest,
) (*pb.UpdateTaskStringMessage, error) {

	log.Println("invoked UpdateTaskString: ", in)
	result, index := SearchingTokenInJWTStore(s.JWTToken, JWTTokensStore)
	if result {
		decodedRequest := s.TaskRes
		result, err := SearchingTaskResult(decodedRequest, GetSolvedExpressionsList(JWTTokensStore[index].UserId), s.w)

		if err != nil {
			msg := "Something went wrong"
			return &pb.UpdateTaskStringMessage{}, errors.New(msg)
		}

		ctx := context.TODO()

		db, err := sql.Open("sqlite3", "DataBase/Store.db")
		if err != nil {
			panic(err)
		}

		err = db.PingContext(ctx)
		if err != nil {
			panic(err)
		}

		encodeResult := EncodeExpression(CalculationStore{ID: JWTTokensStore[index].UserId, Result: result, Status: "Complited"} , CalculationRequest{})

		err = UpdateExpressionLine(ctx, db, JWTTokensStore[index].UserId, "", encodeResult)
		if err != nil {
			panic(err)
		}

		db.Close()

		return &pb.UpdateTaskStringMessage{}, nil
	} else {
		// Sending result
		s.w.Write([]byte("JWT token is not founded"))

		log.Println("\n\nJWT token is not founded") // Succesfuly data
		return &pb.UpdateTaskStringMessage{}, errors.New("JWT token is not founded")
	}
}
