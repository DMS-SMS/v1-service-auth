package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

const (
	emptyString = ""
	emptyInt = 0
)

func (d *_default) GetStudentAuthWithID(studentID string) (auth *model.StudentAuth, err error) {
	auth = new(model.StudentAuth)
	err = d.tx.Where("student_id = ?", studentID).Find(&auth).Error
	return
}

func (d *_default) GetTeacherAuthWithID(teacherID string) (auth *model.TeacherAuth, err error) {
	auth = new(model.TeacherAuth)
	err = d.tx.Where("teacher_id = ?", teacherID).Find(&auth).Error
	return
}

func (d *_default) GetParentAuthWithID(parentID string) (auth *model.ParentAuth, err error) {
	auth = new(model.ParentAuth)
	err = d.tx.Where("parent_id = ?", parentID).Find(&auth).Error
	return
}

func (d *_default) GetStudentUUIDsWithInform(inform *model.StudentInform) (uuidArr []string, err error) {
	cascadeTx := d.tx.New()

	if inform.StudentUUID != emptyString {
		cascadeTx = cascadeTx.Where("student_uuid = ?", inform.StudentUUID)
	}
	if inform.Grade != emptyInt {
		cascadeTx = cascadeTx.Where("grade = ?", inform.Grade)
	}
	if inform.Class != emptyInt {
		cascadeTx = cascadeTx.Where("class = ?", inform.Class)
	}
	if inform.StudentNumber != emptyInt {
		cascadeTx = cascadeTx.Where("class = ?", inform.StudentNumber)
	}
	if inform.Name != emptyString {
		cascadeTx = cascadeTx.Where("name = ?", inform.Name)
	}
	if inform.PhoneNumber != emptyString {
		cascadeTx = cascadeTx.Where("phone_number = ?", inform.PhoneNumber)
	}
	if inform.ProfileURI != emptyString {
		cascadeTx = cascadeTx.Where("profile_uri = ?", inform.ProfileURI)
	}

	informs := make([]*model.StudentInform, 1, 3)
	err = cascadeTx.Find(&informs).Error

	if len(informs) == 0 {
		err = gorm.ErrRecordNotFound
	}

	for _, inform := range informs {
		uuidArr = append(uuidArr, string(inform.StudentUUID))
	}
	return
}

func (d *_default) GetTeacherUUIDsWithInform(inform *model.TeacherInform) (uuidArr []string, err error) {
	cascadeTx := d.tx.New()

	if inform.TeacherUUID != emptyString {
		cascadeTx = cascadeTx.Where("teacher_uuid = ?", inform.TeacherUUID)
	}
	if inform.Grade != emptyInt {
		cascadeTx = cascadeTx.Where("grade = ?", inform.Grade)
	}
	if inform.Class != emptyInt {
		cascadeTx = cascadeTx.Where("class = ?", inform.Class)
	}
	if inform.Name != emptyString {
		cascadeTx = cascadeTx.Where("name = ?", inform.Name)
	}
	if inform.PhoneNumber != emptyString {
		cascadeTx = cascadeTx.Where("phone_number = ?", inform.PhoneNumber)
	}

	informs := make([]*model.TeacherInform, 1, 3)
	err = cascadeTx.Find(&informs).Error

	if len(informs) == 0 {
		err = gorm.ErrRecordNotFound
	}

	for _, inform := range informs {
		uuidArr = append(uuidArr, string(inform.TeacherUUID))
	}
	return
}

func (d *_default) GetParentUUIDsWithInform(inform *model.ParentInform) (uuidArr []string, err error) {
	cascadeTx := d.tx.New()

	if inform.ParentUUID != emptyString {
		cascadeTx = cascadeTx.Where("parent_uuid = ?", inform.ParentUUID)
	}
	if inform.Name != emptyString {
		cascadeTx = cascadeTx.Where("name = ?", inform.Name)
	}
	if inform.PhoneNumber != emptyString {
		cascadeTx = cascadeTx.Where("phone_number = ?", inform.PhoneNumber)
	}

	informs := make([]*model.ParentInform, 1, 3)
	err = cascadeTx.Find(&informs).Error

	if len(informs) == 0 {
		err = gorm.ErrRecordNotFound
	}

	for _, inform := range informs {
		uuidArr = append(uuidArr, string(inform.ParentUUID))
	}
	return
}