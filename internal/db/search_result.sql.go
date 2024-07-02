// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: search_result.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSearchResult = `-- name: CreateSearchResult :exec
INSERT INTO searches_result ( ID, image_url, description, source, price, promotion, link, search_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type CreateSearchResultParams struct {
	ID          uuid.UUID `json:"id"`
	ImageUrl    string    `json:"image_url"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Price       float64   `json:"price"`
	Promotion   bool      `json:"promotion"`
	Link        string    `json:"link"`
	SearchID    uuid.UUID `json:"search_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) CreateSearchResult(ctx context.Context, arg CreateSearchResultParams) error {
	_, err := q.db.ExecContext(ctx, createSearchResult,
		arg.ID,
		arg.ImageUrl,
		arg.Description,
		arg.Source,
		arg.Price,
		arg.Promotion,
		arg.Link,
		arg.SearchID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteSearchResult = `-- name: DeleteSearchResult :exec
DELETE FROM searches_result
WHERE searches_result.id = $1
`

func (q *Queries) DeleteSearchResult(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSearchResult, id)
	return err
}

const getSearchResult = `-- name: GetSearchResult :one
SELECT id, image_url, description, source, price, promotion, link, search_id, created_at, updated_at
FROM searches_result
WHERE searches_result.id = $1
`

func (q *Queries) GetSearchResult(ctx context.Context, id uuid.UUID) (SearchesResult, error) {
	row := q.db.QueryRowContext(ctx, getSearchResult, id)
	var i SearchesResult
	err := row.Scan(
		&i.ID,
		&i.ImageUrl,
		&i.Description,
		&i.Source,
		&i.Price,
		&i.Promotion,
		&i.Link,
		&i.SearchID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSearchesResult = `-- name: GetSearchesResult :many
SELECT id, image_url, description, source, price, promotion, link, search_id, created_at, updated_at
FROM searches_result ORDER BY created_at DESC
`

func (q *Queries) GetSearchesResult(ctx context.Context) ([]SearchesResult, error) {
	rows, err := q.db.QueryContext(ctx, getSearchesResult)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchesResult
	for rows.Next() {
		var i SearchesResult
		if err := rows.Scan(
			&i.ID,
			&i.ImageUrl,
			&i.Description,
			&i.Source,
			&i.Price,
			&i.Promotion,
			&i.Link,
			&i.SearchID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
