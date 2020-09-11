package errors

import "errors"

var (
	StudentUUIDCannotBeChanged = errors.New("student uuid cannot be changed")
	TeacherUUIDCannotBeChanged = errors.New("teacher uuid cannot be changed")
	ParentUUIDCannotBeChanged = errors.New("parent uuid cannot be changed")
)