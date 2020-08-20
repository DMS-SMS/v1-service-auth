package access

import "github.com/jinzhu/gorm"

type Default struct {
	tx *gorm.DB
}

func (d *Default) Begin(db *gorm.DB) {
	d.tx = db.Begin()
}

func (d *Default) Commit() *gorm.DB {
	return d.tx.Commit()
}

func (d *Default) Rollback() *gorm.DB {
	return d.tx.Rollback()
}