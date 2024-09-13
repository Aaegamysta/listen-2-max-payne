package bot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aaegamysta/listen-2-max-payne/internal/db"
	"github.com/aaegamysta/listen-2-max-payne/internal/facade/parser"
	"github.com/aaegamysta/listen-2-max-payne/internal/facade/publisher"
	"github.com/aaegamysta/listen-2-max-payne/internal/twitter"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Bot struct {
	cfg           Config
	logger        *zap.SugaredLogger
	Parser        parser.Interface
	Publisher     publisher.Interface
}

func New(ctx context.Context, env string) *Bot {
	loggerConfig := zap.NewDevelopmentConfig()
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Panicf("something wrong happened while building logger: %v", err)
	}
	sugaredLogger := logger.Sugar()
	configFilePath := fmt.Sprintf("./configs/app.%s.yaml", env)
	file, err := os.Open(configFilePath)
	if err != nil {
		sugaredLogger.Panicf("something wrong happened while opening config file: %v", err)
	}
	defer file.Close()
	fileInBytes, err := io.ReadAll(file)
	if err != nil {
		sugaredLogger.Panicf("something wrong happened while reading config file: %v", err)
	}
	var cfg Config
	err = yaml.Unmarshal(fileInBytes, &cfg)
	if err != nil {
		sugaredLogger.Panicf("something wrong happened while unmarshalling config file: %v", err)
	}
	repo := db.New(ctx, cfg.Database, sugaredLogger)
	parser := parser.New(ctx, cfg.Parser, sugaredLogger, repo)
	twitterClient := twitter.New(ctx, cfg.Twitter, sugaredLogger)
	publisher := publisher.New(sugaredLogger, cfg.Publisher, repo, twitterClient)
	bot := &Bot{
		cfg:           cfg,
		Parser:        parser,
		Publisher:     publisher,
	}
	return bot
}

func (b *Bot) Run(ctx context.Context) {
	f, err := os.Open("./data/excerpts.json")
	if err != nil {
		b.logger.Panicf("something wrong happened while opening excerpts file: %v", err)
	}
	err = b.Parser.ParseAndSaveExcerpts(ctx, f)
	if errors.Is(err, io.EOF) {
		b.logger.Info("excerpts parsed and saved successfully")
	}
	if err != nil {
		b.logger.Panicf("something wrong happened while parsing and saving excerpts: %v", err)
	}
	b.Publisher.StartPublishingExcerpts(ctx)
}
