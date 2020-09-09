package db

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type Accessor interface {
	// 계정 생성 메서드
	CreateStudentAuth(*model.StudentAuth) (result *model.StudentAuth, err error)
	CreateTeacherAuth(*model.TeacherAuth) (result *model.TeacherAuth, err error)
	CreateParentAuth(*model.ParentAuth) (result *model.ParentAuth, err error)

	// 계정 ID로 계정 정보 조회 메서드
	GetStudentAuthWithID(id string) (*model.StudentAuth, error)
	GetTeacherAuthWithID(id string) (*model.TeacherAuth, error)
	GetParentAuthWithID(id string) (*model.ParentAuth, error)

	// 비밀번호 변경 메서드
	ChangeStudentAuthPw(sid string, pw string) error
	ChangeTeacherAuthPw(tid string, pw string) error
	ChangeParentAuthPw(pid string, pw string) error

	// 계성 삭제 메서드 (Soft Delete)
	DeleteStudentAuth(sid uint) error
	DeleteTeacherAuth(tid string) error
	DeleteParentAuth(pid string) error

	// ---

	// 사용자 정보 추가 메서드
	CreateStudentInform(*model.StudentInform) (result *model.StudentInform, err error)
	CreateTeacherInform(*model.TeacherInform) (result *model.TeacherInform, err error)
	CreateParentInform(*model.ParentInform) (result *model.ParentInform, err error)

	// 사용자 정보로 uuid 조회 메서드 (계정 삭제 시 사용)
	GetStudentUUIDWithInform(*model.StudentInform) (sid []string, err error)
	GetTeacherUUIDWithInform(*model.TeacherInform) (tid []string, err error)
	GetParentUUIDWithInform(*model.ParentInform) (pid []string, err error)

	// 계정 UUID로 정보 조회 메서드
	GetStudentInformWithUUID(sid string) (*model.StudentInform, error)
	GetTeacherInformWithUUID(tid string) (*model.TeacherInform, error)
	GetParentInformWithUUID(pid string) (*model.ParentInform, error)

	// 사용자 정보 수정 메서드
	ModifyStudentInform(sid string, modify *model.StudentInform) (result *model.StudentInform, err error)
	ModifyTeacherInform(tid string, modify *model.TeacherInform) (result *model.TeacherInform, err error)
	ModifyParentInform(pid string, modify *model.ParentInform) (result *model.ParentInform, err error)

	// ---

	// 트랜잭션 관련 메서드
	Begin(db *gorm.DB)
	Commit() *gorm.DB
	Rollback() *gorm.DB
}
