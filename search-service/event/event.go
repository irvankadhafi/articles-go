package event

type PubSubPayload struct {
	Event     string `json:"event"`
	ArticleID int    `json:"article_id"`
}
