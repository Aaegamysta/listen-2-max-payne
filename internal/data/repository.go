package data

import (
	"context"
	"fmt"

	"github.com/aaegamysta/listen-2-max-payne/internal/twitter"
	"go.uber.org/zap"
)

type Interface interface {
	CreateTablesIfNotExists(ctx context.Context) error
	BatchInsertExcerpts(excerpts []Excerpt) ([]Excerpt, error)
	GetRandomExcerpt() (Excerpt, error)
	InsertSuccessfullTweetResponse(res twitter.SucessfullTweetResponse) (twitter.Tweet, error)
	InserrtUnsuccessflulTweetResponse(tweet twitter.Tweet) (twitter.Tweet, error)
}

type Impl struct {
	logger           *zap.SugaredLogger
	connectionString string
}

func New(cfg Config, logger *zap.SugaredLogger) Interface {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	impl := &Impl{
		connectionString: connStr,
		logger:           logger,
	}
	return impl
}

func (i *Impl) CreateTablesIfNotExists(ctx context.Context) error {
	panic("unimplemented")
}

func (i *Impl) BatchInsertExcerpts(excerpts []Excerpt) ([]Excerpt, error) {
	panic("unimplemented")
}

func (i *Impl) GetRandomExcerpt() (Excerpt, error) {
	panic("unimplemented")
}

func (i *Impl) InserrtUnsuccessflulTweetResponse(tweet twitter.Tweet) (twitter.Tweet, error) {
	panic("unimplemented")
}

func (i *Impl) InsertSuccessfullTweetResponse(res twitter.SucessfullTweetResponse) (twitter.Tweet, error) {
	panic("unimplemented")
}