package repository

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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

func (r *ProxyPortItemSQLite) CreateSlugUrl(routerPort string) (int, string, error) {
	var s_data int
	random_string := randomString(12)
	// UPDATE proxyPorts SET generatedUrl = 'tmp' WHERE id = 3
	err := r.db.QueryRow("UPDATE proxyPorts SET generatedUrl = ? where id = ? RETURNING id", random_string,
		routerPort).Scan(&s_data)

	if err != nil {
		log.Printf("error in CreateSlugUrl with port %s", routerPort)
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

func (r *ProxyPortItemSQLite) UpdateReconnectInterval(portId string, intervalInMin string) (int, error) {
	var s_data int
	err := r.db.QueryRow("UPDATE proxyPorts SET interval = ? where id = ? RETURNING id", intervalInMin,
		portId).Scan(&s_data)

	if err != nil {
		log.Printf("error in UpdateReconectInterval with port %s , min interval %s", portId, intervalInMin)
		// return 0, err
	}
	// плохой айди - вылетаем

	return s_data, err
}
