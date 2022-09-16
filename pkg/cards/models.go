package cards

import (
	"math"
	"time"
)

// TODO: buckets idea may not be useful or practical
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

// EstimateConfidence estimates confidence in knowing card's information
// based on its recalls history. Returns a percentage value, between 0 and 100.
func (summary RecallSummary) EstimateConfidence() int {
	// card was never recalled correctly yet
	if summary.Ok == 0 {
		return 0
	}

	// weighted difference between correct and failed recalls
	failWeight := 2
	diff := summary.Ok - failWeight*summary.Fail

	flawlessConfidentDifference := 5
	if diff >= flawlessConfidentDifference {
		return 95
	}

	confidentDifference := 10
	if diff >= confidentDifference {
		return 100
	}

	inconfidentDifference := -10
	if diff <= inconfidentDifference {
		return 0
	}

	// inconfidentDifference < diff < confidentDifference
	base := 50.0
	if diff >= 0 {
		progress := (100.0 - base) / float64(confidentDifference) * float64(diff)
		return int(base + progress)
	}

	regress := base / float64(inconfidentDifference) * float64(diff)
	return int(base - regress)
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
