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
