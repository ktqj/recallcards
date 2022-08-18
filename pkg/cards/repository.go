package cards

type CardRepository interface {
	Insert(c Card) error
	ListUsedBuckets() ([]BucketId, error)
	// CountByBucket(b BucketId) (int, error)
	RandomByBucket(b BucketId) (Card, error)
}