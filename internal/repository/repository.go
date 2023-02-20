package repository

import (
	"psn_discount_bot/internal/connector"
	"psn_discount_bot/internal/model"
	"psn_discount_bot/internal/utils/driver/postgres"
)

type Repo struct {
	pg postgres.IDB
}

func New(conn connector.IConnector) *Repo {
	repo := &Repo{
		pg: conn.Pg(),
	}

	//todo
	repo.pg.DB().RegisterModel((*model.UsersGames)(nil))

	return repo
}
