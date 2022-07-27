package cards

type CardRepository interface {
	Insert(c Card) error
}