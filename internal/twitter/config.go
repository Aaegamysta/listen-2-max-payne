package twitter

const MaxTweetLength = 280

type Config struct {
	ConsumerKey    string `yaml:"consumerKey"`
	ConsumerSecret string `yaml:"consumerSecret"`
	AccessToken    string `yaml:"accessToken"`
	AccessSecret   string `yaml:"accessSecret"`
	Endpoint       string `yaml:"endpoint"`
}
