package model

func (sa *StudentAuth) ExceptGormModel() *StudentAuth  {
	return &StudentAuth{
		UUID:       sa.UUID,
		StudentId:  sa.StudentId,
		StudentPw:  sa.StudentPw,
		ParentUUID: sa.ParentUUID,
	}
}

func (ta *TeacherAuth) ExceptGormModel() *TeacherAuth  {
	return &TeacherAuth{
		UUID:      ta.UUID,
		TeacherId: ta.TeacherId,
		TeacherPw: ta.TeacherPw,
	}
}

func (pa *ParentAuth) ExceptGormModel() *ParentAuth  {
	return &ParentAuth{
		UUID:     pa.UUID,
		ParentId: pa.ParentId,
		ParentPw: pa.ParentPw,
	}
}
