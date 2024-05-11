package youtube

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	api "google.golang.org/api/youtube/v3"
	"net/url"
	"thumbs/internal/storage/sqlite"
)

type youtubeService struct {
	log     *zap.Logger
	api     *api.Service
	storage *sqlite.Storage
}

type Service interface {
	Download(ctx context.Context, url string) (string, error)
}

// Можно было сделать лучше, чтобы не привязываться к одной реализации на sqlite

func New(log *zap.Logger, storage *sqlite.Storage, youtubeKey string) Service {
	const op = "youtube.New"

	ctx := context.Background()
	service, err := api.NewService(ctx, option.WithAPIKey(youtubeKey))
	if err != nil {
		log.Panic(op, zap.Error(err))
	}

	return &youtubeService{
		log:     log,
		api:     service,
		storage: storage,
	}
}

func (s *youtubeService) Download(ctx context.Context, url string) (string, error) {
	const op = "youtube.Download"

	videoId, err := s.getVideoIdFromURL(url)
	if err != nil {
		s.log.Error(op, zap.Error(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if videoId == "" {
		return "", fmt.Errorf("%s: invalid url", op)
	}

	thumbnail, err := s.storage.GetThumbnail(ctx, videoId)
	if err != nil {
		s.log.Error(op, zap.Error(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if thumbnail != "" {
		return thumbnail, nil
	}

	call := s.api.Videos.List([]string{"snippet"}).Id(videoId)

	response, err := call.Do()
	if err != nil {
		s.log.Error(op, zap.Error(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("%s: video not found", op)
	}

	thumbnails := response.Items[0].Snippet.Thumbnails

	err = s.storage.SaveThumbnail(ctx, thumbnails.Maxres.Url, videoId)
	if err != nil {
		s.log.Error(op, zap.Error(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return thumbnails.Maxres.Url, nil
}

func (s *youtubeService) getVideoIdFromURL(videoURL string) (string, error) {
	const op = "youtube.getVideoIdFromURL"

	u, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	videoId := u.Query().Get("v")

	return videoId, nil
}
