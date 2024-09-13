package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dghubble/oauth1"
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
	// What are the repercussion of having context.Background backed into the oauth1ConfiguredClient
	// and it doesn't change or rather not creating a
	// oauth1ConfiguredClient with every call
	oauth1ConfiguredClient := oauth1Config.Client(ctx, token)
	return &Impl{
		logger:     logger,
		httpClient: oauth1ConfiguredClient,
		endpoint:   cfg.Endpoint,
	}
}

func (i *Impl) Post(ctx context.Context, e Tweet) (SucessfullTweetResponse, error) {
	tweet := map[string]interface{}{
		"text": e.Text,
	}
	jsonData, err := json.Marshal(tweet)
	if err != nil {
		i.logger.Panicf("failed to marshal tweet: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		i.logger.Panicf("failed to create request %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := i.httpClient.Do(req)
	if err != nil {
		return SucessfullTweetResponse{}, fmt.Errorf("something happened while posting doing request: %w", err)
	}
	if res.Body == nil {
		return SucessfullTweetResponse{}, errors.New("something happened while posting doing request: the response body was empty")
	}
	if res.StatusCode != http.StatusCreated {
		var tweetError TweetError
		b, _ := io.ReadAll(res.Body)
		err = json.Unmarshal(b, &tweetError)
		if err != nil {
			return SucessfullTweetResponse{}, tweetError
		}
		return SucessfullTweetResponse{},
			fmt.Errorf("failed to unmarshall the tweet error response %w here is the stringified response: %s", err, b)
	}
	defer res.Body.Close()
	var successfulTweetRes SucessfullTweetResponse
	if err := json.NewDecoder(res.Body).Decode(&successfulTweetRes); err != nil {
		// the response will be stringified because we do not know the correct struct format to parse it
		b, _ := io.ReadAll(res.Body)
		return SucessfullTweetResponse{},
			fmt.Errorf("failed to unmarshall the successful tweet response, maybe it failed? : %w the body looks like so %s", err, b)
	}
	return successfulTweetRes, nil
}
