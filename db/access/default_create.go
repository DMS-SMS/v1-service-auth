package access

import (
	"auth/db/access/errors"
	"auth/model"
)

func (d *_default) CreateStudentAuth(auth *model.StudentAuth) (*model.StudentAuth, error) {
	result := d.tx.Create(auth)
	if auth, ok := result.Value.(*model.StudentAuth); ok {
		return auth, result.Error
	}
	if result.Error == nil {
		result.Error = errors.StudentAuthAssertionError
	}
	return nil, result.Error
}

func (d *_default) CreateTeacherAuth(auth *model.TeacherAuth) (*model.TeacherAuth, error) {
	result := d.tx.Create(auth)
	if auth, ok := result.Value.(*model.TeacherAuth); ok {
		return auth, result.Error
	}
	if result.Error == nil {
		result.Error = errors.TeacherAuthAssertionError
	}
	return nil, result.Error
}

func (d *_default) CreateParentAuth(auth *model.ParentAuth) (*model.ParentAuth, error) {
	result := d.tx.Create(auth)
	if auth, ok := result.Value.(*model.ParentAuth); ok {
		return auth, result.Error
	}
	if result.Error == nil {
		result.Error = errors.ParentAuthAssertionError
	}
	return nil, result.Error
}

func (d *_default) CreateStudentInform(inform *model.StudentInform) (*model.StudentInform, error) {
	result := d.tx.Create(inform)
	if inform, ok := result.Value.(*model.StudentInform); ok {
		return inform, result.Error
	}
	if result.Error == nil {
		result.Error = errors.StudentInformAssertionError
	}
	return nil, result.Error
}

func (d *_default) CreateTeacherInform(inform *model.TeacherInform) (*model.TeacherInform, error) {
	result := d.tx.Create(inform)
	if inform, ok := result.Value.(*model.TeacherInform); ok {
		return inform, result.Error
	}
	if result.Error == nil {
		result.Error = errors.TeacherInformAssertionError
	}
	return nil, result.Error
}

func (d *_default) CreateParentInform(inform *model.ParentInform) (*model.ParentInform, error) {
	result := d.tx.Create(inform)
	if inform, ok := result.Value.(*model.ParentInform); ok {
		return inform, result.Error
	}
	if result.Error == nil {
		result.Error = errors.ParentInformAssertionError
	}
	return nil, result.Error
}

func (d *_default) AddUnsignedStudent(student *model.UnsignedStudent) (*model.UnsignedStudent, error) {
	result := d.tx.Create(student)
	if inform, ok := result.Value.(*model.UnsignedStudent); ok {
		return inform, result.Error
	}
	if result.Error == nil {
		result.Error = errors.UnsignedStudentAssertionError
	}
	return nil, result.Error
}
