package model

import (
	"github.com/jinzhu/gorm"
)

type TeacherAuth struct {
	gorm.Model
	UUID      string `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=teacher"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	TeacherId string `gorm:"varchar(20);NOT NULL;UNIQUE"` // 4~20자 사이
	TeacherPw string `gorm:"varchar(100):NOT NULL;"`
}

type TeacherInform struct {
	gorm.Model
	TeacherUUID string `gorm:"Type:char(20);UNIQUE;NOT NULL'"`   // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	Name        string `gorm:"Type:varchar(20);NOT NULL"` // 2~4자 사이 한글
	Grade       uint   `gorm:"Type:tinyint(1);"`  // in (1~3)
	Class       uint   `gorm:"Type:tinyint(1);"` // in (1~4)
	PhoneNumber string `gorm:"Type:char(11);UNIQUE;NOT NULL"`  // 11자
}