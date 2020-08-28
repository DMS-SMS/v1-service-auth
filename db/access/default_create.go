package access

import (
	"auth/model"
)

func (d *Default) CreateStudentAuth(auth *model.StudentAuth) (*model.StudentAuth, error) {
	result := d.tx.Create(auth)
	if auth, ok := result.Value.(*model.StudentAuth); ok {
		return auth, result.Error
	}
	if result.Error == nil {
		result.Error = StudentAuthAssertionError
	}
	return nil, result.Error
}

func (d *Default) CreateTeacherAuth(auth *model.TeacherAuth) (*model.TeacherAuth, error) {
	result := d.tx.Create(auth)
	if auth, ok := result.Value.(*model.TeacherAuth); ok {
		return auth, result.Error
	}
	if result.Error == nil {
		result.Error = TeacherAuthAssertionError
	}
	return nil, result.Error
}

func (d *Default) CreateParentAuth(auth *model.ParentAuth) (*model.ParentAuth, error) {
	result := d.tx.Create(auth)
	if auth, ok := result.Value.(*model.ParentAuth); ok {
		return auth, result.Error
	}
	if result.Error == nil {
		result.Error = ParentAuthAssertionError
	}
	return nil, result.Error
}
