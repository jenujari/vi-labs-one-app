package repository

import (
	"errors"
	"server/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	ZERODHA_ACCESS_TOKEN = "ZERODHA_ACCESS_TOKEN"
)

func (r *Repository) SaveZerodhaAccessToken(token string) error {
	in := sqlc.UpsertCacheParams{
		Key:     ZERODHA_ACCESS_TOKEN,
		Value:   pgtype.Text{String: token, Valid: true},
		Created: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	return r.q.UpsertCache(r.ctx, in)
}

func (r *Repository) GetZerodhaAccessToken() (string, error) {
	now := time.Now()

	cache, err := r.q.GetCache(r.ctx, ZERODHA_ACCESS_TOKEN)
	if err != nil {
		return "", err
	}

	diff := now.Sub(cache.Created.Time)

	if diff.Hours() > 8 {
		return "", errors.New("cache expired")
	}

	return cache.Value.String, nil
}
