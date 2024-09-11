package data

import (
	"context"
	"fmt"

	"github.com/aaegamysta/listen-2-max-payne/internal/twitter"
	"github.com/jackc/pgx/v5"
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
	conn, err := pgx.Connect(ctx, i.connectionString)
	if err != nil {
		i.logger.Panicf("something wrong happened while connecting to the database: %v", err)
	}
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		i.logger.Panicf("something wrong happened while starting transaction for creating tables: %v", err)
	}
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS excerpts (
			series text, part text, chapter text, excerpt text,
			PRIMARY KEY (excerpt)
	);`)
	if err != nil {
		i.logger.Panicf("something wrong happened while creating excerpts table: %v", err)
	}
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS successfull_tweet_response (
			posted_on timestamp PRIMARY KEY, tweeted_excerpt TEXT, tweet_id TEXT, edit_history_tweet_ids JSONB,
			FOREIGN KEY (tweeted_excerpt) REFERENCES excerpts(excerpt)
	);`)
	if err != nil {
		i.logger.Panicf("something wrong happened while creating successfull_tweet_response table: %v", err)
	}
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS error_tweet_response (
			post_failed_on timestamp PRIMARY KEY, title TEXT, type TEXT, detail TEXT, status INT, failed_excerpt TEXT,
			FOREIGN KEY (failed_excerpt) REFERENCES excerpts(excerpt)
	);`)
	if err != nil {
		i.logger.Panicf("something wrong happened while creating error_tweet_response table: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		i.logger.Panicf("something wrong happened while commiting table creation statements: %v", err)
	}
	return nil
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