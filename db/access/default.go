package access

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

type _default struct {
	tx *gorm.DB
}

func DefaultReflectType() reflect.Type {
	return reflect.TypeOf(&_default{})
}

func (d *_default) Begin(db *gorm.DB) {
	d.tx = db.Begin()
}

func (d *_default) Commit() *gorm.DB {
	return d.tx.Commit()
}

func (d *_default) Rollback() *gorm.DB {
	return d.tx.Rollback()
}