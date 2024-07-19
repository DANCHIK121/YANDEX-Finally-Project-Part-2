package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// func WriteToDataBase(args [T], action_id int) [T any] {
// 	ctx := context.TODO()

// 	db, err := sql.Open("sqlite3", "DataBase/Store.db")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	switch (action_id) {
// 	case 0:
// 		if err = CreateTables(ctx, db); err != nil {
// 			panic(err)
// 		}
// 		break
// 	case 1:
// 		if err = RegistUser(ctx, db); err != nil {
// 			panic(err)
// 		}
// 		break
// 	}
// }

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
		searchUser = `
			SELECT id, login, password FROM users
		`
	)

	users, err := selectUsers(ctx, db)

	if err != nil {
		panic(err)
	}

	userFinded := false

	for i := 0; i <= len(users)-1; i++ {
		if users[i] == user {
			userFinded = true
			break
		}
	}

	if !userFinded {
		if _, err := db.ExecContext(ctx, registUser, user.Login, user.Password); err != nil {
			return err
		}
	} else {
		http.Error(w, "This user already registered", 505)
		return errors.New("This user already registered")
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

	// 	expressionsTable = `
	// CREATE TABLE IF NOT EXISTS expressions(
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	// 	expression TEXT NOT NULL,
	// 	user_id INTEGER NOT NULL,
	
	// 	FOREIGN KEY (user_id)  REFERENCES expressions (id)
	// );`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	// if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
	// 	return err
	// }

	return nil
}

func selectUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	const (
		searchUser = `
			SELECT login, password FROM users
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
		err := rows.Scan(&u.Login, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
