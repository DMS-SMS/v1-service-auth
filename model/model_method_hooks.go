package model

import (
	"auth/model/validate"
	"auth/tool/mysqlerr"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	emptyString = ""
	emptyInt = 0
)

const (
	validStudentUUID = "student-111111111111"
	validGrade = 2
	validClass = 2
	validStudentNumber = 7
	validName = "박진홍"
	validPhoneNumber = "01088378347"
	validProfileURI = "example.com/profiles/student-111111111111"
)

func (sa *StudentAuth) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(sa)
}

func (ta *TeacherAuth) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(ta)
}

func (pa *ParentAuth) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pa)
}

// 사전에 거를 수 없었던 상황에 대한 오류는 mysql 에러로 반환, 그렇지 않으면 X -> 500으로 처리
func (si *StudentInform) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(si); err != nil {
		return
	}

	query := tx.Where("grade = ? AND class = ? AND student_number = ?", si.Grade, si.Class, si.StudentNumber).Find(&StudentInform{})
	if query.RowsAffected != 0 {
		// number와 같은 key들 상수로 선언 및 관리 필요
		err = mysqlerr.DuplicateEntry(si.StudentNumber.KeyName(), fmt.Sprintf("%d%d%02d", si.Grade, si.Class, si.StudentNumber))
	}

	return
}

func (ti *TeacherInform) BeforeCreate(tx *gorm.DB) (err error) {
	return validate.DBValidator.Struct(ti)
}

func (pi *ParentInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pi)
}