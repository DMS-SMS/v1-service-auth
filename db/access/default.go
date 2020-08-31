package access

import (
	"github.com/jinzhu/gorm"
)

type _default struct {
	None
	tx *gorm.DB
}

func Default() *_default {
	return new(_default)
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