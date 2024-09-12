package publisher

import (
	"context"
	"errors"
	"time"

	"github.com/aaegamysta/listen-2-max-payne/internal/db"
	"github.com/aaegamysta/listen-2-max-payne/internal/twitter"
	"go.uber.org/zap"
)

type Interface interface {
	StartPublishingExcerpts(ctx context.Context)
}

type Impl struct {
	logger            *zap.SugaredLogger
	tweetPeriodPerDay time.Duration
	repository        db.Interface
	twitterClient     twitter.Interface
	// TODO: Double ended queue to prevent previously tweeted
	lruTweets any
}

func New(logger *zap.SugaredLogger, cfg Config, repository db.Interface, twitterClient twitter.Interface) Interface {
	return &Impl{
		logger:            logger,
		repository:        repository,
		twitterClient:     twitterClient,
		tweetPeriodPerDay: time.Duration(cfg.TweetPeriodPerDay),
	}
}

func (i *Impl) StartPublishingExcerpts(ctx context.Context) {
	go func() {
		t := time.NewTicker(1 * i.tweetPeriodPerDay)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				err := i.tweet(ctx)
				if err != nil {
					i.logger.Errorf("failed to tweet excerpt: %v, skipping this one", err)
				}
			}
		}
	}()
}

func (i *Impl) tweet(ctx context.Context) error {
	excerpt, err := i.repository.GetRandomExcerpt(ctx)
	if err != nil {
		i.logger.Errorf("failed to retrieve random excerpt for tweeting: %v", err)
	}
	// TODO: Check if it is present in the cache or queue if present, keep fetching randomly, if not present, tweet it and add it to the cache. If present keep fetching randomly
	successfulTweetRes, err := i.twitterClient.Post(ctx, twitter.Tweet{
		Text: excerpt.Excerpt,
	})
	var unsuccessfullTweetResponse twitter.UnsucessfullTweetResponse
	if errors.As(err, &unsuccessfullTweetResponse) {
		err = i.repository.InsertUnsuccessfulTweetResponse(twitter.UnsucessfullTweetResponse{})
	}
	if err != nil {
		return err
	}
	err = i.repository.InsertSuccessfulTweetResponse(successfulTweetRes)
	if err != nil {
		i.logger.Errorf("failed to insert successful tweet response but it was at least tweeted: %v", err)
	}
	i.logger.Infof("tweeted excerpt: %v on %s", successfulTweetRes.Text, time.Now())
	return nil
}
