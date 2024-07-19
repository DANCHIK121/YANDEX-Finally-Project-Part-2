package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
    "context"
	"net/http"
	"strconv"
	"strings"
    "database/sql"
)

// Struct for request decoding
type CalculationRequest struct {
    ID         int
	Expression string 
}

// Struct for stroing request
type CalculationStore struct {
    ID     int
    Result int
    Status string
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
    ID     int
    Result int
}

// Struct for user registration
type User struct {
	Login    string
	Password string
}

var LocalID int = 0
var ExpressionsList []CalculationStore
var ExpressionRequestsList []CalculationRequest

func ValidateAndDecodeRequest(w http.ResponseWriter, r *http.Request) {
    // If the Content-Type header is present, check that it has the value
    // application/json. Note that we parse and normalize the header to remove 
    // any additional parameters (like charset or boundary information) and normalize
    // it by stripping whitespace and converting to lowercase before we check the
    // value.
    contentType := r.Header.Get("Content-Type")
    if contentType != "" {
        mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
        if mediaType != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return
        }
    }

    if r.Body == nil {
        msg := "Something went wrong"
        http.Error(w, msg, 500)
        return
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

    var storedRequest CalculationStore
    var decodedRequest CalculationRequest

    err := decoder.Decode(&decodedRequest)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        ErrorSwitching(w, syntaxError, unmarshalTypeError, err)
        return
    }

    // Call decode again, using a pointer to an empty anonymous struct as 
    // the destination. If the request body only contained a single JSON 
    // object this will return an io.EOF error. So if we get anything else, 
    // we know that there is additional data in the request body.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
        msg := "Request body must only contain a single JSON object"
        http.Error(w, msg, http.StatusBadRequest)
        return
    }

    if SearchingIdInArray(storedRequest.ID, ExpressionsList) {
        storedRequest.ID = LocalID
        decodedRequest.ID = LocalID
        LocalID++
    } 

    ExpressionsList = append(ExpressionsList, storedRequest)
    ExpressionRequestsList = append(ExpressionRequestsList, decodedRequest)

    log.Println("\n\nSuccesfully") // Succesfuly data
}

func GettingExpressoinsList (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Result Type
	message := "{\"expressions\": "+ string(UnificationExpressionsArray(ExpressionsList, w)) + " }"
	w.Write([]byte(message))

	log.Println("\n\nSuccesfully") // Succesfuly data
}

func GettingExpressoinsListForID (w http.ResponseWriter, r *http.Request) {
    var result string
    var finded bool = false
    var neededExpressionIdIndex int

    // Getting request like :id?ID=2
    neededExpressionId, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        msg := "Something went wrong"
        http.Error(w, msg, 500)
        return
    }

    for i := 0; i <= len(ExpressionsList)-1; i++ {
        if ExpressionsList[i].ID == neededExpressionId {
            finded = true
            neededExpressionIdIndex = i
        }
    }

    if !finded { http.Error(w, "There is no such expression", 404)
    } else {
        result = UnificationExpression[CalculationStore](ExpressionsList[neededExpressionIdIndex], w)
    }

    w.Header().Set("Content-Type", "application/json") // Result Type
	message := "{\"expression\": "+ result + " }"
	w.Write([]byte(message))

	log.Println("\n\nSuccesfully") // Succesfuly data
}

func GettingTask(w http.ResponseWriter, r *http.Request) {
    var task Task

    ct := r.Header.Get("Content-Type")
    if ct != "" {
        mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
        if mediaType != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return
        }

        if r.Body == nil {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
            return
        }

        jsonConsts := ReadJson()
        r.Body = http.MaxBytesReader(w, r.Body, int64(jsonConsts.MaxBytesForReader))
        decoder := json.NewDecoder(r.Body) // Decoder creating
        decoder.DisallowUnknownFields()

        var decodedRequest TaskResult

        err := decoder.Decode(&decodedRequest)
        if err != nil {
            var syntaxError *json.SyntaxError
            var unmarshalTypeError *json.UnmarshalTypeError

            ErrorSwitching(w, syntaxError, unmarshalTypeError, err)
            return
        }

        err = decoder.Decode(&struct{}{})
        if !errors.Is(err, io.EOF) {
            msg := "Request body must only contain a single JSON object"
            http.Error(w, msg, http.StatusBadRequest)
            return
        }

        result, index, err := SearchingTaskResult(decodedRequest, ExpressionsList, w)

        if err != nil {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
            return
        }

        ExpressionsList[index].Result = result
        ExpressionsList[index].Status = "Complited"
        
        log.Println("\n\nSuccesfully") // Succesfuly data
    } else {
        if len(ExpressionRequestsList) <= 0 {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
        } 

        task.ID = ExpressionsList[len(ExpressionsList)-1].ID

        decodedExpression := PostFixDecoding(ExpressionRequestsList[len(ExpressionRequestsList)-1])

        number, err := strconv.Atoi(decodedExpression[0])
        task.Arg1 = number
        if err != nil {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
            return
        }
        number, err = strconv.Atoi(decodedExpression[len(decodedExpression)-1])
        task.Arg2 = number
        if err != nil {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
            return
        }
        task.Operation = decodedExpression[1]

        // Sending result
        w.Header().Set("Content-Type", "application/json") // Result Type
        message := string(UnificationExpression[Task](task, w))
        w.Write([]byte(message))

        log.Println("\n\nSuccesfully") // Succesfuly data
    }
}

func UserRegist(w http.ResponseWriter, r *http.Request) {
    if r.Body == nil {
        msg := "Something went wrong"
        http.Error(w, msg, 500)
        return
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

    var userRegist User

    err := decoder.Decode(&userRegist)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        ErrorSwitching(w, syntaxError, unmarshalTypeError, err)
        return
    }

    // Call decode again, using a pointer to an empty anonymous struct as 
    // the destination. If the request body only contained a single JSON 
    // object this will return an io.EOF error. So if we get anything else, 
    // we know that there is additional data in the request body.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
        msg := "Request body must only contain a single JSON object"
        http.Error(w, msg, http.StatusBadRequest)
        return
    }

    ctx := context.TODO()

    db, err := sql.Open("sqlite3", "DataBase/Store.db")
    if err != nil {
        panic(err)
    }

    err = db.PingContext(ctx)
    if err != nil {
        panic(err)
    }
    
    if err = RegistUser(ctx, db, userRegist, w); err != nil {
        // panic(err)
        // log.Println(err)
    }

    db.Close()

    log.Println("\n\nSuccesfully") // Succesfuly data
}

// Server main logic
func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/internal/task", GettingTask)
    mux.HandleFunc("/api/v1/calculate", ValidateAndDecodeRequest)
	mux.HandleFunc("/api/v1/expressions", GettingExpressoinsList)
    mux.HandleFunc("/api/v1/expressions/{id}", GettingExpressoinsListForID)
    mux.HandleFunc("/api/v1/register", UserRegist)

    // Check database file is exists
    CheckDataBaseTables()

    jsonConsts := ReadJson()
    serverLocalization := fmt.Sprintf("%s:%s", jsonConsts.ServerIP, jsonConsts.ServerPort)
    log.Println(serverLocalization)

    err := http.ListenAndServe(serverLocalization, mux)
    log.Fatal(err)
}
