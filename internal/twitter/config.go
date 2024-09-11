package twitter

const MaxTweetLength = 280

type Config struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessSecret      string `json:"accessSecret"`
	Endpoint         string `json:"endpoint"`
}