package file

import (
	"testing"

	"example.com/recallcards/pkg/cards"
)

var objectsGlobal, _ = readJsonFile[cards.Cards]("./testdata/cards.json")

func BenchmarkCardsLoopingRangeIndexOnlyGlobal(b *testing.B) {
	// BenchmarkCardsLoopingRangeIndexOnly
	// BenchmarkCardsLoopingRangeIndexOnly-12    	1000000000	         0.5007 ns/op	       0 B/op	       0 allocs/op
	for i := 0; i < b.N; i++ {
		count := 0
		for j := range objectsGlobal {
			if objectsGlobal[j].Bucket == cards.DefaultBucket {
				count++
			}
		}
	}
}

func BenchmarkCardsLoopingForIndexGlobal(b *testing.B) {
	// BenchmarkCardsLoopingForIndex
	// BenchmarkCardsLoopingForIndex-12          	1000000000	         0.5083 ns/op	       0 B/op	       0 allocs/op
	for i := 0; i < b.N; i++ {
		count := 0
		for j := 0; j < len(objectsGlobal); j++ {
			if objectsGlobal[j].Bucket == cards.DefaultBucket {
				count++
			}
		}
	}
}

func BenchmarkCardsLoopingFullRangeGlobal(b *testing.B) {
	// BenchmarkCardsLoopingFullRange
	// BenchmarkCardsLoopingFullRange-12         	1000000000	         0.2462 ns/op	       0 B/op	       0 allocs/op
	for i := 0; i < b.N; i++ {
		count := 0
		for _, elem := range objectsGlobal {
			if elem.Bucket == cards.DefaultBucket {
				count++
			}
		}
	}
}

func BenchmarkCardsLoopingRangeIndexOnly(b *testing.B) {
	// BenchmarkCardsLoopingRangeIndexOnly
	// BenchmarkCardsLoopingRangeIndexOnly-12    	1000000000	         0.7389 ns/op	       0 B/op	       0 allocs/op
	var objects, _ = readJsonFile[cards.Cards]("./testdata/cards.json")
	for i := 0; i < b.N; i++ {
		count := 0
		for j := range objects {
			if objects[j].Bucket == cards.DefaultBucket {
				count++
			}
		}
	}
}

func BenchmarkCardsLoopingForIndex(b *testing.B) {
	// BenchmarkCardsLoopingForIndex
	// BenchmarkCardsLoopingForIndex-12          	1000000000	         0.7361 ns/op	       0 B/op	       0 allocs/op
	var objects, _ = readJsonFile[cards.Cards]("./testdata/cards.json")
	for i := 0; i < b.N; i++ {
		count := 0
		for j := 0; j < len(objects); j++ {
			if objects[j].Bucket == cards.DefaultBucket {
				count++
			}
		}
	}
}

func BenchmarkCardsLoopingFullRange(b *testing.B) {
	// BenchmarkCardsLoopingFullRange
	// BenchmarkCardsLoopingFullRange-12         	1000000000	         0.2436 ns/op	       0 B/op	       0 allocs/op
	var objects, _ = readJsonFile[cards.Cards]("./testdata/cards.json")
	for i := 0; i < b.N; i++ {
		count := 0
		for _, elem := range objects {
			if elem.Bucket == cards.DefaultBucket {
				count++
			}
		}
	}
}