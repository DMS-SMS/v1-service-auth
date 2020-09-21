package db

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.LogMode(false)

	//db.DropTableIfExists(&model.AdminAuth{})
	//db.DropTableIfExists(&model.StudentInform{})
	//db.DropTableIfExists(&model.ParentInform{})
	//db.DropTableIfExists(&model.TeacherInform{})
	//db.DropTableIfExists(&model.StudentAuth{})
	//db.DropTableIfExists(&model.ParentAuth{})
	//db.DropTableIfExists(&model.TeacherAuth{})

	if !db.HasTable(&model.AdminAuth{}) {
		db.CreateTable(&model.AdminAuth{})
	}
	if !db.HasTable(&model.StudentAuth{}) {
		db.CreateTable(&model.StudentAuth{})
	}
	if !db.HasTable(&model.StudentInform{}) {
		db.CreateTable(&model.StudentInform{})
	}
	if !db.HasTable(&model.ParentAuth{}) {
		db.CreateTable(&model.ParentAuth{})
	}
	if !db.HasTable(&model.ParentInform{}) {
		db.CreateTable(&model.ParentInform{})
	}
	if !db.HasTable(&model.TeacherAuth{}) {
		db.CreateTable(&model.TeacherAuth{})
	}
	if !db.HasTable(&model.TeacherInform{}) {
		db.CreateTable(&model.TeacherInform{})
	}

	db.AutoMigrate(&model.AdminAuth{}, &model.StudentAuth{}, &model.StudentInform{}, &model.ParentAuth{}, &model.ParentInform{}, &model.TeacherAuth{}, &model.TeacherInform{})
	db.Model(&model.StudentAuth{}).AddForeignKey("parent_uuid", "parent_auths(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&model.StudentInform{}).AddForeignKey("student_uuid", "student_auths(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&model.TeacherInform{}).AddForeignKey("teacher_uuid", "teacher_auths(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&model.ParentInform{}).AddForeignKey("parent_uuid", "parent_auths(uuid)", "RESTRICT", "RESTRICT")

	// 데이터 무결성 제약조건 추가 필요
}
