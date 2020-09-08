package access

import "auth/model"

func (d *_default) GetStudentAuthWithID(id string) (auth *model.StudentAuth, err error) {
	auth = new(model.StudentAuth)
	err = d.tx.Where("student_id = ?", id).Find(&auth).Error
	return
}

func (d *_default) GetTeacherAuthWithID(id string) (auth *model.TeacherAuth, err error) {
	auth = new(model.TeacherAuth)
	err = d.tx.Where("teacher_id = ?", id).Find(&auth).Error
	return
}

func (d *_default) GetParentAuthWithID(id string) (auth *model.ParentAuth, err error) {
	auth = new(model.ParentAuth)
	err = d.tx.Where("parent_id = ?", id).Find(&auth).Error
	return
}