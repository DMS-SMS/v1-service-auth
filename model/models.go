package model

import (
	"github.com/jinzhu/gorm"
)
// ID ASCII만? 이름 한글 유니코드만?

// 학생 계정 테이블
type StudentAuth struct {
	gorm.Model
	UUID       uuid       `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=student,len=20"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	StudentID  studentID  `gorm:"Type:varchar(20);NOT NULL" validate:"min=4,max=20,ascii"`        // 4~20자 사이
	StudentPW  studentPW  `gorm:"Type:varchar(100);NOT NULL"`
	ParentUUID parentUUID `gorm:"Type:char(19);" validate:"uuid=parent,len=19"` 						 // 형식 => 'parent-' + 12자리 랜덤 수 (19자) 만일의 경우를 대비해서 NOT NULL 삭제
}

// 학생 사용자 정보 테이블
type StudentInform struct {
	gorm.Model
	StudentUUID   studentUUID   `gorm:"Type:char(20);UNIQUE;NOT NULL" validate:"uuid=student,len=20"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	Grade         grade         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~3"`                // 1~3 사이 값
	Class         class         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~4"`                // 1~4 사이 값
	StudentNumber studentNumber `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~21"`               // 1~21 사이 값
	Name          name          `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"`       // 2~4자 사이 한글
	PhoneNumber   phoneNumber   `gorm:"Type:char(11);UNIQUE;NOT NULL" validate:"len=11,phone_number"` // 11자
	ProfileURI    profileURI    `gorm:"Type:varchar(150);UNIQUE;NOT NULL"`                            // 제약 조건 나중에 추가 예정
}

// 선생님 계정 테이블
type TeacherAuth struct {
	gorm.Model
	UUID      uuid      `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=teacher,len=20"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	TeacherID teacherID `gorm:"varchar(20);NOT NULL" validate:"min=4,max=20,ascii"`             // 4~20자 사이
	TeacherPW teacherPW `gorm:"varchar(100):NOT NULL"`
}

// 선생님 사용자 정보 테이블
type TeacherInform struct {
	gorm.Model
	TeacherUUID teacherUUID `gorm:"Type:char(20);UNIQUE;NOT NULL'" validate:"uuid=teacher,len=20"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	Name        name        `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"`        // 2~4자 사이 한글
	Grade       grade       `gorm:"Type:tinyint(1);" validate:"range=0~3"`                         // in (1~3)
	Class       class       `gorm:"Type:tinyint(1);" validate:"range=0~4"`                         // in (1~4)
	PhoneNumber phoneNumber `gorm:"Type:char(11);UNIQUE;NOT NULL" validate:"len=11,phone_number"`  // 11자
}

// 부모님 계정 테이블
type ParentAuth struct {
	gorm.Model
	UUID     uuid     `gorm:"PRIMARY_KEY;Type:char(19);UNIQUE;INDEX" validate:"uuid=parent,len=19"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	ParentID parentID `gorm:"Type:varchar(20);NOT NULL" validate:"min=4,max=20,ascii"`       // 4~20자 사이
	ParentPW parentPW `gorm:"Type:varchar(100);NOT NULL"`
}

// 부모님 사용자 정보 테이블
type ParentInform struct {
	gorm.Model
	ParentUUID  parentUUID  `gorm:"Type:char(19);UNIQUE;NOT NULL" validate:"uuid=parent,len=19"`  // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	Name        name        `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"`       // 2~4자 사이 한글
	PhoneNumber phoneNumber `gorm:"Type:char(11);UNIQUE;NOT NULL" validate:"len=11,phone_number"` // 11자
}

// 관리자 계정 테이블
type AdminAuth struct {
	gorm.Model
	UUID    uuid    `gorm:"Type:char(18);UNIQUE;NOT NULL" validate:"uuid=admin,len=18"`
	AdminID adminID `gorm:"varchar(20);NOT NULL;UNIQUE" validate:"min=4,max=20,ascii"`
	AdminPW adminPW `gorm:"varchar(100):NOT NULL;"`
}