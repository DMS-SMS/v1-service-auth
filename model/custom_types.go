package model

import (
	"auth/tool/arraylist"
	"database/sql/driver"
	"errors"
	"fmt"
)

type grade struct {
	value int64
}

var gradeAvailableValues = []int64{1, 2, 3}
var GradeValuesList = arraylist.NewWithInt64(gradeAvailableValues...)
func (g grade) Value() (driver.Value, error) {
	if !GradeValuesList.Contains(g.value) {
		return nil, errors.New(fmt.Sprintf("%d is outside the range of the grade property", g.value))
	}
	return g.value, nil
}

func (g *grade) Scan(v interface{}) error {
	g.value = v.(int64)
	return nil
}