package training

import (
	"math"
	"time"
)

type BucketId uint8

const (
	DefaultBucket BucketId = iota
	BucketA
	BucketB
	ArchiveBucket BucketId = math.MaxUint8 - 1
	DoneBucket    BucketId = math.MaxUint8
)

type RecallAttempt struct {
	RecordedAt time.Time
	Success    bool
	CardId     CardId
}

type CardId int
type Card struct {
	ID          CardId
	Phrase      string
	Translation string
	Bucket      BucketId
	CreatedAt   time.Time
}
