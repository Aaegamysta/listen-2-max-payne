package parser

import (
	"context"
	"encoding/json"
	"github.com/aaegamysta/listen-2-max-payne/internal/db"
	"go.uber.org/zap"
	"io"
)

type Interface interface {
	Parse(ctx context.Context, reader io.Reader, chunks int) <-chan Result
}

type Result struct {
	Excerpt db.Excerpt
	Error   error
}

type impl struct {
	logger *zap.SugaredLogger
}

// Parse implements Interface.
func (i *impl) Parse(ctx context.Context, reader io.Reader, chunks int) <-chan Result {
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
	}()
	return excerptsResultsStream
}

func New(logger *zap.SugaredLogger) Interface {
	return &impl{
		logger: logger,
	}
}
