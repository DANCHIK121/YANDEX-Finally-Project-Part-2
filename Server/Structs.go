package main

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
