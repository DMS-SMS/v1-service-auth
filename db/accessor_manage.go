package db

import (
	"errors"
	"fmt"
	"reflect"
)

type AccessorManage struct {
	accessorType  reflect.Type
	accessorValue reflect.Value
}

func NewAccessorManage(accessor Accessor) (manager AccessorManage, err error) {
	if accessor == nil {
		err = errors.New(fmt.Sprintf("nil parameter is not allowed"))
		return
	}

	accessorType := reflect.TypeOf(accessor)
	accessorValue := reflect.ValueOf(accessor)

	if accessorType.Kind() == reflect.Ptr {
		accessorType = accessorType.Elem()
		accessorValue = accessorValue.Elem()
	}

	manager = AccessorManage{
		accessorType:  accessorType,
		accessorValue: accessorValue,
	}
	return
}

func (atm AccessorManage) BeginTx() (accessor Accessor, err error) {
	if atm.accessorType == nil {
		err = errors.New("please create db.AccessorManage instance object through the constructor")
		return
	}

	newAccessor := reflect.New(atm.accessorType)
	newAccessor.Elem().Set(atm.accessorValue)

	accessor = newAccessor.Interface().(Accessor)
	accessor.BeginTx()
	return
}
