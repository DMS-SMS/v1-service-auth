package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
)

type AccessorManage struct {
	accessorType reflect.Type
	dbForTx *gorm.DB
}

func NewAccessorManage(accessorType reflect.Type, dbForTx *gorm.DB) (manager AccessorManage, err error) {
	if accessorType == nil || dbForTx == nil {
		err = errors.New(fmt.Sprintf("nil parameter is not allowed"))
		return
	}

	if _, ok := reflect.New(accessorType).Elem().Interface().(Accessor); !ok {
		err = errors.New(fmt.Sprintf("type %s is not an implement of db.Accessor", accessorType.String()))
		return
	}

	manager = AccessorManage{
		accessorType: accessorType,
		dbForTx:      dbForTx,
	}
	return
}

func (atm AccessorManage) BeginTx() (accessor Accessor, err error) {
	if atm.accessorType == nil || atm.dbForTx == nil {
		err = errors.New("please create db.AccessorManage instance object through the constructor")
		return
	}

	accessor = reflect.New(atm.accessorType.Elem()).Interface().(Accessor)
	accessor.Begin(atm.dbForTx)
	return
}
