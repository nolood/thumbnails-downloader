package thumbs

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"thumbs/internal/services/youtube"
	thumbsv1 "thumbs/thumbs"
)

type serverAPI struct {
	thumbsv1.UnimplementedDownloaderServer
	youtubeService youtube.Service
}

func Register(gRPC *grpc.Server, youtubeService youtube.Service) {
	thumbsv1.RegisterDownloaderServer(gRPC, &serverAPI{youtubeService: youtubeService})
}

func (s *serverAPI) Download(ctx context.Context, req *thumbsv1.DownloadRequest) (*thumbsv1.DownloadResponse, error) {
	const op = "grpc.thumbs.server.Download"

	// По-хорошему проверять, что пришла ссылка, а не просто текст
	url := req.GetUrl()
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid url")
	}

	url, err := s.youtubeService.Download(ctx, url)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to download thumbnail")
	}

	return &thumbsv1.DownloadResponse{Url: url}, nil

}
