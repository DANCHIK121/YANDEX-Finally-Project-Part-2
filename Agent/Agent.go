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

type Task struct {
    ID             int
    Arg1           int
    Arg2           int
    Operation      string
}

var mx sync.Mutex
var TempTask Task
var CurrentGoRutinesCount int

func GetTaskFromServer() {
	for {
		currentTask, err := DecodeTask()
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

	response := NewRequestToServer()

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

    // Call decode again, using a pointer to an empty anonymous struct as 
    // the destination. If the request body only contained a single JSON 
    // object this will return an io.EOF error. So if we get anything else, 
    // we know that there is additional data in the request body.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
        msg := "Request body must only contain a single JSON object"
        log.Fatalln(msg)
        return
    }

	TempTask = decodedRequest

	SendResultToServer(decodedRequest.ID, CalculatingResult(decodedRequest))
}

func main() {
	EnvironmentVariablesInit()

	// COMPUTING_POWER environment variable
	computingPowerConst := os.Getenv("COMPUTING_POWER")

	mx.Lock()
	CurrentGoRutinesCount = 0
	mx.Unlock()

	if computingPowerConst == "" {
		log.Fatalf("Something went wrong: \n\n %s", fmt.Errorf("COMPUTING_POWER variable is nil"))
	}

	GetTaskFromServer()
}
