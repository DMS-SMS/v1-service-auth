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

// 사전에 거를 수 없었던 상황에 대한 오류는 mysql 에러로 반환, 그렇지 않으면 X -> 500으로 처리
func (si *StudentInform) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(si); err != nil {
		return
	}

	query := tx.Where("grade = ? AND class = ? AND student_number = ?", si.Grade, si.Class, si.StudentNumber).Find(&StudentInform{})
	if query.RowsAffected != 0 {
		// number와 같은 key들 상수로 선언 및 관리 필요
		return mysqlerr.DuplicateEntry("number", fmt.Sprintf("%d-%d-%02d", si.Grade.value, si.Class.value, si.StudentNumber.value))
	}

	return nil
}

func (ti *TeacherInform) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(ti); err != nil {
		return
	}

	query := tx.Where("grade = ? AND class = ?", ti.Grade, ti.Class).Find(&TeacherInform{})
	if query.RowsAffected != 0 {
		// class와 같은 entry들 상수로 선언 및 관리 필요
		// DB의 다형성을 위한 에러 코드 생성 함수 추상화 필요 (나중에)
		return mysqlerr.DuplicateEntry("class", fmt.Sprintf("%d-%d", ti.Grade.value, ti.Class.value))
	}

	return nil
}

func (pi *ParentInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pi)
}