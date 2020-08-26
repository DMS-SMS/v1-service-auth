package model

import (
	"database/sql/driver"
)

type grade struct {
	value int64
}

func (g grade) Value() (driver.Value, error) {
	return g.value, nil
}

func (g *grade) Scan(v interface{}) error {
	g.value = v.(int64)
	return nil
}