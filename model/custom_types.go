package model

import (
	"database/sql/driver"
)

type grade int64
func (g grade) Value() (value driver.Value, err error) {
	value = int64(g)
	if value == 0 { value = nil }
	return
}
func (g *grade) Scan(v interface{}) (_ error) { *g = grade(v.(int64)); return }

type class int64
func (c class) Value() (value driver.Value, err error) {
	value = int64(c)
	if value == 0 { value = nil }
	return
}
func (c *class) Scan(v interface{}) (err error) { *c = class(v.(int64)); return }

type studentNumber int64
func (sn studentNumber) Value() (value driver.Value, err error) {
	value = int64(sn)
	if value == 0 { value = nil }
	return
}
func (sn *studentNumber) Scan(v interface{}) (err error) { *sn = studentNumber(v.(int64)); return }