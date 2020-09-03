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
func (g grade) KeyName() string { return "grade" }

type class int64
func (c class) Value() (value driver.Value, err error) {
	value = int64(c)
	if value == 0 { value = nil }
	return
}
func (c *class) Scan(v interface{}) (err error) { *c = class(v.(int64)); return }
func (c class) KeyName() string { return "class" }


type studentNumber int64
func (sn studentNumber) Value() (value driver.Value, err error) {
	value = int64(sn)
	if value == 0 { value = nil }
	return
}
func (sn *studentNumber) Scan(v interface{}) (err error) { *sn = studentNumber(v.(int64)); return }
func (sn studentNumber) KeyName() string { return "student_number" }

type uuid string
func (u uuid) Value() (driver.Value, error) { return string(u), nil }
func (u *uuid) Scan(v interface{}) (err error) { *u = uuid(v.(string)); return }
func (u uuid) KeyName() string { return "uuid" }

type studentID string
func (si studentID) Value() (driver.Value, error) { return string(si), nil }
func (si *studentID) Scan(v interface{}) (err error) { *si = studentID(v.(string)); return }
func (si studentID) KeyName() string { return "student_id" }

type studentPW string
func (sp studentPW) Value() (driver.Value, error) { return string(sp), nil }
func (sp *studentPW) Scan(v interface{}) (err error) { *sp = studentPW(v.(string)); return }
func (sp studentPW) KeyName() string { return "student_pw" }

type teacherID string
func (ti teacherID) Value() (driver.Value, error) { return string(ti), nil }
func (ti *teacherID) Scan(v interface{}) (err error) { *ti = teacherID(v.(string)); return }
func (ti teacherID) KeyName() string { return "teacher_id" }

type teacherPW string
func (tp teacherPW) Value() (driver.Value, error) { return string(tp), nil }
func (tp *teacherPW) Scan(v interface{}) (err error) { *tp = teacherPW(v.(string)); return }
func (tp teacherPW) KeyName() string { return "teacher_pw" }

type parentID string
func (pi parentID) Value() (driver.Value, error) { return string(pi), nil }
func (pi *parentID) Scan(v interface{}) (err error) { *pi = parentID(v.(string)); return }
func (pi parentID) KeyName() string { return "parent_id" }

type parentPW string
func (pp parentPW) Value() (driver.Value, error) { return string(pp), nil }
func (pp *parentPW) Scan(v interface{}) (err error) { *pp = parentPW(v.(string)); return }
func (pp parentPW) KeyName() string { return "parent_pw" }
