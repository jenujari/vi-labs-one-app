package repository

import (
	"context"
	"server/config"
	"server/helpers"
	"server/sqlc"
)

var repo *Repository

type Repository struct {
	q *sqlc.Queries
	ctx context.Context
}

func InitRepository() {
	pc := helpers.GetMainProcess()
	dbc := config.GetDBC()
	q := sqlc.New(dbc)

	repo = &Repository{
		q:   q,
		ctx: pc.CTX,
	}

	pc.SetContextValue(config.REPO_KEY, repo)
}

func GetRepository() *Repository {
	return repo
}	