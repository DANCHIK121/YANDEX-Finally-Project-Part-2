package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"sync"
	"errors"
	"strings"
	"encoding/json"
)

var mx sync.Mutex
var TempTask Task
var ServerIP string
var ServerHost string
var CurrentJWTToken string
var CurrentGoRutinesCount int

func GetTaskFromServer() {
	for {
		currentTask, err := DecodeTask(CurrentJWTToken)
		if errors.Is(err, io.EOF) {
			log.Println("EOF")
		}
		if currentTask == TempTask {
			log.Println("The task is already being calculated")
		} else {
			mx.Lock()
			CurrentGoRutinesCount++
			mx.Unlock()
			TempTask = currentTask
			go CalculateTaskFromServer()
		}
	}
}

func CalculateTaskFromServer() {
	var err error
	var decodedRequest Task

	response := NewRequestToServer(CurrentJWTToken)

	contentType := response.Request.Header.Get("Content-Type") // Request content type validate
    if contentType != "" {
        mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
        if mediaType != "application/json" {
            msg := "Content-Type header is not application/json"
            log.Fatalln(msg)
            return
        }
    }

	// Creating json decoder
	if response.Body == nil {
		msg := "Something went wrong"
		log.Fatalln(msg)
		return
	}

	decoder := json.NewDecoder(response.Body)
    decoder.DisallowUnknownFields()

	// Decoding Request
    err = decoder.Decode(&decodedRequest)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        ErrorSwitching(syntaxError, unmarshalTypeError, err)
        return
    }

	TempTask = decodedRequest
	log.Println(decodedRequest)
	SendResultToServer(decodedRequest.ID, CalculatingResult(decodedRequest))
}

func main() {
	EnvironmentVariablesInit()

	// COMPUTING_POWER environment variable
	computingPowerConst := os.Getenv("COMPUTING_POWER")

	fmt.Print("Please enter your JWT token here: ")
	fmt.Fscan(os.Stdin, &CurrentJWTToken)

	fmt.Print("Please enter server host here: ")
	fmt.Fscan(os.Stdin, &ServerHost)

	fmt.Print("Please enter server ip here: ")
	fmt.Fscan(os.Stdin, &ServerIP)

	if (CurrentJWTToken != "") && (ServerHost != "") && (ServerIP != "") {

		mx.Lock()
		CurrentGoRutinesCount = 0
		mx.Unlock()

		if computingPowerConst == "" {
			log.Fatalf("Something went wrong: \n\n %s", fmt.Errorf("COMPUTING_POWER variable is nil"))
		}

		GetTaskFromServer()
	}
}
