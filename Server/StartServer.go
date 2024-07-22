package main

import (
    "io"
	"fmt"
	"log"
    "time"
    "errors"
    "context"
	"strconv"
	"strings"
    "net/http"
    "database/sql"
    "encoding/json"
    "github.com/golang-jwt/jwt/v5"
)

var LocalID int = 0
var ExpressionsList []CalculationStore
var ExpressionRequestsList []CalculationRequest
var JWTTokensStore []JWT

func ValidateAndDecodeRequest(w http.ResponseWriter, r *http.Request) {
    decodedRequest, err := ReadRequestJson[CalculationRequest](w, r, "application/json")
    if err != nil {
        panic(err)
    }

    if SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore) {
        if SearchingIdInArray(decodedRequest.ID, ExpressionsList) {
            decodedRequest.ID = LocalID
            LocalID++
        } 

        // ExpressionsList = append(ExpressionsList, decodedRequest)
        // ExpressionRequestsList = append(ExpressionRequestsList, decodedRequest)

        ctx := context.TODO()

        db, err := sql.Open("sqlite3", "DataBase/Store.db")
        if err != nil {
            panic(err)
        }

        err = db.PingContext(ctx)
        if err != nil {
            panic(err)
        }    

        encodedExpression := EncodeExpression(CalculationStore{}, decodedRequest)

        if err = UpdateExpressionLine(ctx, db, encodedExpression, ""); err != nil {
            panic(err)
        }

        db.Close()

        // Sending result
        w.Write([]byte("Succesfully"))

        log.Println("\n\nSuccesfully") // Succesfuly data
    } else {
        // Sending result
        w.Write([]byte("JWT token is not founded"))

        log.Println("\n\nJWT token is not founded") // Succesfuly data
    }
}

func GettingExpressoinsList (w http.ResponseWriter, r *http.Request) {
    decodedRequest, err := ReadRequestJson[OnlyJWTTokens](w, r, "application/json")
    if err != nil {
        panic(err)
    }

    if SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore) {
        // Sending result
        w.Header().Set("Content-Type", "application/json") // Result Type
        message := "{\"expressions\": "+ string(UnificationExpressionsArray(w)) + " }"
        w.Write([]byte(message))

        // Sending result
        w.Write([]byte("Succesfully"))

        log.Println("\n\nSuccesfully") // Succesfuly data
    } else {
        // Sending result
        w.Write([]byte("JWT token is not founded"))

        log.Println("\n\nJWT token is not founded") // Succesfuly data
    }
}

func GettingExpressoinsListForID (w http.ResponseWriter, r *http.Request) {
    var result string
    var finded bool = false
    var neededExpressionIdIndex int

    decodedRequest, err := ReadRequestJson[CalculationRequest](w, r, "application/json")
    if err != nil {
        panic(err)
    }

    if SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore) {
        // Getting request like :id?ID=2
        log.Println(r.PathValue("id"))
        neededExpressionId, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
            return
        }

        array := GetExpressionsList(w)
        for i := 0; i <= len(array)-1; i++ {
            if array[i].ID == neededExpressionId {
                finded = true
                neededExpressionIdIndex = i
            }
        }

        if !finded { http.Error(w, "There is no such expression", 404)
        } else {
            result = UnificationExpression[CalculationStore](GetExpressionsList(w)[neededExpressionIdIndex], w)
        }

        // Sending result
        w.Header().Set("Content-Type", "application/json") // Result Type
        message := "{\"expression\": "+ result + " }"
        w.Write([]byte(message))

        // Sending result
        w.Write([]byte("Succesfully"))

        log.Println("\n\nSuccesfully") // Succesfuly data
    } else {
        // Sending result
        w.Write([]byte("JWT token is not founded"))

        log.Println("\n\nJWT token is not founded") // Succesfuly data
    }
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

        // Sending result
        w.Write([]byte("Succesfully"))
        
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

        // Sending result
        w.Write([]byte("Succesfully"))

        log.Println("\n\nSuccesfully") // Succesfuly data
    }
}

func UserRegist(w http.ResponseWriter, r *http.Request) {
    userRegist, err := ReadRequestJson[User](w, r, "application/json")
    if err != nil {
        panic(err)
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
        panic(err)
    }

    db.Close()

    // Sending result
    w.Write([]byte("Succesfully"))

    log.Println("\n\nSuccesfully") // Succesfuly data
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    userLogin, err := ReadRequestJson[User](w, r, "application/json")
    if err != nil {
        panic(err)
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
    
    if err = UserLogin(ctx, db, userLogin, w); err != nil {
        panic(err)
    }

    db.Close()

    const hmacSampleSecret = "secret"
	now := time.Now()
    jsonConsts := ReadJson()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": userLogin.Login,
		"nbf":  now.Add(time.Minute).Unix(),
		"exp":  now.Add(time.Duration(jsonConsts.ServerJWTTokenEXP) * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		panic(err)
	}

    newUser := JWT{UserId: userLogin.ID, Token: tokenString}
    JWTTokensStore = append(JWTTokensStore, newUser)

    // Sending result
    message := string(fmt.Sprintf("Your JWT token: %s\n\nSuccesfully", tokenString))
    w.Write([]byte(message))

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
    mux.HandleFunc("/api/v1/login", LoginUser)

    // Check database file is exists
    CheckDataBaseTables()

    jsonConsts := ReadJson()
    serverLocalization := fmt.Sprintf("%s:%s", jsonConsts.ServerIP, jsonConsts.ServerPort)
    log.Println(serverLocalization)

    err := http.ListenAndServe(serverLocalization, mux)
    log.Fatal(err)
}
