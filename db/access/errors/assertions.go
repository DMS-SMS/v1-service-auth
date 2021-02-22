package errors

import (
	"errors"
	"fmt"
)

const (
	InterfaceAssertionErrorFormat = "interface conversion error, interface: %s"
)

var (
	StudentAuthAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.StudentAuth"))
	TeacherAuthAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.TeacherAuth"))
	ParentAuthAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.ParentAuth"))
	StudentInformAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.StudentInform"))
	TeacherInformAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.TeacherInform"))
	ParentInformAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.ParentInform"))
	UnsignedStudentAssertionError = errors.New(fmt.Sprintf(InterfaceAssertionErrorFormat, "*model.UnsignedStudent"))
)