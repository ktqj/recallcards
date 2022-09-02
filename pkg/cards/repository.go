package cards

type CardRepository interface {
	InsertCard(c Card) error
	InsertRecallAttempt(r RecallAttempt) error
	ListUsedBuckets() ([]BucketId, error)
	RandomCardByBucket(b BucketId) (Card, error)
	CountRecallAttempts(cid CardId) int
}
