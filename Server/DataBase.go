package main

import (
	"os"
	"fmt"
	"log"
	"errors"
	"context"
	"net/http"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CheckDataBaseTables() {
	if _, err := os.Stat("/DataBase/Store.db"); errors.Is(err, os.ErrNotExist) {
		ctx := context.TODO()

		db, err := sql.Open("sqlite3", "DataBase/Store.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = db.PingContext(ctx)
		if err != nil {
			panic(err)
		}
		
		if err = CreateTables(ctx, db); err != nil {
			panic(err)
		}
	} else {
		log.Println("DataBase file is finded seccfully")
	}
}

func RegistUser(ctx context.Context, db *sql.DB, user User, w http.ResponseWriter) error {
	const (
		registUser = `
			INSERT INTO users (login, password) values ($1, $2)
		`
		creatingNewExpressionsCell = `
			INSERT INTO expressions (solved_expressions, expression, user_id) values ($1, $2, $3)
		`
	)

	users, err := SelectUsers(ctx, db)

	if err != nil {
		panic(err)
	}

	userFinded := false

	for i := 0; i <= len(users)-1; i++ {
		user.ID = users[i].ID
		if users[i] == user {
			userFinded = true
			break
		}
	}

	if !userFinded {
		if _, err := db.ExecContext(ctx, registUser, user.Login, user.Password); err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, creatingNewExpressionsCell, "", "", user.ID+1); err != nil {
			return err
		}
	} else {
		// Sending result
		w.Write([]byte("This user already registered"))
	}

	return nil
}

func UserLogin(ctx context.Context, db *sql.DB, user User, w http.ResponseWriter) error {
	const (
		searchUser = `
			SELECT id, login, password FROM users
		`
	)

	users, err := SelectUsers(ctx, db)

	if err != nil {
		panic(err)
	}

	userFinded := false

	for i := 0; i <= len(users)-1; i++ {
		user.ID = users[i].ID
		if users[i] == user {
			userFinded = true
			break
		}
	}

	if !userFinded {
		// Sending result
		w.Write([]byte("User is not founded"))
	}

	return nil
}

func UpdateExpressionLine(ctx context.Context, db *sql.DB, user_id int, expression string, solved_expression string) error {
	updateSolvedExpressionLine := fmt.Sprintf(`
		UPDATE expressions SET solved_expressions = solved_expressions || "%s" WHERE user_id = %s;
	`, solved_expression, fmt.Sprint(user_id))

	updateExpressionLine := fmt.Sprintf(`
		UPDATE expressions SET expression = expression || "%s" WHERE user_id = %s;
	`, expression, fmt.Sprint(user_id))

	if solved_expression == "" {
		if _, err := db.ExecContext(ctx, updateExpressionLine); err != nil {
			return err
		}
	} else {
		if _, err := db.ExecContext(ctx, updateSolvedExpressionLine); err != nil {
			return err
		}
	}

	return nil
}

func CreateTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			login TEXT,
			password TEXT NOT NULL
		);`

		expressionsTable = `
		CREATE TABLE IF NOT EXISTS expressions(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			expression         TEXT NOT NULL,
			solved_expressions TEXT NOT NULL,
			user_id            INTEGER NOT NULL,
		
			FOREIGN KEY (user_id)  REFERENCES expressions (id)
		);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}

func SelectUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	const (
		searchUser = `
			SELECT id, login, password FROM users
		`
	)

	var users []User
	rows, err := db.QueryContext(ctx, searchUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func SelectSolvedExpression(ctx context.Context, db *sql.DB, user_id int) (string, error) {
	searchExpression := fmt.Sprintf(`
		SELECT solved_expressions FROM expressions WHERE user_id = %s
	`, fmt.Sprint(user_id))

	var expression Expression
	rows, err := db.QueryContext(ctx, searchExpression)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&expression.Expression)
		if err != nil {
			return "", err
		}
	}

	return expression.Expression, nil
}

func SelectExpression(ctx context.Context, db *sql.DB, user_id int) (string, error) {
	searchExpression := fmt.Sprintf(`
		SELECT expression FROM expressions WHERE user_id = %s
	`, fmt.Sprint(user_id))

	var expression Expression
	rows, err := db.QueryContext(ctx, searchExpression)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&expression.Expression)
		if err != nil {
			return "", err
		}
	}

	return expression.Expression, nil
}

func SelectPastUserID(ctx context.Context, db *sql.DB) (int, error) {
	const (
		searchPastUserID = `
			SELECT user_id FROM expressions
		`
	)

	var users []PastUserID
	rows, err := db.QueryContext(ctx, searchPastUserID)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		u := PastUserID{}
		err := rows.Scan(&u.ID)
		if err != nil {
			return -1, err
		}
		users = append(users, u)
	}

	return users[len(users)-1].ID, nil
}

func SelectUserForLogin(ctx context.Context, db *sql.DB, login string) (int, error) {
	searchPastUserID := fmt.Sprintf(`
		SELECT id FROM users WHERE login = "%s"
	`, login)

	var users []User
	rows, err := db.QueryContext(ctx, searchPastUserID)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID)
		if err != nil {
			return -1, err
		}
		users = append(users, u)
	}

	return users[len(users)-1].ID, nil
}
