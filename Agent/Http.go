package main

import (
	"os"
	"log"
	"time"
	"bytes"
	"strconv"
	"net/http"
	"encoding/json"
)

type TaskResult struct {
    ID     int
    Result int
}

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
	TIME_SUBTRACTION_MS, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
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

func NewRequestToServer() *http.Response {
	client := http.Client{}
	request, err := http.NewRequest("POST", "http://localhost:8000/internal/task", nil)

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
	message := "{ "+ string(UnificationTaskResult(taskResult)) + " }"

	log.Println(message)

	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(jsonData)

	response, _ := client.Post("http://localhost:8000/internal/task", "application/json", body)
	response.Body.Close()

	mx.Lock()
	CurrentGoRutinesCount--
	mx.Unlock()

	log.Println("\n\nSuccesfully") // Succesfuly data
}
