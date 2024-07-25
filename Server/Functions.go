package main

import (
	"io"
	"fmt"
	"log"
	"errors"
	"context"
	"strconv"
	"strings"
	"net/http"
	"database/sql"
	"encoding/json"
)


func ErrorSwitching (w http.ResponseWriter, syntaxError *json.SyntaxError, unmarshalTypeError *json.UnmarshalTypeError, err error, isTest bool) string {
	switch {
	// Catch any syntax errors in the JSON and send an error message
	// which interpolates the location of the problem to make it
	// easier for the client to fix.
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		if !isTest { http.Error(w, msg, http.StatusUnprocessableEntity) }
		return msg

	// In some circumstances Decode() may also return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
	// is an open issue regarding this at
	// https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := "Request body contains badly-formed JSON"
		if !isTest { http.Error(w, msg, http.StatusUnprocessableEntity) }
		return msg

	// Catch any type errors, like trying to assign a string in the
	// JSON request body to a int field in our Person struct. We can
	// interpolate the relevant field name and position into the error
	// message to make it easier for the client to fix.
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		if !isTest { http.Error(w, msg, http.StatusUnprocessableEntity) }
		return msg

	// Catch the error caused by extra unexpected fields in the request
	// body. We extract the field name from the error message and
	// interpolate it in our custom error message. There is an open
	// issue at https://github.com/golang/go/issues/29035 regarding
	// turning this into a sentinel error.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		if !isTest { http.Error(w, msg, http.StatusUnprocessableEntity) }
		return msg

	// An io.EOF error is returned by Decode() if the request body is
	// empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		if !isTest { http.Error(w, msg, 500) }
		return msg

	// Catch the error caused by the request body being too large. Again
	// there is an open issue regarding turning this into a sentinel
	// error at https://github.com/golang/go/issues/30715.
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		if !isTest { http.Error(w, msg, 500) }
		return msg

	// Otherwise default to logging the error and sending a 500 Internal
	// Server Error response.
	default:
		log.Print(err.Error())
		if !isTest { http.Error(w, http.StatusText(http.StatusInternalServerError), 500) }
		return err.Error()
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

func GetExpressionsList(user_id int) []CalculationRequest {
	var result []CalculationRequest

	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "DataBase/Store.db")
	if err != nil {
		panic(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	temp := ""
	if temp, err = SelectExpression(ctx, db, user_id); err != nil {
		panic(err)
	}

	db.Close()

	separation := strings.Split(temp, ";")
	for i := 0; i <= len(separation)-1; i++ {
		_, temp := DecodeExpression("", separation[i])
		result = append(result, temp)
	}

	return result
}

func GetSolvedExpressionsList(user_id int) []CalculationStore {
	var result []CalculationStore

	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "DataBase/Store.db")
	if err != nil {
		panic(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	temp := ""
	if temp, err = SelectSolvedExpression(ctx, db, user_id); err != nil {
		panic(err)
	}

	db.Close()

	separation := strings.Split(temp, ";")
	for i := 0; i <= len(separation)-1; i++ {
		temp, _ := DecodeExpression(separation[i], "")
		result = append(result, temp)
	}

	return result
}

func UnificationExpressionsArray(w http.ResponseWriter, user_id int) string {
	result := "[ "

	array := GetExpressionsList(user_id)
	for i := 0; i <= len(array)-1; i++ {
		jsonData, err := json.Marshal(&array[i]) // Encoding calculation requests

		if err != nil {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500) // Marshaling error
		}

		if i == len(array)-1 {
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

func SearchingTaskResult(taskResult TaskResult, calculationStore []CalculationStore, w http.ResponseWriter) (int, error) {
	var err error
	var finded bool = false

	for i := 0; i <= len(calculationStore)-1; i++ {
		if taskResult.ID == calculationStore[i].ID {
			finded = true
		}
	}

	if !finded {
		http.Error(w, "There is no such task", 404) // Not succesfuly data
		return 0, fmt.Errorf("error")
	}

	return taskResult.Result, err
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

func SearchingTokenInJWTStore(token string, array []JWT) (bool, int) {
	var index int
	var finded bool = false

	for i := 0; i <= len(array)-1; i++ {
		if array[i].Token == token {
			finded = true
			index = i
			break
		}
	}

	return finded, index
}

func ReadRequestJson[T any] (w http.ResponseWriter, r *http.Request, content_type string, is_get_task bool) (T, error) {
	// If the Content-Type header is present, check that it has the value
    // application/json. Note that we parse and normalize the header to remove 
    // any additional parameters (like charset or boundary information) and normalize
    // it by stripping whitespace and converting to lowercase before we check the
    // value.
	var decodedRequest T

	if is_get_task != true {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
			if mediaType != content_type {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return decodedRequest, errors.New(msg)
			}
		}
	}

    if r.Body == nil {
        msg := "Something went wrong"
        http.Error(w, msg, 500)
        return decodedRequest, errors.New(msg)
    }

    // Use http.MaxBytesReader to enforce a maximum read of 1MB from the
    // response body. A request body larger than that will now result in
    // Decode() returning a "http: request body too large" error.
    jsonConsts := ReadJson()
    r.Body = http.MaxBytesReader(w, r.Body, int64(jsonConsts.MaxBytesForReader))

    // Setup the decoder and call the DisallowUnknownFields() method on it.
    // This will cause Decode() to return a "json: unknown field ..." error
    // if it encounters any extra unexpected fields in the JSON. Strictly
    // speaking, it returns an error for "keys which do not match any
    // non-ignored, exported fields in the destination".
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    err := decoder.Decode(&decodedRequest)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        ErrorSwitching(w, syntaxError, unmarshalTypeError, err, false)
        return decodedRequest, errors.New(err.Error())
    }

    // Call decode again, using a pointer to an empty anonymous struct as 
    // the destination. If the request body only contained a single JSON 
    // object this will return an io.EOF error. So if we get anything else, 
    // we know that there is additional data in the request body.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
        msg := "Request body must only contain a single JSON object"
        http.Error(w, msg, http.StatusBadRequest)
        return decodedRequest, errors.New(msg)
    }

	return decodedRequest, nil
}

func EncodeExpression(expression CalculationStore, solved_expression CalculationRequest) string {
	var result string = ""

	if (solved_expression == CalculationRequest{}) {
		result = fmt.Sprintf("%s,%s,%s;", fmt.Sprint(expression.ID), fmt.Sprint(expression.Result), expression.Status)
	} else {
		result = fmt.Sprintf("%s,%s;", fmt.Sprint(solved_expression.ID), fmt.Sprint(solved_expression.Expression))
	}

	return result
}

func DecodeExpression(encodedExpression string, solvedAndEncodedExpression string) (CalculationStore, CalculationRequest) {
	var err error
	var resultExpression CalculationStore
	var resultSolvedExpression CalculationRequest

	if solvedAndEncodedExpression == "" {
		if encodedExpression != "" {
			comma_separation := strings.Split(encodedExpression, ",")
			resultExpression.ID, err = strconv.Atoi(comma_separation[0])

			if err != nil {
				panic(err)
			}

			resultExpression.Result, err = strconv.Atoi(comma_separation[1])

			if err != nil {
				panic(err)
			}

			if comma_separation[2] == "" {comma_separation[2] = " "}
			resultExpression.Status = comma_separation[2]

			return resultExpression, CalculationRequest{}
		} else {
			return CalculationStore{}, CalculationRequest{}
		}
	} else {
		if solvedAndEncodedExpression != "" {
			comma_separation := strings.Split(solvedAndEncodedExpression, ",")
			resultSolvedExpression.ID, err = strconv.Atoi(comma_separation[0])

			if err != nil {
				panic(err)
			}

			resultSolvedExpression.Expression = comma_separation[1]

			return CalculationStore{}, resultSolvedExpression
		} else {
			return CalculationStore{}, CalculationRequest{}
		}
	}
}
