package event

type PubSubPayload struct {
	Event     string `json:"event"`
	ArticleID int    `json:"article_id"`
}

const (
	// ArticleSubjectAdd :nodoc:
	ArticleSubjectAdd = "ARTICLE_ADD"

	// ArticleSubjectRemove :nodoc:
	ArticleSubjectRemove = "ARTICLE_REMOVE"
)
