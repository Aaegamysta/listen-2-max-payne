package twitter

import (
	"context"
	"encoding/json"
	"github.com/dghubble/oauth1"
	"net/http"

	"go.uber.org/zap"
)

type Interface interface {
	Post(ctx context.Context, tweet Tweet) (SucessfullTweetResponse, error)
}

type Impl struct {
	logger     *zap.SugaredLogger
	httpClient *http.Client
	endpoint   string
}

func New(ctx context.Context, cfg Config, logger *zap.SugaredLogger) Interface {
	oauth1Config := oauth1.NewConfig(cfg.ConsumerKey, cfg.ConsumerSecret)
	token := oauth1.NewToken(cfg.AccessToken, cfg.AccessSecret)
	// What are the repercussion of having context.Background backed into the oauth1ConfiguredClient and it doesn't change or rather not creating a oauth1ConfiguredClient with every call
	oauth1ConfiguredClient := oauth1Config.Client(context.Background(), token)
	return &Impl{
		logger:     logger,
		httpClient: oauth1ConfiguredClient,
		endpoint:   cfg.Endpoint,
	}
}

func (i *Impl) Post(ctx context.Context, e Tweet) (SucessfullTweetResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.endpoint, nil)
	if err != nil {
		i.logger.Panicf("failed to create request", err)
	}
	res, err := i.httpClient.Do(req)
	if err != nil {
		panic("handle this error")
	}
	if res.StatusCode != http.StatusCreated {
		panic("handle this error")
	}
	if res.Body == nil {
		panic("handle this error")
	}
	defer res.Body.Close()
	var successfulTweetRes SucessfullTweetResponse
	if err := json.NewDecoder(res.Body).Decode(&successfulTweetRes); err != nil {
		panic("handle this error")
	}
	return successfulTweetRes, nil
}
