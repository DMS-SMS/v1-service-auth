package access

import (
	"auth/model"
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