package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveThumbnail(ctx context.Context, thumbnailUrl string, videoId string) error {
	const op = "storage.sqlite.SaveThumbnail"

	stmt, err := s.db.Prepare("INSERT INTO thumbs (url, video_id) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, thumbnailUrl, videoId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetThumbnail(ctx context.Context, videoID string) (string, error) {
	const op = "storage.sqlite.GetThumbnail"

	stmt, err := s.db.Prepare("SELECT url from thumbs WHERE video_id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, videoID)

	var url string

	err = row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
