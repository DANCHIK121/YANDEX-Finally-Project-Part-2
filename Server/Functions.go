package main

import (
	"io"
	"log"
	"fmt"
	"errors"
	"strings"
	"net/http"
	"encoding/json"
)

func ErrorSwitching (w http.ResponseWriter, syntaxError *json.SyntaxError, unmarshalTypeError *json.UnmarshalTypeError, err error ) {
	switch {
	// Catch any syntax errors in the JSON and send an error message
	// which interpolates the location of the problem to make it
	// easier for the client to fix.
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		http.Error(w, msg, http.StatusUnprocessableEntity)

	// In some circumstances Decode() may also return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
	// is an open issue regarding this at
	// https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := "Request body contains badly-formed JSON"
		http.Error(w, msg, http.StatusUnprocessableEntity)

	// Catch any type errors, like trying to assign a string in the
	// JSON request body to a int field in our Person struct. We can
	// interpolate the relevant field name and position into the error
	// message to make it easier for the client to fix.
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		http.Error(w, msg, http.StatusUnprocessableEntity)

	// Catch the error caused by extra unexpected fields in the request
	// body. We extract the field name from the error message and
	// interpolate it in our custom error message. There is an open
	// issue at https://github.com/golang/go/issues/29035 regarding
	// turning this into a sentinel error.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		http.Error(w, msg, http.StatusUnprocessableEntity)

	// An io.EOF error is returned by Decode() if the request body is
	// empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		http.Error(w, msg, 500)

	// Catch the error caused by the request body being too large. Again
	// there is an open issue regarding turning this into a sentinel
	// error at https://github.com/golang/go/issues/30715.
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		http.Error(w, msg, 500)

	// Otherwise default to logging the error and sending a 500 Internal
	// Server Error response.
	default:
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
	}
}

func UnificationExpression[T any](expression T, w http.ResponseWriter) string {

	jsonData, err := json.Marshal(&expression) // Encoding calculation requests

	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), 500) // Marshaling error
	}

	return string(jsonData)
}

func UnificationExpressionsArray(expressoinsArray []CalculationStore, w http.ResponseWriter) string {
	result := "[ "

	for i := 0; i <= len(expressoinsArray)-1; i++ {
		jsonData, err := json.Marshal(&expressoinsArray[i]) // Encoding calculation requests

		if err != nil {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500) // Marshaling error
		}

		if i == len(expressoinsArray)-1 {
			result += string(jsonData)
		} else { result += string(jsonData) + ", " }
	}

	result += " ]" // End of calculation requests list

	return result
}

func PostFixDecoding(calculationRequest CalculationRequest) []string {
	var decodedExpression []string

	// Splitting expression value
	decodedExpression = strings.Split(calculationRequest.Expression, " ")

	return decodedExpression
}

func SearchingTaskResult(taskResult TaskResult, calculationStore []CalculationStore, w http.ResponseWriter) (int, int, error) {
	var err error
	var index int
	var finded bool = false

	for i := 0; i <= len(calculationStore)-1; i++ {
		if taskResult.ID == calculationStore[i].ID {
			index = i
			finded = true
		}
	}

	if !finded {
		http.Error(w, "There is no such task", 404) // Not succesfuly data
		return 0, 0, fmt.Errorf("error")
	}

	return taskResult.Result, index, err
}

func SearchingIdInArray(id int, array []CalculationStore) bool {
	var finded bool = false

	for i := 0; i <= len(array)-1; i++ {
		if array[i].ID == id {
			finded = true
		}
	}

	return finded
}
