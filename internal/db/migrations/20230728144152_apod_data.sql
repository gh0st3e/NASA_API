-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS apod_data (
       id SERIAL NOT NULL,
       date DATE NOT NULL,
       explanation TEXT NOT NULL,
       hdurl TEXT NOT NULL,
       media_type TEXT NOT NULL,
       service_version TEXT NOT NULL,
       title TEXT NOT NULL,
       url TEXT NOT NULL,
       image BYTEA NOT NULL,
       CONSTRAINT apod_data_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE TABLE IF EXISTS apod_data;
-- +goose StatementEnd
