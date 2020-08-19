package model

import "github.com/jinzhu/gorm"

type StudentAuth struct {
	gorm.Model
	UUID       string `gorm:"PRIMARY_KEY;Type:char(20);INDEX"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	StudentId  string `gorm:"Type:varchar(20);NOT NULL"`  // 4~20자 사이
	StudentPw  string `gorm:"Type:varchar(100);NOT NULL"`
	ParentUUID string `gorm:"Type:char(19);NOT NULL"`	  // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
}

type StudentInform struct {
	gorm.Model
	StudentUUID   string `gorm:"Type:char(20);NOT NULL"`   // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	Grade         uint   `gorm:"Type:tinyint(1);NOT NULL"` // in (1~3)
	Class         uint   `gorm:"Type:tinyint(1);NOT NULL"` // in (1~4)
	StudentNumber uint   `gorm:"Type:tinyint(1);NOT NULL"` // in (1~21)
	Name          string `gorm:"Type:varchar(4);NOT NULL"` // 2~4자 사이 한글
	PhoneNumber   string `gorm:"Type:char(11);NOT NULL"`   // 11자
	ProfileUri    string `gorm:"Type:varchar(150);NOT NULL"`
}