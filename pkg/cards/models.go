package cards

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

// https://github.com/golang/go/issues/25922#issuecomment-1065971260
//go:generate go run golang.org/x/tools/cmd/stringer -type=BucketId

type RecallAttempt struct {
	RecordedAt time.Time
	Success    bool
	CardId     CardId
}
type Recalls []RecallAttempt

type RecallSummary struct {
	Ok   int
	Fail int
}

type CardId int
type Card struct {
	ID          CardId
	Phrase      string
	Translation string
	Bucket      BucketId
	CreatedAt   time.Time
}
type Cards []Card
