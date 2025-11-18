package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Shortener struct {
	db *sql.DB
}

func NewShortener(db *sql.DB) *Shortener {
	return &Shortener{db: db}
}

func (s *Shortener) SaveShortUrl(ctx context.Context, shortUrl, origin string) error {
	if origin == "" {
		return fmt.Errorf("origin link dont given")
	}
	_, err := s.db.ExecContext(ctx, "insert into links(short_url, origin_url) values($1, $2)", shortUrl, origin)
	if err != nil {
		return err
	}
	return nil
}
func (s *Shortener) GetOrigUrl(ctx context.Context, shortUrl string) (string, error) {
	var data struct {
		origin string
		status bool
	}
	if shortUrl == "" {
		return "", fmt.Errorf("short link dont given")
	}
	cleanShortUrl := strings.TrimPrefix(shortUrl, "nozdrin-")
	row := s.db.QueryRowContext(ctx, "select origin_url, status from links where short_url=$1", cleanShortUrl)

	if err := row.Scan(&data.origin, &data.status); err != nil {
		return "", err
	}
	if !data.status {
		return "", fmt.Errorf("short link inactive")
	}
	return data.origin, nil

}
