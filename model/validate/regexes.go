package validate

import "regexp"

const (
	adminUUIDRegexString = "^admin-\\d{12}"
	studentUUIDRegexString = "^student-\\d{12}"
	teacherUUIDRegexString = "^teacher-\\d{12}"
	parentUUIDRegexString = "^parent-\\d{12}"
	phoneNumberRegexString = "\\d{11}"
)

var (
	adminUUIDRegex = regexp.MustCompile(adminUUIDRegexString)
	studentUUIDRegex = regexp.MustCompile(studentUUIDRegexString)
	teacherUUIDRegex = regexp.MustCompile(teacherUUIDRegexString)
	parentUUIDRegex = regexp.MustCompile(parentUUIDRegexString)
	phoneNumberRegex = regexp.MustCompile(phoneNumberRegexString)
)