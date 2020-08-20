package db

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

type AccessorManage struct {
	accessor Accessor
	db *gorm.DB
}

func NewAccessorManage(accessor Accessor, db *gorm.DB) AccessorManage {
	return AccessorManage{
		accessor: accessor,
		db:       db,
	}
}

func (atm AccessorManage) BeginTx() (accessor Accessor) {
	t := reflect.TypeOf(atm.accessor)
	accessor = reflect.New(t).Interface().(Accessor)
	accessor.Begin(atm.db)
	return
}
