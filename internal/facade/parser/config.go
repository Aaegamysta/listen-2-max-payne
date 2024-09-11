package parser

type Config struct {
	Chunks int `yaml:"batchInsertChunkSize"`
}
