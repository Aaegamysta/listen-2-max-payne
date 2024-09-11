package db

import (
	"context"
	"fmt"

	"github.com/aaegamysta/listen-2-max-payne/internal/twitter"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Interface interface {
	CreateTablesIfNotExists(ctx context.Context) error
	BatchInsertExcerpts(ctxc context.Context, excerpts []Excerpt) ([]Excerpt, error)
	GetRandomExcerpt(ctx context.Context) (Excerpt, error)
	InsertSuccessfulTweetResponse(res twitter.SucessfullTweetResponse) (twitter.Tweet, error)
	InsertUnsuccessfulTweetResponse(tweet twitter.Tweet) (twitter.Tweet, error)
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

func (repository *Impl) CreateTablesIfNotExists(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, repository.connectionString)
	if err != nil {
		repository.logger.Panicf("something wrong happened while connecting to the database: %v", err)
	}
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repository.logger.Panicf("something wrong happened while starting transaction for creating tables: %v", err)
	}
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS excerpts (
			series int, part text, chapter text, excerpt text,
			PRIMARY KEY (excerpt)
	);`)
	if err != nil {
		repository.logger.Panicf("something wrong happened while creating excerpts table: %v", err)
	}
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS successfull_tweet_response (
			posted_on timestamp PRIMARY KEY, tweeted_excerpt TEXT, tweet_id TEXT, edit_history_tweet_ids JSONB,
			FOREIGN KEY (tweeted_excerpt) REFERENCES excerpts(excerpt)
	);`)
	if err != nil {
		repository.logger.Panicf("something wrong happened while creating successfull_tweet_response table: %v", err)
	}
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS error_tweet_response (
			post_failed_on timestamp PRIMARY KEY, title TEXT, type TEXT, detail TEXT, status INT, failed_excerpt TEXT,
			FOREIGN KEY (failed_excerpt) REFERENCES excerpts(excerpt)
	);`)
	if err != nil {
		repository.logger.Panicf("something wrong happened while creating error_tweet_response table: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		repository.logger.Panicf("something wrong happened while commiting table creation statements: %v", err)
	}
	return nil
}

func (repository *Impl) BatchInsertExcerpts(ctx context.Context, excerpts []Excerpt) ([]Excerpt, error) {
	conn, err := pgx.Connect(context.Background(), repository.connectionString)
	if err != nil {
		return nil, fmt.Errorf("something wrong happened while acquiring connection to the database: %w", err)
	}
	_, err = conn.CopyFrom(context.Background(),
		pgx.Identifier{"excerpts"}, []string{"series", "part", "chapter", "excerpt"},
		pgx.CopyFromSlice(len(excerpts), func(i int) ([]any, error) {
			if len(excerpts[i].Excerpt) > 280 {
				repository.logger.Warnf("found an excerpt that will be not be tweetable because it is more than 28 characters %s", excerpts[i].Excerpt)
			}
			return []any{excerpts[i].Series, excerpts[i].Part, excerpts[i].Chapter, excerpts[i].Excerpt}, nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("something wrong happened while batch inserting excerpts: %w", err)
	}
	return excerpts, nil
}

func (repository *Impl) GetRandomExcerpt(context.Context) (Excerpt, error) {
	conn, err := pgx.Connect(context.Background(), repository.connectionString)
	if err != nil {
		return Excerpt{}, fmt.Errorf("something wrong happened while acquiring connection to the database: %w", err)
	}
	var e Excerpt
	for len(e.Excerpt) < twitter.MaxTweetLength {
		row := conn.QueryRow(context.Background(), "SELECT * FROM excerpts ORDER BY random()")
		err = row.Scan(&e.Series, &e.Part, &e.Series, &e.Excerpt)
		if err != nil {
			return Excerpt{}, fmt.Errorf("something wrong happened while fetching a random excerpt: %w", err)
		}
	}
	return e, nil
}

func (repository *Impl) InsertUnsuccessfulTweetResponse(tweet twitter.Tweet) (twitter.Tweet, error) {
	panic("unimplemented")
}

func (repository *Impl) InsertSuccessfulTweetResponse(res twitter.SucessfullTweetResponse) (twitter.Tweet, error) {
	panic("unimplemented")
}
