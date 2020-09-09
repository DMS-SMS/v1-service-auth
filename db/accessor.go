package db

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type Accessor interface {
	// 계정 생성 메서드
	CreateStudentAuth(*model.StudentAuth) (auth *model.StudentAuth, err error)
	CreateTeacherAuth(*model.TeacherAuth) (auth *model.TeacherAuth, err error)
	CreateParentAuth(*model.ParentAuth) (auth *model.ParentAuth, err error)

	// 계정 ID로 계정 정보 조회 메서드
	GetStudentAuthWithID(studentID string) (*model.StudentAuth, error)
	GetTeacherAuthWithID(teacherID string) (*model.TeacherAuth, error)
	GetParentAuthWithID(parentID string) (*model.ParentAuth, error)

	// 비밀번호 변경 메서드
	ChangeStudentAuthPw(uuid string, studentPW string) error
	ChangeTeacherAuthPw(uuid string, teacherPW string) error
	ChangeParentAuthPw(uuid string, parentPW string) error

	// 계성 삭제 메서드 (Soft Delete)
	DeleteStudentAuth(uuid string) error
	DeleteTeacherAuth(uuid string) error
	DeleteParentAuth(uuid string) error

	// ---

	// 사용자 정보 추가 메서드
	CreateStudentInform(*model.StudentInform) (result *model.StudentInform, err error)
	CreateTeacherInform(*model.TeacherInform) (result *model.TeacherInform, err error)
	CreateParentInform(*model.ParentInform) (result *model.ParentInform, err error)

	// 사용자 정보로 uuid 조회 메서드 (계정 삭제 시 사용)
	GetStudentUUIDsWithInform(*model.StudentInform) (uuidArr []string, err error)
	GetTeacherUUIDsWithInform(*model.TeacherInform) (uuidArr []string, err error)
	GetParentUUIDsWithInform(*model.ParentInform) (uuidArr []string, err error)

	// 계정 UUID로 정보 조회 메서드
	GetStudentInformWithUUID(uuid string) (*model.StudentInform, error)
	GetTeacherInformWithUUID(uuid string) (*model.TeacherInform, error)
	GetParentInformWithUUID(uuid string) (*model.ParentInform, error)

	// 사용자 정보 수정 메서드
	ModifyStudentInform(uuid string, modify *model.StudentInform) (result *model.StudentInform, err error)
	ModifyTeacherInform(uuid string, modify *model.TeacherInform) (result *model.TeacherInform, err error)
	ModifyParentInform(uuid string, modify *model.ParentInform) (result *model.ParentInform, err error)

	// ---

	// 트랜잭션 관련 메서드
	Begin(db *gorm.DB)
	Commit() *gorm.DB
	Rollback() *gorm.DB
}
