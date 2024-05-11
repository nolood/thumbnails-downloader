package grpc

import (
	thumbsv1 "cli-thumbs/cli-thumbs"
	"context"
	"fmt"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	api thumbsv1.DownloaderClient
}

// Было бы не плохо добавить логгер

func New(ctx context.Context, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{api: thumbsv1.NewDownloaderClient(cc)}, nil
}

func (c *Client) Download(ctx context.Context, url string) (string, error) {
	const op = "grpc.Download"

	resp, err := c.api.Download(ctx, &thumbsv1.DownloadRequest{Url: url})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetUrl(), nil
}
