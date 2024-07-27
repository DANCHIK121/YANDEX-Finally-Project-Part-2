package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CalculatingResult(task Task) int {
	var result int = task.Arg1

	// TIME_ADDITION_MS - the time of the addition operation in milliseconds
	// TIME_SUBTRACTION_MS - the time of the subtraction operation in milliseconds
	// TIME_MULTIPLICATIONS_MS - the time of the multiplication operation in milliseconds
	// TIME_DIVISIONS_MS - the time of the division operation in milliseconds

	// Calculating actions environment variables
	TIME_ADDITION_MS, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if (err != nil) {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}
	TIME_DIVISIONS_MS, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if (err != nil) {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}
	TIME_SUBTRACTION_MS, err := strconv.Atoi(os.Getenv("TIME_SUBSTRACTION_MS"))
	if (err != nil) {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}
	TIME_MULTIPLICATIONS_MS, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if (err != nil) {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}

	switch (task.Operation) {
	// Valide action signs
	case "+":
		result += task.Arg2
		time.Sleep(time.Duration(TIME_ADDITION_MS))
	case "-": 
		result -= task.Arg2 
		time.Sleep(time.Duration(TIME_SUBTRACTION_MS))
	case "*":
		result *= task.Arg2
		time.Sleep(time.Duration(TIME_MULTIPLICATIONS_MS))
	case "/":
		result /= task.Arg2
		time.Sleep(time.Duration(TIME_DIVISIONS_MS))
	case " ":
		log.Fatalln("Somthing went wrong") // Not succesfuly data
	default:
		result = task.Arg1
	}

	return result
}

func NewRequestToServer(jwt_token string) *http.Response {
	client := http.Client{}

	var jsonStr = []byte(fmt.Sprintf(`{"JWTToken":"%s"}`, jwt_token))
	request, err := http.NewRequest("POST", "http://server-conteiner:8888/internal/task", bytes.NewBuffer(jsonStr))

    if err != nil {
        log.Fatal(err)
    }
	
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}

	return response
}

func SendResultToServer(id, result int) {
	var taskResult TaskResult

	client := http.Client{Timeout: time.Duration(3) * time.Second}

	taskResult.ID = id
	taskResult.Result = result
	taskResult.JWTToken = CurrentJWTToken
	// message := fmt.Sprintf("{ "+ string(UnificationTaskResult(taskResult)) + ", \"JWTToken\": \"%s\" }", CurrentJWTToken)

	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(jsonData)

	response, _ := client.Post("http://server-conteiner:8888/internal/task", "application/json", body)
	response.Body.Close()

	mx.Lock()
	CurrentGoRutinesCount--
	mx.Unlock()

	log.Println("\n\nSuccesfully") // Succesfuly data
}
