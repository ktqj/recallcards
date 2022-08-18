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
  DoneBucket BucketId = math.MaxUint8
)

type RecallAttempt struct {
  ts time.Time
  success bool
}

type CardId string
type Card struct {
  ID CardId
  Phrase string
  Translation string
  Bucket BucketId
  Created_at time.Time
}
