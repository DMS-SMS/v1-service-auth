package model

import (
	"github.com/jinzhu/gorm"
)

// 학생 계정 테이블
type StudentAuth struct {
	gorm.Model
	UUID       string `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=student"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	StudentId  string `gorm:"Type:varchar(20);NOT NULL;UNIQUE"`                               // 4~20자 사이
	StudentPw  string `gorm:"Type:varchar(100);NOT NULL"`
	ParentUUID string `gorm:"Type:char(19);"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자) 만일의 경우를 대비해서 NOT NULL 삭제
}

// 학생 사용자 정보 테이블
type StudentInform struct {
	gorm.Model
	StudentUUID   string `gorm:"Type:char(20);UNIQUE;NOT NULL"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	Grade         grade  `gorm:"Type:tinyint(1);NOT NULL"`
	Class         uint   `gorm:"Type:tinyint(1);NOT NULL"`      // in (1~4)
	StudentNumber uint   `gorm:"Type:tinyint(1);NOT NULL"`      // in (1~21)
	Name          string `gorm:"Type:varchar(4);NOT NULL"`      // 2~4자 사이 한글
	PhoneNumber   string `gorm:"Type:char(11);UNIQUE;NOT NULL"` // 11자
	ProfileUri    string `gorm:"Type:varchar(150);UNIQUE;NOT NULL"`
}

// ---

// 선생님 계정 테이블
type TeacherAuth struct {
	gorm.Model
	UUID      string `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=teacher"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	TeacherId string `gorm:"varchar(20);NOT NULL;UNIQUE"`                                    // 4~20자 사이
	TeacherPw string `gorm:"varchar(100):NOT NULL;"`
}

// 선생님 사용자 정보 테이블
type TeacherInform struct {
	gorm.Model
	TeacherUUID string `gorm:"Type:char(20);UNIQUE;NOT NULL'"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	Name        string `gorm:"Type:varchar(20);NOT NULL"`      // 2~4자 사이 한글
	Grade       uint   `gorm:"Type:tinyint(1);"`               // in (1~3)
	Class       uint   `gorm:"Type:tinyint(1);"`               // in (1~4)
	PhoneNumber string `gorm:"Type:char(11);UNIQUE;NOT NULL"`  // 11자
}

// ---

// 부모님 계정 테이블
type ParentAuth struct {
	gorm.Model
	UUID     string `gorm:"PRIMARY_KEY;Type:char(19);UNIQUE;INDEX" validate:"uuid=parent"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	ParentId string `gorm:"Type:varchar(20);NOT NULL;UNIQUE"`                              // 4~20자 사이
	ParentPw string `gorm:"Type:varchar(100);NOT NULL"`
}

// 부모님 사용자 정보 테이블
type ParentInform struct {
	gorm.Model
	ParentUUID  string `gorm:"Type:char(19);UNIQUE;NOT NULL"` // 형식 => 'parent-' + 12자리 랜덤 수
	Name        string `gorm:"Type:char(4);NOT NULL"`         // 2~4자 사이
	PhoneNumber string `gorm:"Type:char(11);UNIQUE;NOT NULL"` // 11자
}