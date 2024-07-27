package main

import (
	"encoding/json"
)

// Agent file

// Struct for decoding tasks
type Task struct {
    ID             int
    Arg1           int
    Arg2           int
    Operation      string
}

// HTTP file

// Struct for encode task with result
type TaskResult struct {
    ID       int
    Result   int
	JWTToken string
}

// Functions test file

// Strcut for test UnificationTaskResult
type TestStructForUnificationTaskResult struct {
    in  TaskResult
    out string
}

// Struct for test DecodeTask
type TestStructForDecodeTask struct {
	in        string
	outResult Task
	outError  error
}

// Struct for test ErrorSwitching
type TestStructForErrorSwitching struct {
	inSyntaxError        *json.SyntaxError
	inUnmarshalTypeError *json.UnmarshalTypeError
	inError              error
	out                  string
}

// Struct for test EnvironmentVariablesInit
type TestStructFotEnvironmentVariablesInit struct {
	ComputingPower        int 
	TimeAdditionMS        int 
	TimeSubstrationMS     int 
	TimeMultiplicationsMS int 
	TimeDivisionsMS       int 
}

// Json test file

// Struct for test UnificationTaskResult
type TestStructForReadJson struct {
    out JsonConstsFile
}

// Http test file

// Struct for test CalculatingResult
type TestStructForCalculatingResult struct {
	in  Task
	out int
}

// Json file

// Struct for decoding json
type JsonConstsFile struct {
	ComputingPower        int `json:"computing_power"`
	TimeAdditionMS        int `json:"time_addition_ms"`
	TimeSubstrationMS     int `json:"time_substraction_ms"`
	TimeMultiplicationsMS int `json:"time_multiplications_ms"`
	TimeDivisionsMS       int `json:"time_divisions_ms"`
}
