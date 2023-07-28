package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-faster/errors"
	"time"

	"github.com/gh0st3e/NASA_API/internal/entity"
)

const (
	CtxTimeout = 5 * time.Second
)

type Store struct {
	db *sql.DB
}

func MewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveApod(ctx context.Context, apod entity.Apod) error {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	query := `INSERT INTO apod_data(date,explanation,hdurl,media_type,service_version,title,url,image)
						VALUES($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := s.db.ExecContext(ctx, query,
		apod.Date,
		apod.Explanation,
		apod.HdUrl,
		apod.MediaType,
		apod.ServiceVersion,
		apod.Title,
		apod.Url,
		apod.Image)

	return err
}

func (s *Store) RetrieveAPODByDate(ctx context.Context, date string) (*entity.Apod, error) {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	query := `SELECT date,explanation,hdurl,media_type,service_version,title,url,image
				FROM apod_data
				WHERE date=$1`

	apod := &entity.Apod{}

	err := s.db.QueryRowContext(ctx, query, date).Scan(
		&apod.Date,
		&apod.Explanation,
		&apod.HdUrl,
		&apod.MediaType,
		&apod.ServiceVersion,
		&apod.Title,
		&apod.Url,
		&apod.Image)

	if errors.Is(err, sql.ErrNoRows) {
		return apod, fmt.Errorf("no such apod")
	}

	return apod, err
}

func (s *Store) RetrieveAllApods(ctx context.Context) ([]entity.Apod, error) {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	query := `SELECT date,explanation,hdurl,media_type,service_version,title,url,image
				FROM apod_data`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var apods []entity.Apod

	for rows.Next() {
		apod := entity.Apod{}

		err := rows.Scan(
			&apod.Date,
			&apod.Explanation,
			&apod.HdUrl,
			&apod.Url,
			&apod.MediaType,
			&apod.ServiceVersion,
			&apod.Title,
			&apod.Url,
			&apod.Image)
		if err != nil {
			return nil, err
		}

		apods = append(apods, apod)
	}

	return apods, nil
}
