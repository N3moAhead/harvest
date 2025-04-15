package item

type ItemType int

const (
	Carrot ItemType = iota
	Potato
)

func NewCarrot() *Item {
	newItem := newItemBase()
	newItem.Type = Carrot
	return newItem
}

func NewPotato() *Item {
	newItem := newItemBase()
	newItem.Type = Potato
	return newItem
}
