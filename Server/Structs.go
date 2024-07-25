package main

import (
    "encoding/json"
)

// Start server file

// Struct for request decoding
type CalculationRequest struct {
    ID         int
	Expression string 
    JWTToken   string
}

// Struct for stroing request
type CalculationStore struct {
    ID         int
    Result     int
    Status     string
	JWTToken   JWT
}

// Struct for getting task
type Task struct {
    ID             int
    Arg1           int
    Arg2           int
    Operation      string
}

// Struct for getting result for task
type TaskResult struct {
    ID       int
    Result   int
    JWTToken string
}

// Struct for user registration
type User struct {
    ID       int
	Login    string
	Password string
}

// Struct for select past user id 
type PastUserID struct {
    ID       int
}

// Struct for JWT tokens
type JWT struct {
    UserId int
    Token  string
}

// Struct fot only JWT tokens
type OnlyJWTTokens struct {
    JWTToken string
}

// Functions test file

// Strcut for test UnificationTaskResult
type TestStructForUnificationTaskResult struct {
    in  TaskResult
    out string
}

// Struct for test ErrorSwitching
type TestStructForErrorSwitching struct {
	inSyntaxError        *json.SyntaxError
	inUnmarshalTypeError *json.UnmarshalTypeError
	inError              error
	out                  string
}

// Struct for test GetExpressionsList
type TestStructForGetExpressionsList struct {
	in  int
	out CalculationRequest
}

// Struct for test GetSolvedExpressionsList
type TestStructForGetSolvedExpressionsList struct {
	in  int
	out CalculationStore
}

// Struct for test UnificationExpressionsArray
type TestStructForUnificationExpressionsArray struct {
	in  int
	out string
}

// Struct for test PostFixDecoding
type TestStructForPostFixDecoding struct {
	in  CalculationRequest
	out []string
}

// Struct for test PostFixDecoding
type TestStructForSearchingTaskResult struct {
	inTaskResult       TaskResult
    inCalculationStore []CalculationStore
	outInt             int
    outError           error
}

// Struct for test SearchingIdInArray
type TestStructForSearchingIdInArray struct {
	inID    int
    inArray []CalculationStore
	out     bool
}

// Struct for test SearchingTokenInJWTStore
type TestStructForSearchingTokenInJWTStore struct {
	inToken string
    inArray []JWT
	outBool bool
    outInt  int
}

// Struct for test EncodeExpression
type TestStructForEncodeExpression struct {
	inExpression       CalculationStore
    inSolvedExpression CalculationRequest
    out                string
}

// Struct for test DecodeExpression
type TestStructForDecodeExpression struct {
	inExpression        string
    inSolvedExpression  string
    outExpression       CalculationStore
    outSolvedExpression CalculationRequest
}

// Json test file

// Struct for test UnificationTaskResult
type TestStructForReadJson struct {
    out JsonConstsFile
}

// Database file

// Struct for searching expression
type Expression struct {
	Expression string
}

// Json file

// Struct for decoding json
type JsonConstsFile struct {
	ServerIP          string `json:"server_ip"`
	ServerPort        string `json:"server_port"`
	ServerJWTTokenEXP int `json:"jwt_token_exp"`
	MaxBytesForReader int `json:"server_max_bytes"`
}
