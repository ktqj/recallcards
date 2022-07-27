package storage

import (
	"fmt"
	"strconv"

	"example.com/recallcards/pkg/cards"
)

var idGenerator = make(chan cards.CardId)

func idCounter() {
	for i := 0; ; i++ {
		idGenerator <- cards.CardId(strconv.Itoa(i))
	}
}

func init() {
	go idCounter()
}

type mem struct {
	st map[cards.CardId]cards.Card
}

func (cr *mem) Insert(c cards.Card) error {
	c.ID = <-idGenerator
	cr.st[c.ID] = c
	fmt.Printf("%v", cr.st)
	return nil
}

func NewMemoryRepository() *mem {
	return &mem{st: make(map[cards.CardId]cards.Card)}
}


