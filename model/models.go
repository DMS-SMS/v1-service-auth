package model

import (
	"github.com/jinzhu/gorm"
)

// 학생 계정 테이블
type StudentAuth struct {
	gorm.Model
	UUID       uuid       `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=student,len=20"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	StudentID  studentID  `gorm:"Type:varchar(20);NOT NULL" validate:"min=4,max=20,ascii"`               // 4~20자 사이
	StudentPW  studentPW  `gorm:"Type:varchar(100);NOT NULL"`
	ParentUUID parentUUID `gorm:"Type:char(19);" validate:"uuid=parent"` 						         // 형식 => 'parent-' + 12자리 랜덤 수 (19자) 만일의 경우를 대비해서 NOT NULL 삭제
}

// 학생 사용자 정보 테이블
type StudentInform struct {
	gorm.Model
	StudentUUID   studentUUID   `gorm:"Type:char(20);UNIQUE;NOT NULL" validate:"uuid=student,len=20"` // 형식 => 'student-' + 12자리 랜덤 수 (20자)
	Grade         grade         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~3"`                // 1~3 사이 값
	Class         class         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~4"`                // 1~4 사이 값
	StudentNumber studentNumber `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~21"`               // 1~21 사이 값
	Name          name          `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"`       // 2~4자 사이 한글
	PhoneNumber   phoneNumber   `gorm:"Type:char(11);NOT NULL" validate:"len=11,phone_number"`        // 11자
	ProfileURI    profileURI    `gorm:"Type:varchar(150);NOT NULL"`                                   // 제약 조건 나중에 추가 예정
	ParentStatus  parentStatus  `gorm:"varchar(30);default:OK_CONN_OK_NOTIFY;NOT NULL"`
}

// 계정 생성 전 사전에 인증된 사용자 정보 테이블
type UnsignedStudent struct {
	gorm.Model
	AuthCode      authCode      `gorm:"Type:int(11);NOT NULL" validate:"range=100000~999999"`   // 6자리 숫자
	Grade         grade         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~3"`          // 1~3 사이 값
	Class         class         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~4"`          // 1~4 사이 값
	StudentNumber studentNumber `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~21"`         // 1~21 사이 값
	Name          name          `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"` // 2~4자 사이 한글
	PhoneNumber   phoneNumber   `gorm:"Type:char(11);NOT NULL" validate:"len=11,phone_number"`  // 11자
	PreProfileURI preProfileURI `gorm:"Type:varchar(150);NOT NULL"`
}

// 선생님 계정 테이블
type TeacherAuth struct {
	gorm.Model
	UUID      uuid      `gorm:"PRIMARY_KEY;Type:char(20);UNIQUE;INDEX" validate:"uuid=teacher,len=20"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	TeacherID teacherID `gorm:"varchar(20);NOT NULL" validate:"min=4,max=20,ascii"`                    // 4~20자 사이
	TeacherPW teacherPW `gorm:"varchar(100):NOT NULL"`
	Certified certified `gorm:"bool;default:false;NOT NULL"`
}

// 선생님 사용자 정보 테이블
type TeacherInform struct {
	gorm.Model
	TeacherUUID teacherUUID `gorm:"Type:char(20);UNIQUE;NOT NULL'" validate:"uuid=teacher,len=20"` // 형식 => 'teacher-' + 12자리 랜덤 수 (20자)
	Name        name        `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4"`               // 2~4자 사이 (원래 한글, PICK에서는 아니라서 지움)
	Grade       grade       `gorm:"Type:tinyint(1);" validate:"range=0~3"`                         // in (1~3)
	Class       class       `gorm:"Type:tinyint(1);" validate:"range=0~4"`                         // in (1~4)
	PhoneNumber phoneNumber `gorm:"Type:char(11)" validate:"phone_number"`                         // 휴대전화 형식
}

// 부모님 계정 테이블
type ParentAuth struct {
	gorm.Model
	UUID     uuid     `gorm:"PRIMARY_KEY;Type:char(19);UNIQUE;INDEX" validate:"uuid=parent,len=19"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	ParentID parentID `gorm:"Type:varchar(20);NOT NULL" validate:"min=4,max=20,ascii"`              // 4~20자 사이
	ParentPW parentPW `gorm:"Type:varchar(100);NOT NULL"`
}

// 부모님 사용자 정보 테이블
type ParentInform struct {
	gorm.Model
	ParentUUID  parentUUID  `gorm:"Type:char(19);UNIQUE;NOT NULL" validate:"uuid=parent,len=19"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	Name        name        `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"`      // 2~4자 사이 한글
	PhoneNumber phoneNumber `gorm:"Type:char(11)" validate:"phone_number"`                       // 11자
}

// 학부모 자녀 정보 테이블
type ParentChildren struct {
	gorm.Model
	ParentUUID    parentUUID    `gorm:"Type:char(19);NOT NULL" validate:"uuid=parent,len=19"` // 형식 => 'parent-' + 12자리 랜덤 수 (19자)
	Grade         grade         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~3"`               // 1~3 사이 값
	Class         class         `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~4"`               // 1~4 사이 값
	StudentNumber studentNumber `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~21"`              // 1~21 사이 값
	Name          name          `gorm:"Type:varchar(4);NOT NULL" validate:"min=2,max=4,korean"`      // 2~4자 사이 한글
	StudentUUID   studentUUID   `gorm:"Type:char(20)" validate:"uuid=student"`
}

// 관리자 계정 테이블
type AdminAuth struct {
	gorm.Model
	UUID    uuid    `gorm:"Type:char(18);UNIQUE;NOT NULL" validate:"uuid=admin,len=18"`
	AdminID adminID `gorm:"varchar(20);NOT NULL;UNIQUE" validate:"min=4,max=20,ascii"`
	AdminPW adminPW `gorm:"varchar(100):NOT NULL;"`
}
