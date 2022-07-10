package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type ProxyServerContainer struct {
}

func nativeExec() string {
	cmd, err := exec.Command("/bin/sh", "/Users/mac/Desktop/job.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string(cmd)
	return output
}
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func rebootProxyRouter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "starting reboot, %s!\n", ps.ByName("url"))
	// get change url by select generated url
	slug := ps.ByName("url")
	params := r.URL.Query()
	//myParam := params["my-query-param"]
	fmt.Println(slug, params)
	// select changeUrl from testDB where generatedUrl like 'slug'
}

//func Hello(w http.ResponseWriter, r *http.Request) {
//	params := httprouter.ParamsFromContext(r.Context())
//params := r.Context().Value(httprouter.ParamsKey)
//	fmt.Fprintf(w, "hello, %s!\n", params.ByName("name"))
//}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./testDB.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	statement, error := db.Prepare("CREATE TABLE IF NOT EXISTS proxyPorts (id INTEGER PRIMARY KEY AUTOINCREMENT , changeUrl TEXT, generatedUrl TEXT)")
	statement.Exec()
	return db, error
}

type proxyPortDbObject struct {
	id           int
	changeUrl    string
	generatedUrl string
}

func testScan(db *sql.DB) {
	var s_data string
	err := db.QueryRow("select changeUrl from proxyPorts where generatedUrl = ?", "zxcc").Scan(&s_data)
	if err != nil {
		log.Println("error in testScan", err.Error())
	}

	//log.Println(err2)
	log.Println(s_data)
}

func execProxyCommand() {
	ctx, cancel := context.WithTimeout(context.Background(), 5500*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
		fmt.Println("err!")
	}
	fmt.Println("not error")
}
func main() {
	// fmt.Printf(nativeExec())
	db, err := initDB()
	if err != nil {
		log.Fatal("err to connect to ddb")
	}

	//rows, _ := db.Query("select changeUrl from proxyPorts")
	//
	//var s string
	//for rows.Next() {
	//	rows.Scan(&s)
	//	//fmt.Println(strconv.Itoa())
	//}
	//fmt.Println(s)
	testScan(db)
	//router := httprouter.New()
	//router.GET("/", Index)
	//router.GET("/reboot/:url", rebootProxyRouter)
	//
	//log.Fatal(http.ListenAndServe(":8080", router))
}
