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

type grade struct { value int64 }
func Grade(v int64) grade { return grade{v} }
func (g grade) Value() (value driver.Value, err error) {
	value = g.value
	if g.value == int64(0) { return }
	if !GradeValuesList.Contains(g.value) {
		err = errors.New(fmt.Sprintf("%d is outside the range of the grade property", g.value))
	}
	return
}
func (g *grade) Scan(v interface{}) (err error) {
	g.value = v.(int64)
	return
}


type class struct { value int64 }
func Class(v int64) class { return class{v} }
func (c class) Value() (value driver.Value, err error) {
	value = c.value
	if c.value == int64(0) { return }
	if !ClassValuesList.Contains(c.value) {
		err = errors.New(fmt.Sprintf("%d is outside the range of the class property", c.value))
	}
	return
}
func (c *class) Scan(v interface{}) (err error) {
	c.value = v.(int64)
	return
}

type studentNumber struct { value int64 }
func StudentNumber(v int64) studentNumber { return studentNumber{v} }
func (sn studentNumber) Value() (value driver.Value, err error) {
	value = sn.value
	if !StudentNumberValuesList.Contains(sn.value) {
		err = errors.New(fmt.Sprintf("%d is outside the range of the student_number property", sn.value))
	}
	return
}
func (sn studentNumber) Scan(v interface{}) (err error) {
	sn.value = v.(int64)
	return
}