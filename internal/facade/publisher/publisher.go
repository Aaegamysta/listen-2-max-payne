package publisher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aaegamysta/listen-2-max-payne/internal/db"
	"github.com/aaegamysta/listen-2-max-payne/internal/queue"
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
	doubleEndedQueue queue.Dequeue
}

func New(logger *zap.SugaredLogger, cfg Config, repository db.Interface, twitterClient twitter.Interface) Interface {
	return &Impl{
		logger:            logger,
		repository:        repository,
		twitterClient:     twitterClient,
		tweetPeriodPerDay: time.Duration(cfg.TweetPeriodPerDay),
		doubleEndedQueue:  queue.New(20),
	}
}

func (i *Impl) StartPublishingExcerpts(ctx context.Context) {
	go func() {
		// for purpose of testing setting it as a minute
		t := time.NewTicker(1 * time.Minute)
		// t := time.NewTicker(1 * i.tweetPeriodPerDay)
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
	for i.doubleEndedQueue.Peek().Excerpt == excerpt.Excerpt {
		excerpt, err = i.repository.GetRandomExcerpt(ctx)
		if err != nil {
			return fmt.Errorf(`failed to continuously retrieve random excerpt because front element is equal to the fetched excerpt: %w`,
				err)
		}
	}
	_, _ = i.doubleEndedQueue.Dequeue()
	_ = i.doubleEndedQueue.Enqueue(excerpt)
	successfulTweetRes, err := i.twitterClient.Post(ctx, twitter.Tweet{
		Text: excerpt.Excerpt,
	})

	var unsuccessfullTweetResponse twitter.TweetError
	if errors.As(err, &unsuccessfullTweetResponse) {
		// here an error can be generated if the unsuccsesful tweet response is not sent
		return i.repository.InsertUnsuccessfulTweetResponse(ctx, excerpt, twitter.TweetError{})
	} else if err != nil {
		return err
	}

	err = i.repository.InsertSuccessfulTweetResponse(ctx, successfulTweetRes)
	if err != nil {
		return fmt.Errorf("failed to insert successful tweet response but it was at least tweeted: %w", err)
	}
	i.logger.Infof("tweeted excerpt: %v on %s", successfulTweetRes.Text, time.Now())
	return nil
}
