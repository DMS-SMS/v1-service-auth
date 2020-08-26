package model

import "auth/model/validate"

func (sa *StudentAuth) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(sa)
}

func (ta *TeacherAuth) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(ta)
}

func (pa *ParentAuth) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pa)
}

func (si *StudentInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(si)
}

func (ti *TeacherInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(ti)
}

func (pi *ParentInform) BeforeCreate() (err error) {
	return validate.DBValidator.Struct(pi)
}