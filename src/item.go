package fzf

import (
	"github.com/junegunn/fzf/src/util"
)

// Item represents each input line. 56 bytes.
type Item struct {
	text        util.Chars    // 32 = 24 + 1 + 1 + 2 + 4
	transformed *[]Token      // 8
	origText    *[]byte       // 8
}

// Index returns ordinal index of the Item
func (item *Item) Index() int32 {
	return item.text.Index
}

var minItem = Item{text: util.Chars{Index: -1}}

func (item *Item) TrimLength() uint16 {
	return item.text.TrimLength()
}
