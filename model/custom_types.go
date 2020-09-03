package model

import (
	"auth/tool/arraylist"
	"database/sql/driver"
	"errors"
	"fmt"
)

var (
	gradeAvailableValues = []int64{1, 2, 3}
	GradeValuesList = arraylist.NewWithInt64(gradeAvailableValues...)

	classAvailableValues = []int64{1, 2, 3, 4}
	ClassValuesList = arraylist.NewWithInt64(classAvailableValues...)

	studentNumberAvailableValues = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
	StudentNumberValuesList = arraylist.NewWithInt64(studentNumberAvailableValues...)
)

type grade struct {
	value int64
}

type class struct {
	value int64
}

type studentNumber struct {
	value int64
}

// ---

func Grade(v int64) grade {
	return grade{v}
}

func Class(v int64) class {
	return class{v}
}

func StudentNumber(v int64) studentNumber {
	return studentNumber{v}
}

// ---

func (g grade) Value() (driver.Value, error) {
	if g.value == 0 {
		return nil, nil
	}

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
	if c.value == 0 {
		return nil, nil
	}

	if !ClassValuesList.Contains(c.value) {
		return nil, errors.New(fmt.Sprintf("%d is outside the range of the class property", c.value))
	}
	return c.value, nil
}

func (c *class) Scan(v interface{}) error {
	c.value = v.(int64)
	return nil
}

func (sn studentNumber) Value() (driver.Value, error) {
	if !StudentNumberValuesList.Contains(sn.value) {
		return nil, errors.New(fmt.Sprintf("%d is outside the range of the student_number property", sn.value))
	}
	return sn.value, nil
}

func (sn studentNumber) Scan(v interface{}) error {
	sn.value = v.(int64)
	return nil
}