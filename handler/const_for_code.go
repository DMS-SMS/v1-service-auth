package handler

const (
	CodeStudentIDDuplicate = -101
	CodeParentUUIDNoExist = -102
	CodeStudentNumberDuplicate = -103
	CodeStudentPhoneNumberDuplicate = -104

	CodeTeacherIDDuplicate = -201
	CodeTeacherPhoneNumberDuplicate = -202

	CodeParentIDDuplicate = -301
	CodeParentPhoneNumberDuplicate = -302

	CodeStudentIDNoExist = -401
	CodeIncorrectStudentPWForLogin = -402

	CodeTeacherIDNoExist = -411
	CodeIncorrectTeacherPWForLogin = -411

	CodeParentIDNoExist = -421
	CodeIncorrectParentPWForLogin = -422

	CodeAdminIDNoExist = -431
	CodeIncorrectAdminPWForLogin = -432

	CodeStudentWithThatInformNoExist = -501
	CodeTeacherWithThatInformNoExist = -511
	CodeParentWithThatInformNoExist = -521

	CodeIncorrectStudentPWForChange = -701

	CodeIncorrectTeacherPWForChange = -801

	CodeIncorrectParentPWForChange = -901
)
