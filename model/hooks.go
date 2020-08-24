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
