package model

import (
	"database/sql/driver"
)

type Grade int64
func (g Grade) Value() (value driver.Value, err error) {
	value = int64(g)
	if value == 0 { value = nil }
	return
}
func (g *Grade) Scan(src interface{}) (_ error) { *g = Grade(src.(int64)); return }
func (g Grade) KeyName() string { return "grade" }

type Class int64
func (c Class) Value() (value driver.Value, err error) {
	value = int64(c)
	if value == 0 { value = nil }
	return
}
func (c *Class) Scan(src interface{}) (err error) { *c = Class(src.(int64)); return }
func (c Class) KeyName() string { return "class" }


type StudentNumber int64
func (sn StudentNumber) Value() (value driver.Value, err error) {
	value = int64(sn)
	if value == 0 { value = nil }
	return
}
func (sn *StudentNumber) Scan(src interface{}) (err error) { *sn = StudentNumber(src.(int64)); return }
func (sn StudentNumber) KeyName() string { return "student_number" }

type UUID string
func (u UUID) Value() (driver.Value, error) { return string(u), nil }
func (u *UUID) Scan(src interface{}) (err error) { *u = UUID(src.(string)); return }
func (u UUID) KeyName() string { return "UUID" }

type StudentID string
func (si StudentID) Value() (driver.Value, error) { return string(si), nil }
func (si *StudentID) Scan(src interface{}) (err error) { *si = StudentID(src.(string)); return }
func (si StudentID) KeyName() string { return "student_id" }

type StudentPW string
func (sp StudentPW) Value() (driver.Value, error) { return string(sp), nil }
func (sp *StudentPW) Scan(src interface{}) (err error) { *sp = StudentPW(src.(string)); return }
func (sp StudentPW) KeyName() string { return "student_pw" }

type TeacherID string
func (ti TeacherID) Value() (driver.Value, error) { return string(ti), nil }
func (ti *TeacherID) Scan(src interface{}) (err error) { *ti = TeacherID(src.(string)); return }
func (ti TeacherID) KeyName() string { return "teacher_id" }

type TeacherPW string
func (tp TeacherPW) Value() (driver.Value, error) { return string(tp), nil }
func (tp *TeacherPW) Scan(src interface{}) (err error) { *tp = TeacherPW(src.(string)); return }
func (tp TeacherPW) KeyName() string { return "teacher_pw" }

type ParentID string
func (pi ParentID) Value() (driver.Value, error) { return string(pi), nil }
func (pi *ParentID) Scan(src interface{}) (err error) { *pi = ParentID(src.(string)); return }
func (pi ParentID) KeyName() string { return "parent_id" }

type ParentPW string
func (pp ParentPW) Value() (driver.Value, error) { return string(pp), nil }
func (pp *ParentPW) Scan(src interface{}) (err error) { *pp = ParentPW(src.(string)); return }
func (pp ParentPW) KeyName() string { return "parent_pw" }

type StudentUUID string
func (su StudentUUID) Value() (driver.Value, error) { return string(su), nil }
func (su *StudentUUID) Scan(src interface{}) (err error) { *su = StudentUUID(src.(string)); return }
func (su StudentUUID) KeyName() string { return "student_uuid" }

type TeacherUUID string
func (tu TeacherUUID) Value() (driver.Value, error) { return string(tu), nil }
func (tu *TeacherUUID) Scan(src interface{}) (err error) { *tu = TeacherUUID(src.(string)); return }
func (tu TeacherUUID) KeyName() string { return "teacher_uuid" }

type ParentUUID string
func (pu ParentUUID) Value() (value driver.Value, err error) {
	value = string(pu)
	if value == "" { value = nil }
	return
}
func (pu *ParentUUID) Scan(src interface{}) (err error) { *pu = ParentUUID(src.(string)); return }
func (pu ParentUUID) KeyName() string { return "parent_uuid" }

type Name string
func (n Name) Value() (driver.Value, error) { return string(n), nil }
func (n *Name) Scan(src interface{}) (err error) { *n = Name(src.(string)); return }
func (n Name) KeyName() string { return "name" }

type PhoneNumber string
func (pn PhoneNumber) Value() (driver.Value, error) { return string(pn), nil }
func (pn *PhoneNumber) Scan(src interface{}) (err error) { *pn = PhoneNumber(src.(string)); return }
func (pn PhoneNumber) KeyName() string { return "phone_number" }

type ProfileURI string
func (pu ProfileURI) Value() (driver.Value, error) { return string(pu), nil }
func (pu *ProfileURI) Scan(src interface{}) (err error) { *pu = ProfileURI(src.(string)); return }
func (pu ProfileURI) KeyName() string { return "profile_uri" }