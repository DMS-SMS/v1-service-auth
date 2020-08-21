package model

import "github.com/jinzhu/gorm"

type ParentAuth struct {
	gorm.Model
	UUID     string `gorm:"PRIMARY_KEY;Type:char(19);UNIQUE;INDEX"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	ParentId string `gorm:"Type:varchar(20);NOT NULL;UNIQUE"`  // 4~20자 사이
	ParentPw string `gorm:"Type:varchar(100);NOT NULL"`
}

type ParentInform struct {
	gorm.Model
	ParentUUID  string `gorm:"Type:char(19);UNIQUE;NOT NULL"` // 형식 => 'parent-' + 12자리 랜덤 수
	Name        string `gorm:"Type:char(4);NOT NULL"`  // 2~4자 사이
	PhoneNumber string `gorm:"Type:char(11);UNIQUE;NOT NULL"` // 11자
}