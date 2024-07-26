package Agent

import (
	"io"
	"os"
    "errors"
	"strconv"
    "testing"
    "encoding/json"
)

func TestUnificationTaskResult(t *testing.T) {
	testCases := []TestStructForUnificationTaskResult{
        {TaskResult{ID: 0, Result: 2, JWTToken: ""}, `{"ID":0,"Result":2,"JWTToken":""}`},    // 1
        {TaskResult{ID: 22222, Result: 5, JWTToken: "ssf"}, `{"ID":22222,"Result":5,"JWTToken":"ssf"}`},  // 2
        {TaskResult{ID: 33, Result: -5, JWTToken: "ss/s"}, `{"ID":33,"Result":-5,"JWTToken":"ss/s"}`}, // 3
    }

	for test := range testCases {
        result := UnificationTaskResult(testCases[test].in)

        if result != testCases[test].out {
            t.Error()
        }
    }
}

func TestDecodeTask(t *testing.T) {
	// Please past here admin token

	testCases := []TestStructForDecodeTask {
        {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE4MTExNTIsImlhdCI6MTcyMTgwOTk1MiwibmFtZSI6IkFkbWluIiwibmJmIjoxNzIxODEwMDEyfQ.31cFX6wEvp8I_FwGifXP-ih3JjjIRX6lBaEpfI9nKa0", Task{ID: 0, Arg1: 5, Arg2: 2, Operation: "-"}, nil},    // 1
    }

	for test := range testCases {
        result, err := DecodeTask(testCases[test].in)

        if (result != testCases[test].outResult) || (err != testCases[test].outError) {
            t.Error()
        }
    }
}

func TestErrorSwitching(t *testing.T) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	testCases := []TestStructForErrorSwitching {
        {syntaxError, unmarshalTypeError, io.ErrUnexpectedEOF, "Request body contains badly-formed JSON"},                         // 1
		{syntaxError, unmarshalTypeError, io.EOF, "Request body must not be empty"},                                               // 2
		{syntaxError, unmarshalTypeError, errors.New("http: request body too large"), "Request body must not be larger than 1MB"}, // 3
		{syntaxError, unmarshalTypeError, errors.New("Test error"), "Test error"},                                                 // 4
    }

	for test := range testCases {
        result := ErrorSwitching(testCases[test].inSyntaxError, testCases[test].inUnmarshalTypeError, testCases[test].inError)

        if (result != testCases[test].out) {
            t.Error()
        }
    }
}

func TestEnvironmentVariablesInit(t *testing.T) {
    // Read json file
    jsonConsts := ReadJson()

	testCases := TestStructFotEnvironmentVariablesInit {
        ComputingPower: jsonConsts.ComputingPower,
        TimeAdditionMS: jsonConsts.TimeAdditionMS,
        TimeSubstrationMS: jsonConsts.TimeSubstrationMS,
        TimeMultiplicationsMS: jsonConsts.TimeMultiplicationsMS,
        TimeDivisionsMS: jsonConsts.TimeDivisionsMS,
    }

    EnvironmentVariablesInit()

    // Read values
    resultComputingPower, err :=      strconv.Atoi(os.Getenv("COMPUTING_POWER"))
    resultTimeAddition, err :=        strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
    resultTimeSubstraction, err :=    strconv.Atoi(os.Getenv("TIME_SUBSTRACTION_MS"))
    resultTimeMultiplications, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
    resultTimeDivisions, err :=       strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))


    if (resultComputingPower != testCases.ComputingPower) || (resultTimeAddition != testCases.TimeAdditionMS) ||
       (resultTimeSubstraction != testCases.TimeSubstrationMS) || (resultTimeMultiplications != testCases.TimeMultiplicationsMS) ||
       (resultTimeDivisions != testCases.TimeDivisionsMS) {
        if err == nil {
            t.Error()
        } else {
            panic(err)
        }
    }
}
