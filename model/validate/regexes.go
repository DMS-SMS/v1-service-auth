package validate

import "regexp"

const (
	studentUUIDRegexString = "^student-\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d"
	teacherUUIDRegexString = "^teacher-\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d"
	parentUUIDRegexString = "^parent-\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d"
	phoneNumberRegexString = "\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d"
)

var (
	studentUUIDRegex = regexp.MustCompile(studentUUIDRegexString)
	teacherUUIDRegex = regexp.MustCompile(teacherUUIDRegexString)
	parentUUIDRegex = regexp.MustCompile(parentUUIDRegexString)
	phoneNumberRegex = regexp.MustCompile(phoneNumberRegexString)
)