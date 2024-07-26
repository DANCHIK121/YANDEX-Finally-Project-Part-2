package main

import (
	"fmt"
	"log"
    "time"
    "context"
	"strconv"
    "net/http"
    "database/sql"
    "github.com/golang-jwt/jwt/v5"
)

var LocalID int = 0
var LocalIDForDataBase int = 0
var JWTTokensStore []JWT
var UseGRPC bool

func ValidateAndDecodeRequest(w http.ResponseWriter, r *http.Request) {
    decodedRequest, err := ReadRequestJson[CalculationRequest](w, r, "application/json", false)
    if err != nil {
        panic(err)
    }

    result, index := SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore)
    if result {
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

        if err = UpdateExpressionLine(ctx, db, JWTTokensStore[index].UserId, encodedExpression, ""); err != nil {
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
    decodedRequest, err := ReadRequestJson[OnlyJWTTokens](w, r, "application/json", false)
    if err != nil {
        panic(err)
    }

    result, index := SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore)
    if result {
        // Sending result
        w.Header().Set("Content-Type", "application/json") // Result Type
        message := "{\"expressions\": "+ string(UnificationExpressionsArray(w, JWTTokensStore[index].UserId)) + " }"
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

    decodedRequest, err := ReadRequestJson[CalculationRequest](w, r, "application/json", false)
    if err != nil {
        panic(err)
    }

    searchResult, index := SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore)
    if searchResult {
        // Getting request like :id?ID=2
        log.Println(r.PathValue("id"))
        neededExpressionId, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            msg := "Something went wrong"
            http.Error(w, msg, 500)
            return
        }

        array := GetSolvedExpressionsList(JWTTokensStore[index].UserId)
        for i := 0; i <= len(array)-1; i++ {
            if array[i].ID == neededExpressionId {
                finded = true
                neededExpressionIdIndex = i
            }
        }

        if !finded { http.Error(w, "There is no such expression", 404)
        } else {
            result = UnificationExpression[CalculationStore](GetSolvedExpressionsList(JWTTokensStore[index].UserId)[neededExpressionIdIndex], w)
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
    if ct == "application/json" {
        decodedRequest, err := ReadRequestJson[TaskResult](w, r, "application/json", false)
        if err != nil {
            panic(err)
        }
        if UseGRPC {
            GRPCmain(decodedRequest.JWTToken, w, decodedRequest)
        } else {
            result, index := SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore)
            if result {
                result, err := SearchingTaskResult(decodedRequest, GetSolvedExpressionsList(JWTTokensStore[index].UserId), w)

                if err != nil {
                    msg := "Something went wrong"
                    http.Error(w, msg, 500)
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

                encodeResult := EncodeExpression(CalculationStore{ID: JWTTokensStore[index].UserId, Result: result, Status: "Complited"} , CalculationRequest{})

                err = UpdateExpressionLine(ctx, db, JWTTokensStore[index].UserId, "", encodeResult)
                if err != nil {
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
    } else {
        decodedRequest, err := ReadRequestJson[OnlyJWTTokens](w, r, "", true)
        if err != nil {
            panic(err)
        }

        result, index := SearchingTokenInJWTStore(decodedRequest.JWTToken, JWTTokensStore)
        if result {
            if len(GetExpressionsList(JWTTokensStore[index].UserId)) <= 0 {
                msg := "Something went wrong"
                http.Error(w, msg, 500)
            }

            task.ID = GetExpressionsList(JWTTokensStore[index].UserId)[len(GetExpressionsList(JWTTokensStore[index].UserId))-1].ID

            decodedExpression := PostFixDecoding(GetExpressionsList(JWTTokensStore[index].UserId)[len(GetExpressionsList(JWTTokensStore[index].UserId))-2])

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
            w.Write([]byte("\n\nSuccesfully"))

            log.Println("\n\nSuccesfully") // Succesfuly data
        } else {
            // Sending result
            w.Write([]byte("JWT token is not founded"))

            log.Println("\n\nJWT token is not founded") // Succesfuly data
        }
    }
}

func UserRegist(w http.ResponseWriter, r *http.Request) {
    userRegist, err := ReadRequestJson[User](w, r, "application/json", false)
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
    userLogin, err := ReadRequestJson[User](w, r, "application/json", false)
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


    ctx = context.TODO()

    db, err = sql.Open("sqlite3", "DataBase/Store.db")
    if err != nil {
        panic(err)
    }

    err = db.PingContext(ctx)
    if err != nil {
        panic(err)
    }
    
    IDForDataBase, err := SelectUserForLogin(ctx, db, userLogin.Login)
    if err != nil {
        panic(err)
    }

    db.Close()

    newUser := JWT{UserId: IDForDataBase, Token: tokenString}
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

    UseGRPC = false
    jsonConsts := ReadJson()
    serverLocalization := fmt.Sprintf("%s:%s", jsonConsts.ServerIP, jsonConsts.ServerPort)
    log.Println(serverLocalization)

    err := http.ListenAndServe(serverLocalization, mux)
    log.Fatal(err)
}
