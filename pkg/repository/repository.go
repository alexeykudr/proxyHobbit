package repository

import (
	"database/sql"
)

type ProxyPorts interface {
	// return id of slugfield in db or error
	GenerateSlug(portId string) (int, string, error)
	GetIdBySlug(slug string) (int, error)
	UpdateReconnectInterval(portId int, minutes string) (int, error)
	CreateSimpleUser(id, username, password string) (error)
}

type Repository struct {
	ProxyPorts
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProxyPorts: NewProxyPortItem(db),
	}
}
