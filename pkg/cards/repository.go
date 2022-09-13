package cards

type Repository interface {
	InsertCard(c Card) error
	InsertRecallAttempt(r RecallAttempt) error
	CardById(cid CardId) (Card, error)
	ListCardIds() ([]CardId, error)
	ListCards() ([]Card, error)
	ListUsedBuckets() ([]BucketId, error)
	RandomCardByBucket(b BucketId) (Card, error)
	RecallSummary(cid CardId) RecallSummary
}
