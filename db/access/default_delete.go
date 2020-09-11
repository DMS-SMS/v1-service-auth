package access

import (
	"auth/model"
)

func (d *_default) DeleteStudentAuth(uuid string) (err error) {
	err = d.tx.Delete(&model.StudentAuth{UUID: model.UUID(uuid)}).Error
	return
}

func (d *_default) DeleteTeacherAuth(uuid string) (err error) {
	err = d.tx.Delete(&model.TeacherAuth{UUID: model.UUID(uuid)}).Error
	return
}

func (d *_default) DeleteParentAuth(uuid string) (err error) {
	err = d.tx.Delete(&model.ParentAuth{UUID: model.UUID(uuid)}).Error
	return
}