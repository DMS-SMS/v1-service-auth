package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type grade struct {
	value int64
}

type class struct {
	value int64
}

// ---

func Grade(v int64) grade {
	return grade{v}
}

func Class(v int64) class {
	return class{v}
}

// ---

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

func (c class) Value() (driver.Value, error) {
	if !ClassValuesList.Contains(c.value) {
		return nil, errors.New(fmt.Sprintf("%d is outside the range of the class property", c.value))
	}
	return c.value, nil
}

func (c *class) Scan(v interface{}) error {
	c.value = v.(int64)
	return nil
}