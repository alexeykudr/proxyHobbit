package repository

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	// "strconv"
	"time"
)

type ProxyPortItemSQLite struct {
	db *sql.DB
}

func NewProxyPortItem(db *sql.DB) *ProxyPortItemSQLite {
	return &ProxyPortItemSQLite{db: db}
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func randomInt(length int) int {
	return rand.Intn(length)
}

func (r *ProxyPortItemSQLite) GenerateSlug(routerPort int) (int, string, error) {
	var s_data int
	random_string := randomString(12)
	// UPDATE proxyPorts SET generatedUrl = 'tmp' WHERE id = 3

	// router_port_id := 11

	// a, _ := strconv.Atoi(routerPort)

	err := r.db.QueryRow("UPDATE proxyPorts SET generatedUrl = ? where router_id = ? RETURNING id",
		random_string, routerPort).Scan(&s_data)

	if err != nil {
		log.Printf("error in GenerateSlug with port %d", routerPort)
		return 0, "error", err
	}
	return s_data, random_string, nil
}

func (r *ProxyPortItemSQLite) GetIdBySlug(slug string) (int, error) {
	var s_data int
	err := r.db.QueryRow("select id from proxyPorts where generatedUrl = ?", slug).Scan(&s_data)
	if err != nil {
		log.Printf("error in GetIdBySlug with slug %s", slug)
		return 0, err
	}
	return s_data, nil
}

func (r *ProxyPortItemSQLite) UpdateReconnectInterval(portId int, intervalInMin string) (int, error) {
	var s_data int
	err := r.db.QueryRow("UPDATE proxyPorts SET interval = ? where router_id = ? RETURNING id", intervalInMin,
		portId).Scan(&s_data)

	if err != nil {
		log.Printf("error in UpdateReconectInterval with port %d , min interval %s", portId, intervalInMin)
		return 0, err
	}
	// плохой айди - вылетаем

	return s_data, err
}

func (r *ProxyPortItemSQLite) CreateSimpleUser(username string, password string) (string, error) {
	// 	INSERT INTO users (username, password)
	// VALUES ("testuser2", "1234");
	var username_data string
	// Calling Sprintf() function
	// s := fmt.Sprintf("%s is a %s Portal.\n", name, dept)

	stmt, err := r.db.Exec("INSERT INTO users(username, password) VALUES(?, ?) RETURNING id", username, password)
	fmt.Println(stmt)
	fmt.Println(err)
	// convert stmt to id
	// https://pkg.go.dev/database/sql#DB.Exec
	return username_data, nil
}
