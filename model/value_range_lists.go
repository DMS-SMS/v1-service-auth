package model

import "auth/tool/arraylist"

var (
	gradeAvailableValues = []int64{1, 2, 3}
	GradeValuesList = arraylist.NewWithInt64(gradeAvailableValues...)

	classAvailableValues = []int64{1, 2, 3, 4}
	ClassValuesList = arraylist.NewWithInt64(classAvailableValues...)

	studentNumberAvailableValues = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
	StudentNumberValuesList = arraylist.NewWithInt64(studentNumberAvailableValues...)
)