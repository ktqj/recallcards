package cards

import (
	"math"
	"time"
	// _ "golang.org/x/tools/cmd/stringer"
)

type BucketId uint8

const (
	DefaultBucket BucketId = iota
	BucketA
	BucketB
	ArchiveBucket BucketId = math.MaxUint8 - 1
	DoneBucket    BucketId = math.MaxUint8
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=BucketId

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
