package fzf

import (
	"github.com/reinhrst/fzf-lib/util"
)

// Item represents each input line. 56 bytes.
type Item struct {
	text util.Chars // 32 = 24 + 1 + 1 + 2 + 4
}

// Index returns ordinal index of the Item
func (item *Item) Index() int32 {
	return item.text.Index
}

var minItem = Item{text: util.Chars{Index: -1}}

func (item *Item) TrimLength() uint16 {
	return item.text.TrimLength()
}
