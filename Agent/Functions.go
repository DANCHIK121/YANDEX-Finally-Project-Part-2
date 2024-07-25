package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"errors"
	"strings"
	"net/http"
	"encoding/json"
)

func EnvironmentVariablesInit() {
	// Read json
	jsonConsts := ReadJson()

	// Setting environment variables
	var err error

	err = os.Setenv("COMPUTING_POWER", fmt.Sprintf("%s", fmt.Sprint(jsonConsts.ComputingPower)))
	if err != nil {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}

	// Calculating actions
	err = os.Setenv("TIME_ADDITION_MS", fmt.Sprintf("%s", fmt.Sprint(jsonConsts.TimeAdditionMS))) // TIME_ADDITION_MS - the time of the addition operation in milliseconds
	if err != nil {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}
	err = os.Setenv("TIME_SUBSTRACTION_MS", fmt.Sprintf("%s", fmt.Sprint(jsonConsts.TimeSubstrationMS))) // TIME_SUBTRACTION_MS - the time of the subtraction operation in milliseconds
	if err != nil {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}
	err = os.Setenv("TIME_MULTIPLICATIONS_MS", fmt.Sprintf("%s", fmt.Sprint(jsonConsts.TimeMultiplicationsMS))) // TIME_MULTIPLICATIONS_MS - the time of the multiplication operation in milliseconds
	if err != nil {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}
	err = os.Setenv("TIME_DIVISIONS_MS", fmt.Sprintf("%s", fmt.Sprint(jsonConsts.TimeDivisionsMS))) // TIME_DIVISIONS_MS - the time of the division operation in milliseconds
	if err != nil {
		log.Fatalf("Something went wrong: \n\n %s", err)
	}

	log.Println("Succefuly")
}

func ErrorSwitching(syntaxError *json.SyntaxError, unmarshalTypeError *json.UnmarshalTypeError, err error ) string {
	switch {
	// Catch any syntax errors in the JSON and send an error message
	// which interpolates the location of the problem to make it
	// easier for the client to fix.
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		log.Println(msg)
		return msg

	// In some circumstances Decode() may also return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
	// is an open issue regarding this at
	// https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := "Request body contains badly-formed JSON"
		log.Println(msg)
		return msg

	// Catch any type errors, like trying to assign a string in the
	// JSON request body to a int field in our Person struct. We can
	// interpolate the relevant field name and position into the error
	// message to make it easier for the client to fix.
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		log.Println(msg)
		return msg

	// Catch the error caused by extra unexpected fields in the request
	// body. We extract the field name from the error message and
	// interpolate it in our custom error message. There is an open
	// issue at https://github.com/golang/go/issues/29035 regarding
	// turning this into a sentinel error.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		log.Println(msg)
		return msg

	// An io.EOF error is returned by Decode() if the request body is
	// empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		log.Println(msg)
		return msg

	// Catch the error caused by the request body being too large. Again
	// there is an open issue regarding turning this into a sentinel
	// error at https://github.com/golang/go/issues/30715.
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		log.Println(msg)
		return msg

	// Otherwise default to logging the error and sending a 500 Internal
	// Server Error response.
	default:
		log.Print(err.Error())
		return err.Error()
	}
}

func UnificationTaskResult(taskResult TaskResult) string {

	jsonData, err := json.Marshal(&taskResult) // Encoding calculation requests

	if err != nil {
		log.Print(err.Error())
		log.Println(http.StatusText(http.StatusInternalServerError)) // Marshaling error
		return ""
	}

	return string(jsonData)
}

func DecodeTask(jwt_token string) (Task, error) {
	var err error
	var result Task

	response := NewRequestToServer(jwt_token)

	contentType := response.Request.Header.Get("Content-Type") // Request content type validate
    if contentType != "" {
        mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
        if mediaType != "application/json" {
            msg := "Content-Type header is not application/json"
            log.Println(msg)
            return Task{}, fmt.Errorf("Content-Type header is not application/json")
        }
    }

	// Creating json decoder
	if response.Body == nil {
		msg := "something went wrong"
		log.Println(msg)
		return Task{}, fmt.Errorf(msg)
	}

	decoder := json.NewDecoder(response.Body)
    decoder.DisallowUnknownFields()

	// Decoding Request
    err = decoder.Decode(&result)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        ErrorSwitching(syntaxError, unmarshalTypeError, err)
        return Task{}, err
    }

	return result, nil
}
