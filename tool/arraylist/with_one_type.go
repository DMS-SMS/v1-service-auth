package arraylist

import (
	"github.com/emirpasic/gods/lists/arraylist"
)

func NewWithInt64(values ...int64) *arraylist.List {
	list := make([]interface{}, len(values))
	for i, v := range values {
		list[i] = v
	}
	return arraylist.New(list...)
}