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

	validStudentUUID = "student-111111111111"
	validTeacherUUID = "teacher-111111111111"
	validParentUUID = "parent-111111111111"
	validGrade = 2
	validClass = 2
	validStudentNumber = 7
	validName = "박진홍"
	validPhoneNumber = "01088378347"
	validProfileURI = "example.com/profiles/student-111111111111"
)

func (sa *StudentAuth) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(sa); err != nil {
		return
	}

	query := tx.Where("student_id = ?", sa.StudentID).Find(&StudentAuth{})
	if query.RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(StudentAuthInstance.StudentID.KeyName(), string(sa.StudentID))
	}
	return
}

func (ta *TeacherAuth) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(ta); err != nil {
		return
	}

	query := tx.Where("teacher_id = ?", ta.TeacherID).Find(&TeacherAuth{})
	if query.RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(TeacherAuthInstance.TeacherID.KeyName(), string(ta.TeacherID))
	}
	return
}

func (pa *ParentAuth) BeforeCreate(tx *gorm.DB) (err error) {
	if err = validate.DBValidator.Struct(pa); err != nil {
		return
	}

	query := tx.Where("parent_id = ?", pa.ParentID).Find(&ParentAuth{})
	if query.RowsAffected != 0 {
		err = mysqlerr.DuplicateEntry(ParentAuthInstance.ParentID.KeyName(), string(pa.ParentID))
	}
	return
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

func (si *StudentInform) BeforeUpdate(tx *gorm.DB) (err error) {
	informForValidate := si.DeepCopy()

	if informForValidate.StudentUUID == emptyString { informForValidate.StudentUUID = validStudentUUID }
	if informForValidate.Grade == emptyInt          { informForValidate.Grade = validGrade }
	if informForValidate.Class == emptyInt          { informForValidate.Class = validClass }
	if informForValidate.StudentNumber == emptyInt  { informForValidate.StudentNumber = validStudentNumber }
	if informForValidate.Name == emptyString        { informForValidate.Name = validName }
	if informForValidate.PhoneNumber == emptyString { informForValidate.PhoneNumber = validPhoneNumber }
	if informForValidate.ProfileURI == emptyString  { informForValidate.ProfileURI = validProfileURI }

	if err = validate.DBValidator.Struct(informForValidate); err != nil {
		return
	}

	if si.Grade != emptyInt && si.Class != emptyInt && si.StudentNumber != emptyInt {
		studentNumberTable := tx.Where("grade = ? AND class = ? AND student_number = ?", si.Grade, si.Class, si.StudentNumber).Find(&StudentInform{})
		if studentNumberTable.RowsAffected != 0 {
			err = mysqlerr.DuplicateEntry(si.StudentNumber.KeyName(), fmt.Sprintf("%d%d%02d", si.Grade, si.Class, si.StudentNumber))
		}
	}
	return
}

func (ti *TeacherInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(ti)
}

func (ti *TeacherInform) BeforeUpdate() (err error) {
	informForValidate := ti.DeepCopy()

	if informForValidate.TeacherUUID == emptyString { informForValidate.TeacherUUID = validTeacherUUID }
	if informForValidate.Grade == emptyInt          { informForValidate.Grade = validGrade }
	if informForValidate.Class == emptyInt          { informForValidate.Class = validClass }
	if informForValidate.Name == emptyString        { informForValidate.Name = validName }
	if informForValidate.PhoneNumber == emptyString { informForValidate.PhoneNumber = validPhoneNumber }

	return validate.DBValidator.Struct(informForValidate)
}

func (pi *ParentInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pi)
}

func (pi *ParentInform) BeforeUpdate() (err error) {
	informForValidate := pi.DeepCopy()

	if informForValidate.ParentUUID == emptyString  { informForValidate.ParentUUID = validParentUUID }
	if informForValidate.Name == emptyString        { informForValidate.Name = validName }
	if informForValidate.PhoneNumber == emptyString { informForValidate.PhoneNumber = validPhoneNumber }

	return validate.DBValidator.Struct(informForValidate)
}
