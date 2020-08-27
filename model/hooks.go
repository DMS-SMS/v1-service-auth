package model

import (
	"auth/model/validate"
	"auth/tool/mysqlerr"
	"fmt"
	"github.com/jinzhu/gorm"
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

func (si *StudentInform) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(si); err != nil {
		return
	}

	query := tx.Where("grade = ? AND class = ? AND student_number = ?", si.Grade, si.Class, si.StudentNumber).Find(&StudentInform{})
	if query.RowsAffected != 0 {
		// number와 같은 entry들 상수로 선언 및 관리 필요
		return mysqlerr.DuplicateEntry("number", fmt.Sprintf("%d%d%02d", si.Grade.value, si.Class.value, si.StudentNumber.value))
	}

	return nil
}

func (ti *TeacherInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(ti)
	// 학년 반 중복 확인 추가
}

func (pi *ParentInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pi)
}