package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type repo interface {
	SaveShortUrl(ctx context.Context, shortUrl, origin string) error
	GetOrigUrl(ctx context.Context, shortUrl string) (string, error)
}
type ShortenerUrlService struct {
	shortenerUrl repo
}

func NewShortenerUrl(short repo) *ShortenerUrlService {
	return &ShortenerUrlService{shortenerUrl: short}
}
func (s *ShortenerUrlService) CreateShortUrl(ctx context.Context, origin string) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6
	shortUrl := make([]byte, codeLength)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < codeLength; i++ {
		shortUrl[i] = charset[r.Intn(len(charset))]
	}

	if err := s.shortenerUrl.SaveShortUrl(ctx, string(shortUrl), origin); err != nil {
		return "", err
	}
	return fmt.Sprintf("nozdrin-%s", string(shortUrl)), nil
}
func (s *ShortenerUrlService) GetOriginUrl(ctx context.Context, shortUrl string) (string, error) {
	origin, err := s.shortenerUrl.GetOrigUrl(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	return origin, nil
}
