package parser

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func Test_parse(t *testing.T) {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		t.Fatalf("failed to create logger isntance: %v", err)
	}
	f, err := os.Open("./testdata/excerpts.json")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	p := New(logger.Sugar(), nil)
	pImpl := p.(*impl)
	Result := pImpl.parse(context.Background(), f, 10)
	start := time.Now()
	log.Printf("started parsing at %s", start)
	counter := 0
	for result := range Result {
		counter++
		if result.Error != nil {
			t.Fatalf("something wrong happened while parsing ")
		}
	}
	finish := time.Now()
	if parsingDuration := finish.Sub(start) / time.Millisecond; parsingDuration == 0 {
		log.Printf("finished parsing at %s parsing total %d taking less than 1 millisecond", finish, counter)
	} else {
		log.Printf("finished parsing at %s parsing total %d taking %d microsecond", finish, counter, parsingDuration)
	}
}
