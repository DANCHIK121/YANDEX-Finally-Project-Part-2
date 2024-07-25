package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestErrorSwitching(t *testing.T) {
	var w http.ResponseWriter
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	testCases := []TestStructForErrorSwitching {
        {syntaxError, unmarshalTypeError, io.ErrUnexpectedEOF, "Request body contains badly-formed JSON"},                         // 1
		{syntaxError, unmarshalTypeError, io.EOF, "Request body must not be empty"},                                               // 2
		{syntaxError, unmarshalTypeError, errors.New("http: request body too large"), "Request body must not be larger than 1MB"}, // 3
		{syntaxError, unmarshalTypeError, errors.New("Test error"), "Test error"},                                                 // 4
    }

	for test := range testCases {
        result := ErrorSwitching(w, testCases[test].inSyntaxError, testCases[test].inUnmarshalTypeError, testCases[test].inError, true)

        if (result != testCases[test].out) {
            t.Error()
        }
    }
}

func TestUnificationTaskResult(t *testing.T) {
	var w http.ResponseWriter
	testCases := []TestStructForUnificationTaskResult{
        {TaskResult{ID: 0, Result: 2, JWTToken: ""}, `{"ID":0,"Result":2,"JWTToken":""}`},    // 1
        {TaskResult{ID: 22222, Result: 5, JWTToken: "ssf"}, `{"ID":22222,"Result":5,"JWTToken":"ssf"}`},  // 2
        {TaskResult{ID: 33, Result: -5, JWTToken: "ss/s"}, `{"ID":33,"Result":-5,"JWTToken":"ss/s"}`}, // 3
    }

	for test := range testCases {
        result := UnificationExpression[TaskResult](testCases[test].in, w)

        if result != testCases[test].out {
            t.Error()
        }
    }
}

func TestGetExpressionsList(t *testing.T) {
	testCases := TestStructForGetExpressionsList {
		// Please past here current expression data
		1, CalculationRequest{ID: 0, Expression: "5 - 2", JWTToken: ""},    // 1
    }

	result := []CalculationRequest{}
    result = GetExpressionsList(testCases.in)

    if result[0] != testCases.out {
        t.Error()
    }
}

func TestGetSolvedExpressionsList(t *testing.T) {
	testCases := TestStructForGetSolvedExpressionsList {
		// Please past here current expression data
		1, CalculationStore{ID: 2, Result: -4, Status: "Complited"},    // 1
    }

	result := []CalculationStore{}
    result = GetSolvedExpressionsList(testCases.in)

    if result[0] != testCases.out {
        t.Error()
    }
}

func TestUnificationExpressionsArray(t *testing.T) {
    var w http.ResponseWriter
	testCases := TestStructForUnificationExpressionsArray {
		// Please past here current expression data
		1, `[ {"ID":0,"Expression":"5 - 2","JWTToken":""}, {"ID":0,"Expression":"","JWTToken":""} ]`,    // 1
    }

    result := UnificationExpressionsArray(w, testCases.in)

    if result != testCases.out {
        t.Error()
    }
}

func TestPostFixDecoding(t *testing.T) {
	testCases := []TestStructForPostFixDecoding {
		{CalculationRequest{ID: 0, Expression: "5 - 2"}, []string{"5", "-", "2"}},     // 1
        {CalculationRequest{ID: 0, Expression: "10 + 2"}, []string{"10", "+", "2"}},   // 2
        {CalculationRequest{ID: 0, Expression: "3 * 3"}, []string{"3", "*", "3"}},     // 3
        {CalculationRequest{ID: 0, Expression: "52 / 55"}, []string{"52", "/", "55"}}, // 4
        {CalculationRequest{ID: 0, Expression: "-35 * 2"}, []string{"-35", "*", "2"}}, // 5
    }

    for test := range testCases {
        result := PostFixDecoding(testCases[test].in)

        for symbol := range result {
            if result[symbol] != testCases[test].out[symbol] {
                t.Error()
            }
        }
    }
}

func TestSearchingTaskResult(t *testing.T) {
    var w http.ResponseWriter
	testCases := []TestStructForSearchingTaskResult {
		{TaskResult{ID: 0, Result: 3}, []CalculationStore{{ID: 0, Result: 3, Status: "Complited"}}, 3, nil},    // 1
        {TaskResult{ID: 0, Result: 22}, []CalculationStore{{ID: 0, Result: 22, Status: "Complited"}}, 22, nil}, // 2
        {TaskResult{ID: 0, Result: -5}, []CalculationStore{{ID: 0, Result: -5, Status: "Complited"}}, -5, nil}, // 3
    }

    for test := range testCases {
        result, err := SearchingTaskResult(testCases[test].inTaskResult, testCases[test].inCalculationStore, w)

        if err != testCases[test].outError {
            t.Error()
        }

        if testCases[test].outInt != result {
            t.Error()
        }
    }
}

func TestSearchingIdInArray(t *testing.T) {
	testCases := []TestStructForSearchingIdInArray {
		{1, []CalculationStore{{ID: 0, Result: 3, Status: "Complited"}}, false},                                             // 1
        {0, []CalculationStore{{ID: 0, Result: 22, Status: "Complited"}}, true},                                             // 2
        {22, []CalculationStore{{ID: 0, Result: -5, Status: "Complited"}, {ID: 22, Result: -5, Status: "Complited"}}, true}, // 3
    }

    for test := range testCases {
        result := SearchingIdInArray(testCases[test].inID, testCases[test].inArray)

        if testCases[test].out != result {
            t.Error()
        }
    }
}

func TestSearchingTokenInJWTStore(t *testing.T) {
	testCases := []TestStructForSearchingTokenInJWTStore {
		{"first_test", []JWT{{1, "first_test"}}, true, 0},                      // 1
        {"second_test", []JWT{{1, "first_test"}}, false, 0},                    // 2
        {"third_test", []JWT{{1, "first_test"}, {2, "second_test"}}, false, 0}, // 3
    }

    for test := range testCases {
        result, index := SearchingTokenInJWTStore(testCases[test].inToken, testCases[test].inArray)

        if testCases[test].outBool != result {
            t.Error()
        }

        if testCases[test].outInt != index {
            t.Error()
        }
    }
}

func TestEncodeExpression(t *testing.T) {
	testCases := []TestStructForEncodeExpression {
		{CalculationStore{ID: 0, Result: 5, Status: "Complited"}, CalculationRequest{}, "0,5,Complited;"}, // 1
        {CalculationStore{}, CalculationRequest{ID: 0, Expression: "5 - 2"}, "0,5 - 2;"},            // 2
    }

    for test := range testCases {
        result := EncodeExpression(testCases[test].inExpression, testCases[test].inSolvedExpression)

        if testCases[test].out != result {
            t.Error()
        }
    }
}

func TestDecodeExpression(t *testing.T) {
	testCases := []TestStructForDecodeExpression {
		{"0,5,Complited", "", CalculationStore{ID: 0, Result: 5, Status: "Complited"}, CalculationRequest{}}, // 1
        {"", "0,5 - 2",       CalculationStore{}, CalculationRequest{ID: 0, Expression: "5 - 2"}},            // 2
    }

    for test := range testCases {
        result, second_result := DecodeExpression(testCases[test].inExpression, testCases[test].inSolvedExpression)

        if (testCases[test].outExpression != result) || (testCases[test].outSolvedExpression != second_result) {
            t.Error()
        }
    }
}
