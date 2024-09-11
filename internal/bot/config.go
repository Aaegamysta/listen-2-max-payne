package bot

import (
	"github.com/aaegamysta/listen-2-max-payne/internal/db"
	"github.com/aaegamysta/listen-2-max-payne/internal/facade/parser"
	"github.com/aaegamysta/listen-2-max-payne/internal/facade/publisher"
	"github.com/aaegamysta/listen-2-max-payne/internal/twitter"
)

type Config struct {
	Database db.Config `yaml:"psql"`
	Twitter twitter.Config `yaml:"twitter"`
	Parser parser.Config `yaml:"parser"`
	Publisher publisher.Config `yaml:"publisher"`
}
