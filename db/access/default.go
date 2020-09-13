package access

import (
	"github.com/jinzhu/gorm"
)

type _default struct {
	tx *gorm.DB
}

func Default(tx *gorm.DB) *_default {
	return &_default{tx: tx}
}

func (d *_default) BeginTx() {
	d.tx = d.tx.Begin()
}

func (d *_default) Commit() *gorm.DB {
	return d.tx.Commit()
}

func (d *_default) Rollback() *gorm.DB {
	return d.tx.Rollback()
}