package model

import "auth/tool/arraylist"

var (
	gradeAvailableValues = []int64{1, 2, 3}
	GradeValuesList = arraylist.NewWithInt64(gradeAvailableValues...)
	classAvailableValues = []int64{1, 2, 3, 4}
	ClassValuesList = arraylist.NewWithInt64(classAvailableValues...)
)