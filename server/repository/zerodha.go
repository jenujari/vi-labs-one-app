package repository

import (
	"errors"
	"server/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	ZERODHA_AUTH_TOKEN = "ZERODHA_AUTH_TOKEN"
)

func (r *Repository) SaveZerodhaAuth(token string) error {
	in_ := sqlc.UpsertCacheParams{
		Key:     ZERODHA_AUTH_TOKEN,
		Value:   pgtype.Text{String: token, Valid: true},
		Created: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	return r.q.UpsertCache(r.ctx, in_)
}

func (r *Repository) GetZerodhaAuth() (string, error) {
	now := time.Now()

	cache, err := r.q.GetCache(r.ctx, ZERODHA_AUTH_TOKEN)
	if err != nil {
		return "", err
	}

	diff := now.Sub(cache.Created.Time)

	if diff.Hours() > 8 {
		return "", errors.New("cache expired")
	}

	return cache.Value.String, nil
}
