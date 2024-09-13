package twitter

import "fmt"

type SucessfullTweetResponse struct {
	Data                TweetData `json:"data"`
	Text                string    `json:"text"`
	EditHistoryTweetIDs []string  `json:"edit_history_tweet_ids"`
}

type TweetData struct {
	ID string `json:"id"`
}

type TweetError struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

func (e TweetError) Error() string {
	return fmt.Sprintf("title: %s, type: %s, detail: %s, status: %d", e.Title, e.Type, e.Detail, e.Status)
}
