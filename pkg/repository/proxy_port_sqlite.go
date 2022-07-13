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

func (r *ProxyPortItemSQLite) GenerateSlug(routerPort string) (int, string, error) {
	var s_data int
	random_string := randomString(12)
	// router_port_id := 11

	// a, _ := strconv.Atoi(routerPort)

	err := r.db.QueryRow("UPDATE proxyPorts SET generatedUrl = ? where router_id = ? RETURNING id",
		random_string, routerPort).Scan(&s_data)

	if err != nil {
		log.Printf("error in GenerateSlug with port %s", routerPort)
		return 0, "error", err
	}
	return s_data, random_string, nil
}

func (r *ProxyPortItemSQLite) GetIdBySlug(slug string) (int, error) {
	var s_data int
	err := r.db.QueryRow("select router_id from proxyPorts where generatedUrl = ?", slug).Scan(&s_data)
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

func (r *ProxyPortItemSQLite) CreateSimpleUser(username string, password string) (error) {
	stmt, err := r.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	if _, err := stmt.Exec(username, password); err != nil {
		log.Println(err.Error())
	}
	// https://pkg.go.dev/database/sql#DB.Exec
	return nil
}
