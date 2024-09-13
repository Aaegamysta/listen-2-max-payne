package parser

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aaegamysta/listen-2-max-payne/internal/db"
	"go.uber.org/zap"
)

type Interface interface {
	ParseAndSaveExcerpts(ctx context.Context, reader io.Reader) error
}

type Result struct {
	Excerpt db.Excerpt
	Error   error
}

type impl struct {
	logger               *zap.SugaredLogger
	repository           db.Interface
	batchInsertChunkSize int
}

func (i *impl) ParseAndSaveExcerpts(ctx context.Context, reader io.Reader) error {
	excerptResultsStream := i.parse(ctx, reader, i.batchInsertChunkSize)
	err := i.save(ctx, i.batchInsertChunkSize, excerptResultsStream)
	if err != nil {
		return err
	}
	return io.EOF
}

func (i *impl) parse(_ context.Context, reader io.Reader, chunks int) <-chan Result {
	excerptsResultsStream := make(chan Result, chunks)
	decoder := json.NewDecoder(reader)
	go func() {
		defer close(excerptsResultsStream)
		// Read the opening brace of the JSON object
		_, err := decoder.Token()
		if err != nil {
			i.logger.Panicf("something wrong happened while reading opening brace:", err)
			return
		}

		// Read the "excerpts" key
		_, err = decoder.Token()
		if err != nil {
			i.logger.Panicf("something wrong happened while parsing excerpts key:", err)
			return
		}

		// Read the opening bracket of the excerpts array
		_, err = decoder.Token()
		if err != nil {
			i.logger.Panicf("something wrong happened while parsing opening bracket:", err)
			return
		}
		for decoder.More() {
			var e db.Excerpt
			err := decoder.Decode(&e)
			excerptsResultsStream <- Result{
				Excerpt: e,
				Error:   err,
			}
		}
		i.logger.Infof("finishing parsing all excerpts")
	}()
	return excerptsResultsStream
}

func (i *impl) save(ctx context.Context, chunks int, excerptsResults <-chan Result) error {
	excerpts := make([]db.Excerpt, 0)
	counter := 0
	for result := range excerptsResults {
		if counter == chunks {
			_, err := i.repository.BatchInsertExcerpts(ctx, excerpts)
			if err != nil {
				return err
			}
			excerpts = make([]db.Excerpt, 0)
			counter = 0
		}
		if result.Error != nil {
			return result.Error
		}
		excerpts = append(excerpts, result.Excerpt)
		counter++
	}
	_, err := i.repository.BatchInsertExcerpts(ctx, excerpts)
	if err != nil {
		return err
	}
	i.logger.Infof("finishing saving all excerpts")
	return nil
}

func New(_ context.Context, cfg Config, logger *zap.SugaredLogger, repo db.Interface) Interface {
	return &impl{
		logger:               logger,
		repository:           repo,
		batchInsertChunkSize: cfg.BatchInsertChunkSize,
	}
}
