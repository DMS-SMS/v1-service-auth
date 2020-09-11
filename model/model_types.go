package model

import (
	"auth/tool/random"
	"database/sql/driver"
)

var (
	nullReplaceValueForGrade int64
	nullReplaceValueForClass int64
	nullReplaceValueForParentUUID string
)

func init() {
	nullReplaceValueForGrade = random.Int64WithLength(10)
	nullReplaceValueForClass = random.Int64WithLength(10)
	nullReplaceValueForParentUUID = random.StringConsistOfIntWithLength(10)
}

// Grade 필드에서 사용할 사용자 정의 타입
type grade int64
func Grade(i int64) grade { return grade(i) }
func (g grade) Value() (value driver.Value, err error) {
	value = int64(g)
	if value == int64(0) { value = nil }
	return
}
func (g *grade) Scan(src interface{}) (_ error) { *g = grade(src.(int64)); return }
func (g grade) KeyName() string { return "grade" }
func (g grade) NullReplaceValue() int64 { return nullReplaceValueForGrade  }

// Class 필드에서 사용할 사용자 정의 타입
type class int64
func Class(i int64) class { return class(i) }
func (c class) Value() (value driver.Value, err error) {
	value = int64(c)
	if value == int64(0) { value = nil }
	return
}
func (c *class) Scan(src interface{}) (err error) { *c = class(src.(int64)); return }
func (c class) KeyName() string { return "class" }
func (c class) NullReplaceValue() int64 { return nullReplaceValueForClass  }

// StudentNumber 필드에서 사용할 사용자 정의 타입
type studentNumber int64
func StudentNumber(i int64) studentNumber { return studentNumber(i) }
func (sn studentNumber) Value() (driver.Value, error) { return int64(sn), nil }
func (sn *studentNumber) Scan(src interface{}) (err error) { *sn = studentNumber(src.(int64)); return }
func (sn studentNumber) KeyName() string { return "student_number" }

// UUID 필드에서 사용할 사용자 정의 타입
type uuid string
func UUID(s string) uuid { return uuid(s) }
func (u uuid) Value() (driver.Value, error) { return string(u), nil }
func (u *uuid) Scan(src interface{}) (err error) { *u = uuid(src.([]uint8)); return }
func (u uuid) KeyName() string { return "uuid" }

// StudentID 필드에서 사용할 사용자 정의 타입
type studentID string
func StudentID(s string) studentID { return studentID(s) }
func (si studentID) Value() (driver.Value, error) { return string(si), nil }
func (si *studentID) Scan(src interface{}) (err error) { *si = studentID(src.([]uint8)); return }
func (si studentID) KeyName() string { return "student_id" }

// StudentPW 필드에서 사용할 사용자 정의 타입
type studentPW string
func StudentPW(s string) studentPW { return studentPW(s) }
func (sp studentPW) Value() (driver.Value, error) { return string(sp), nil }
func (sp *studentPW) Scan(src interface{}) (err error) { *sp = studentPW(src.([]uint8)); return }
func (sp studentPW) KeyName() string { return "student_pw" }

// TeacherID 필드에서 사용할 사용자 정의 타입
type teacherID string
func TeacherID(s string) teacherID { return teacherID(s) }
func (ti teacherID) Value() (driver.Value, error) { return string(ti), nil }
func (ti *teacherID) Scan(src interface{}) (err error) { *ti = teacherID(src.([]uint8)); return }
func (ti teacherID) KeyName() string { return "teacher_id" }

// TeacherPW 필드에서 사용할 사용자 정의 타입
type teacherPW string
func TeacherPW(s string) teacherPW { return teacherPW(s) }
func (tp teacherPW) Value() (driver.Value, error) { return string(tp), nil }
func (tp *teacherPW) Scan(src interface{}) (err error) { *tp = teacherPW(src.([]uint8)); return }
func (tp teacherPW) KeyName() string { return "teacher_pw" }

// ParentID 필드에서 사용할 사용자 정의 타입
type parentID string
func ParentID(s string) parentID { return parentID(s) }
func (pi parentID) Value() (driver.Value, error) { return string(pi), nil }
func (pi *parentID) Scan(src interface{}) (err error) { *pi = parentID(src.([]uint8)); return }
func (pi parentID) KeyName() string { return "parent_id" }

// ParentPW 필드에서 사용할 사용자 정의 타입
type parentPW string
func ParentPW(s string) parentPW { return parentPW(s) }
func (pp parentPW) Value() (driver.Value, error) { return string(pp), nil }
func (pp *parentPW) Scan(src interface{}) (err error) { *pp = parentPW(src.([]uint8)); return }
func (pp parentPW) KeyName() string { return "parent_pw" }

// StudentUUID 필드에서 사용할 사용자 정의 타입
type studentUUID string
func StudentUUID(s string) studentUUID { return studentUUID(s) }
func (su studentUUID) Value() (driver.Value, error) { return string(su), nil }
func (su *studentUUID) Scan(src interface{}) (err error) { *su = studentUUID(src.([]uint8)); return }
func (su studentUUID) KeyName() string { return "student_uuid" }

// TeacherUUID 필드에서 사용할 사용자 정의 타입
type teacherUUID string
func TeacherUUID(s string) teacherUUID { return teacherUUID(s) }
func (tu teacherUUID) Value() (driver.Value, error) { return string(tu), nil }
func (tu *teacherUUID) Scan(src interface{}) (err error) { *tu = teacherUUID(src.([]uint8)); return }
func (tu teacherUUID) KeyName() string { return "teacher_uuid" }

// ParentUUID 필드에서 사용할 사용자 정의 타입
type parentUUID string
func ParentUUID(s string) parentUUID { return parentUUID(s) }
func (pu parentUUID) Value() (value driver.Value, err error) {
	value = string(pu)
	if value == "" { value = nil }
	return
}
func (pu *parentUUID) Scan(src interface{}) (err error) { *pu = parentUUID(src.([]uint8)); return }
func (pu parentUUID) KeyName() string { return "parent_uuid" }
func (pu parentUUID) NullReplaceValue() string { return nullReplaceValueForParentUUID }

// Name 필드에서 사용할 사용자 정의 타입
type name string
func Name(s string) name { return name(s) }
func (n name) Value() (driver.Value, error) { return string(n), nil }
func (n *name) Scan(src interface{}) (err error) { *n = name(src.([]uint8)); return }
func (n name) KeyName() string { return "name" }

// PhoneNumber 필드에서 사용할 사용자 정의 타입
type phoneNumber string
func PhoneNumber(s string) phoneNumber { return phoneNumber(s) }
func (pn phoneNumber) Value() (driver.Value, error) { return string(pn), nil }
func (pn *phoneNumber) Scan(src interface{}) (err error) { *pn = phoneNumber(src.([]uint8)); return }
func (pn phoneNumber) KeyName() string { return "phone_number" }

// ProfileURI 필드에서 사용할 사용자 정의 타입
type profileURI string
func ProfileURI(s string) profileURI { return profileURI(s) }
func (pu profileURI) Value() (driver.Value, error) { return string(pu), nil }
func (pu *profileURI) Scan(src interface{}) (err error) { *pu = profileURI(src.([]uint8)); return }
func (pu profileURI) KeyName() string { return "profile_uri" }