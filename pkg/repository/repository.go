package repository

import "database/sql"

type ProxyPorts interface{
	// return id of slugfield in db or error
	CreateSlugUrl(routerPort string)(int, string, error)
	GetIdBySlug(slug string) (int, error)
	UpdateReconnectInterval(portId string , minutes string)(int, error)
}

type Repository struct{
	ProxyPorts
}

func NewRepository(db *sql.DB) *Repository{
	return &Repository{
		ProxyPorts: NewProxyPortItem(db),
	}
}
