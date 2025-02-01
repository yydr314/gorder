package convertor

import "sync"

var (
	itemConvertor *ItemConvertor
	itemOnce      sync.Once
)

var (
	itemWithQuantityConvertor *ItemWithQuantityConvertor
	itemWithQuantityOnce      sync.Once
)

func NewItemConvertor() *ItemConvertor {
	itemOnce.Do(func() {
		itemConvertor = new(ItemConvertor)
	})
	return itemConvertor
}

func NewItemWihhQuantityConvertor() *ItemWithQuantityConvertor {
	itemWithQuantityOnce.Do(func() {
		itemWithQuantityConvertor = new(ItemWithQuantityConvertor)
	})
	return itemWithQuantityConvertor
}
