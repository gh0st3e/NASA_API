package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/pkg/errors"
)

const (
	CtxTimeout = 5 * time.Second
)

type ApodStore struct {
	db *sql.DB
}

func NewApodStore(db *sql.DB) *ApodStore {
	return &ApodStore{db: db}
}

func (a *ApodStore) SaveApod(ctx context.Context, apod *entity.Apod) error {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	query := `INSERT INTO apod_data(date,explanation,hdurl,media_type,service_version,title,url)
						VALUES($1,$2,$3,$4,$5,$6,$7)`

	_, err := a.db.ExecContext(ctx, query,
		apod.Date,
		apod.Explanation,
		apod.HdURL,
		apod.MediaType,
		apod.ServiceVersion,
		apod.Title,
		apod.URL)

	return err
}

func (a *ApodStore) RetrieveApodByDate(ctx context.Context, date string) (*entity.Apod, error) {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	query := `SELECT date,explanation,hdurl,media_type,service_version,title,url
				FROM apod_data
				WHERE date=$1`

	apod := &entity.Apod{}

	err := a.db.QueryRowContext(ctx, query, date).Scan(
		&apod.Date,
		&apod.Explanation,
		&apod.HdURL,
		&apod.MediaType,
		&apod.ServiceVersion,
		&apod.Title,
		&apod.URL)

	if errors.Is(err, sql.ErrNoRows) {
		return apod, fmt.Errorf("no such apod")
	}

	return apod, err
}

func (a *ApodStore) RetrieveAllApods(ctx context.Context) ([]entity.Apod, error) {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	query := `SELECT date,explanation,hdurl,media_type,service_version,title,url
				FROM apod_data`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var apods []entity.Apod

	for rows.Next() {
		apod := entity.Apod{}

		err := rows.Scan(
			&apod.Date,
			&apod.Explanation,
			&apod.HdURL,
			&apod.MediaType,
			&apod.ServiceVersion,
			&apod.Title,
			&apod.URL)
		if err != nil {
			return nil, err
		}

		apods = append(apods, apod)
	}

	return apods, nil
}
