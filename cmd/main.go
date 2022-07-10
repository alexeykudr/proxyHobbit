package main

import (
	"awesomeProject/pkg/handler"
	"awesomeProject/pkg/repository"
	"database/sql"
	// "fmt"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

// TODO
// init sqllite db
// start server
// create pkg/service,handler,repo


func NewSqlLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./testDB.db")
	if err != nil {
		log.Println("Error with open DB", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error with ping DB", err.Error())
		return nil, err
	}

	return db, nil
}

func main(){
	db, _ := NewSqlLiteDB()
	repo := repository.NewRepository(db)
	handler := handler.NewHandler(repo)	
	
	log.Fatal(http.ListenAndServe(":8080", handler.InitRoutes() ))
}
