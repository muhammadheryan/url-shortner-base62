package url_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	appurl "github.com/muhammadheryan/url-shortner-base62/application/url"
	"github.com/muhammadheryan/url-shortner-base62/constant"
	urlmocks "github.com/muhammadheryan/url-shortner-base62/mocks/repository/url"
	"github.com/muhammadheryan/url-shortner-base62/model"
	cerr "github.com/muhammadheryan/url-shortner-base62/utils/errors"
	"github.com/stretchr/testify/mock"
)

func TestURLApp_CreateURLShortner(t *testing.T) {
	type fields struct {
		urlRepo *urlmocks.URLRepository
	}
	type args struct {
		ctx context.Context
		req *model.CreateURLShortnerRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockCall func(f fields)
		want     *model.GetURLResponse
		wantErr  bool
	}{
		{
			name: "success: normalize URL, create then update",
			fields: fields{
				urlRepo: urlmocks.NewURLRepository(t),
			},
			args: args{
				ctx: context.Background(),
				req: &model.CreateURLShortnerRequest{OriginalURL: "example.com"},
			},
			mockCall: func(f fields) {
				f.urlRepo.
					On("Create", mock.Anything, mock.MatchedBy(func(ent *model.URLEntity) bool {
						return ent.OriginalURL == "https://example.com"
					})).
					Return(&model.URLEntity{
						ID:          1,
						OriginalURL: "https://example.com",
						CreatedAt:   time.Now(),
					}, nil).
					Once()

				// ID=1 -> shortURL "00001" (minLength=5)
				f.urlRepo.
					On("Update", mock.Anything, mock.MatchedBy(func(ent *model.URLEntity) bool {
						return ent.ID == 1 && ent.ShortURL == "00001"
					})).
					Return(&model.URLEntity{
						ID:          1,
						ShortURL:    "00001",
						OriginalURL: "https://example.com",
						CreatedAt:   time.Now(),
						UpdatedAt:   nil,
					}, nil).
					Once()
			},
			want: &model.GetURLResponse{
				ShortURL:    "00001",
				OriginalURL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "error: repository Create returns error -> ErrInternal",
			fields: fields{
				urlRepo: urlmocks.NewURLRepository(t),
			},
			args: args{
				ctx: context.Background(),
				req: &model.CreateURLShortnerRequest{OriginalURL: "foo.com"},
			},
			mockCall: func(f fields) {
				f.urlRepo.
					On("Create", mock.Anything, mock.AnythingOfType("*model.URLEntity")).
					Return(nil, errors.New("db down")).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error: repository Update returns error -> ErrInternal",
			fields: fields{
				urlRepo: urlmocks.NewURLRepository(t),
			},
			args: args{
				ctx: context.Background(),
				req: &model.CreateURLShortnerRequest{OriginalURL: "bar.com"},
			},
			mockCall: func(f fields) {
				f.urlRepo.
					On("Create", mock.Anything, mock.AnythingOfType("*model.URLEntity")).
					Return(&model.URLEntity{
						ID:          10,
						OriginalURL: "https://bar.com",
						CreatedAt:   time.Now(),
					}, nil).
					Once()

				f.urlRepo.
					On("Update", mock.Anything, mock.AnythingOfType("*model.URLEntity")).
					Return(nil, errors.New("update failed")).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockCall != nil {
				ttFields := tt.fields
				// Ensure mock expectations are set before calling app
				tt.mockCall(ttFields)
			}
			app := appurl.NewURLApplication(tt.fields.urlRepo)

			got, err := app.CreateURLShortner(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateURLShortner() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				var ce cerr.CustomError
				if !errors.As(err, &ce) {
					t.Fatalf("error type = %T, want CustomError", err)
				}
				if ce.ErrorCode() != constant.ErrorTypeCode[constant.ErrInternal] {
					t.Fatalf("error code = %s, want %s", ce.ErrorCode(), constant.ErrorTypeCode[constant.ErrInternal])
				}
				return
			}

			if got.ShortURL != tt.want.ShortURL || got.OriginalURL != tt.want.OriginalURL {
				t.Fatalf("CreateURLShortner() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestURLApp_GetURLByShortURL(t *testing.T) {
	type fields struct {
		urlRepo *urlmocks.URLRepository
	}
	type args struct {
		ctx      context.Context
		shortURL string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockCall func(f fields)
		want     *model.GetURLResponse
		wantErr  bool
	}{
		{
			name: "success: found entity",
			fields: fields{
				urlRepo: urlmocks.NewURLRepository(t),
			},
			args: args{
				ctx:      context.Background(),
				shortURL: "0000Z",
			},
			mockCall: func(f fields) {
				now := time.Now()
				f.urlRepo.
					On("Get", mock.Anything, &model.URLFilter{ShortURL: "0000Z"}).
					Return(&model.URLEntity{
						ID:          99,
						UserID:      0,
						ShortURL:    "0000Z",
						OriginalURL: "https://golang.org",
						CreatedAt:   now,
						UpdatedAt:   nil,
					}, nil).
					Once()
			},
			want: &model.GetURLResponse{
				ShortURL:    "0000Z",
				OriginalURL: "https://golang.org",
			},
			wantErr: false,
		},
		{
			name: "not found: repo returns (nil, nil) -> ErrNotFound",
			fields: fields{
				urlRepo: urlmocks.NewURLRepository(t),
			},
			args: args{
				ctx:      context.Background(),
				shortURL: "xxxxx",
			},
			mockCall: func(f fields) {
				f.urlRepo.
					On("Get", mock.Anything, &model.URLFilter{ShortURL: "xxxxx"}).
					Return(nil, nil).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error: repository returns error -> ErrInternal",
			fields: fields{
				urlRepo: urlmocks.NewURLRepository(t),
			},
			args: args{
				ctx:      context.Background(),
				shortURL: "errxx",
			},
			mockCall: func(f fields) {
				f.urlRepo.
					On("Get", mock.Anything, &model.URLFilter{ShortURL: "errxx"}).
					Return(nil, errors.New("query failed")).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockCall != nil {
				ttFields := tt.fields
				// Ensure mock expectations are set before calling app
				tt.mockCall(ttFields)
			}
			app := appurl.NewURLApplication(tt.fields.urlRepo)

			got, err := app.GetURLByShortURL(tt.args.ctx, tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetURLByShortURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				var ce cerr.CustomError
				if !errors.As(err, &ce) {
					t.Fatalf("error type = %T, want CustomError", err)
				}
				if tt.args.shortURL == "xxxxx" && ce.ErrorCode() != constant.ErrorTypeCode[constant.ErrNotFound] {
					t.Fatalf("error code = %s, want %s", ce.ErrorCode(), constant.ErrorTypeCode[constant.ErrNotFound])
				}
				if tt.args.shortURL == "errxx" && ce.ErrorCode() != constant.ErrorTypeCode[constant.ErrInternal] {
					t.Fatalf("error code = %s, want %s", ce.ErrorCode(), constant.ErrorTypeCode[constant.ErrInternal])
				}
				return
			}

			want := &model.GetURLResponse{
				ShortURL:    tt.want.ShortURL,
				OriginalURL: tt.want.OriginalURL,
			}
			gotComparable := &model.GetURLResponse{
				ShortURL:    got.ShortURL,
				OriginalURL: got.OriginalURL,
			}
			if !reflect.DeepEqual(gotComparable, want) {
				t.Fatalf("GetURLByShortURL() = %+v, want %+v", gotComparable, want)
			}
		})
	}
}
