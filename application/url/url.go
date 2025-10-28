package url

import (
	"context"
	"log"
	"strings"

	"github.com/muhammadheryan/url-shortner-base62/constant"
	"github.com/muhammadheryan/url-shortner-base62/model"
	"github.com/muhammadheryan/url-shortner-base62/repository/url"
	"github.com/muhammadheryan/url-shortner-base62/utils/errors"
)

type URLAppImpl struct {
	URLRepository url.URLRepository
}

type URLApp interface {
	CreateURLShortner(ctx context.Context, req *model.CreateURLShortnerRequest) (*model.GetURLResponse, error)
	GetURLByShortURL(ctx context.Context, shortURL string) (*model.GetURLResponse, error)
}

func NewURLApplication(URLRepository url.URLRepository) URLApp {
	return &URLAppImpl{
		URLRepository: URLRepository,
	}
}

func (u *URLAppImpl) CreateURLShortner(ctx context.Context, req *model.CreateURLShortnerRequest) (*model.GetURLResponse, error) {
	// check http or https
	if !strings.HasPrefix(req.OriginalURL, "http://") && !strings.HasPrefix(req.OriginalURL, "https://") {
		req.OriginalURL = "https://" + req.OriginalURL
	}

	// Create in database to get ID
	createdURL, err := u.URLRepository.Create(ctx, &model.URLEntity{
		UserID:      0, // You might want to get this from context or auth
		OriginalURL: req.OriginalURL,
	})
	if err != nil {
		log.Println("[CreateURLShortner] err Create", err)
		return nil, errors.SetCustomError(constant.ErrInternal)
	}

	// Generate short URL from ID
	shortURL := createBase62Converter(createdURL.ID)

	// Update the URL entity with short URL
	createdURL.ShortURL = shortURL
	updatedURL, err := u.URLRepository.Update(ctx, createdURL)
	if err != nil {
		log.Println("[CreateURLShortner] err Update", err)
		return nil, errors.SetCustomError(constant.ErrInternal)
	}

	// Return response
	return &model.GetURLResponse{
		ShortURL:    updatedURL.ShortURL,
		OriginalURL: updatedURL.OriginalURL,
		CreatedAt:   updatedURL.CreatedAt,
		UpdatedAt:   updatedURL.UpdatedAt,
	}, nil
}

func (u *URLAppImpl) GetURLByShortURL(ctx context.Context, shortURL string) (*model.GetURLResponse, error) {
	urlEntity, err := u.URLRepository.Get(ctx, &model.URLFilter{
		ShortURL: shortURL,
	})
	if err != nil {
		log.Println("[GetURLByShortURL] err Get", err)
		return nil, errors.SetCustomError(constant.ErrInternal)
	}

	if urlEntity == nil {
		return nil, errors.SetCustomError(constant.ErrNotFound)
	}

	// Return response
	return &model.GetURLResponse{
		ShortURL:    urlEntity.ShortURL,
		OriginalURL: urlEntity.OriginalURL,
		CreatedAt:   urlEntity.CreatedAt,
		UpdatedAt:   urlEntity.UpdatedAt,
	}, nil
}

func createBase62Converter(id uint64) (shortURL string) {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const minLength = 5

	if id == 0 {
		return strings.Repeat(string(base62Chars[0]), minLength)
	}

	var result []byte
	for id > 0 {
		result = append(result, base62Chars[id%62])
		id /= 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	shortURL = string(result)
	if len(shortURL) < minLength {
		padding := strings.Repeat(string(base62Chars[0]), minLength-len(shortURL))
		shortURL = padding + shortURL
	}

	return shortURL
}
