package url

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/muhammadheryan/url-shortner-base62/model"
)

type SQL struct {
	conn *sqlx.DB
}

type URLRepository interface {
	Create(ctx context.Context, req *model.URLEntity) (*model.URLEntity, error)
	Update(ctx context.Context, req *model.URLEntity) (*model.URLEntity, error)
	Get(ctx context.Context, filter *model.URLFilter) (*model.URLEntity, error)
}

func NewURLRepository(conn *sqlx.DB) URLRepository {
	return &SQL{conn: conn}
}

const (
	insertURLQuery = `INSERT INTO url (user_id, original_url, created_at) VALUES (?, ?, NOW())`
	updateURLQuery = `UPDATE url SET short_url = ?, original_url = ?, updated_at = NOW() WHERE id = ?`
	getURLBase     = `SELECT id, user_id, short_url, original_url, created_at, updated_at FROM url WHERE true`
)

func (s *SQL) Create(ctx context.Context, data *model.URLEntity) (*model.URLEntity, error) {
	result, err := s.conn.ExecContext(ctx, insertURLQuery, data.UserID, data.OriginalURL)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	data.ID = uint64(lastID)

	return data, nil
}

func (s *SQL) Update(ctx context.Context, data *model.URLEntity) (*model.URLEntity, error) {
	_, err := s.conn.ExecContext(ctx, updateURLQuery, data.ShortURL, data.OriginalURL, data.ID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *SQL) Get(ctx context.Context, filter *model.URLFilter) (*model.URLEntity, error) {
	query := getURLBase
	args := make([]any, 0, 2)

	if filter.ID != 0 {
		query += " AND id = ?"
		args = append(args, filter.ID)
	}
	if filter.ShortURL != "" {
		query += " AND short_url = ?"
		args = append(args, filter.ShortURL)
	}

	var entity model.URLEntity
	if err := s.conn.QueryRowxContext(ctx, query, args...).StructScan(&entity); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
