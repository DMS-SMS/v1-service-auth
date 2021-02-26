package access

import (
	"auth/model"
)

func (d *_default) DeleteStudentAuth(uuid string) (err error) {
	err = d.tx.Where("uuid = ?", uuid).Delete(&model.StudentAuth{}).Error
	return
}

func (d *_default) DeleteTeacherAuth(uuid string) (err error) {
	err = d.tx.Where("uuid = ?", uuid).Delete(&model.TeacherAuth{}).Error
	return
}

func (d *_default) DeleteParentAuth(uuid string) (err error) {
	err = d.tx.Where("uuid = ?", uuid).Delete(&model.ParentAuth{}).Error
	return
}

func (d *_default) DeleteStudentInform(studentUUID string) (err error) {
	err = d.tx.Where("student_uuid = ?", studentUUID).Delete(&model.StudentInform{}).Error
	return
}

func (d *_default) DeleteTeacherInform(teacherUUID string) (err error) {
	err = d.tx.Where("teacher_uuid = ?", teacherUUID).Delete(&model.TeacherInform{}).Error
	return
}

func (d *_default) DeleteParentInform(parentUUID string) (err error) {
	err = d.tx.Where("parent_uuid = ?", parentUUID).Delete(&model.ParentInform{}).Error
	return
}

func (d *_default) DeleteUnsignedStudent(authCode int64) (err error) {
	err = d.tx.Where("auth_code = ?", authCode).Delete(&model.UnsignedStudent{}).Error
	return
}
