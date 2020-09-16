package db

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type Accessor interface {
	// 계정 생성 메서드
	CreateStudentAuth(auth *model.StudentAuth) (resultAuth *model.StudentAuth, err error)
	CreateTeacherAuth(auth *model.TeacherAuth) (resultAuth *model.TeacherAuth, err error)
	CreateParentAuth(auth *model.ParentAuth) (resultAuth *model.ParentAuth, err error)

	// 계정 ID로 계정 정보 조회 메서드
	GetStudentAuthWithID(studentID string) (*model.StudentAuth, error)
	GetTeacherAuthWithID(teacherID string) (*model.TeacherAuth, error)
	GetParentAuthWithID(parentID string) (*model.ParentAuth, error)

	// 계정 UUID로 계정 정보 조회 메서드
	GetStudentAuthWithUUID(uuid string) (*model.StudentAuth, error)
	GetTeacherAuthWithUUID(uuid string) (*model.TeacherAuth, error)
	GetParentAuthWithUUID(uuid string) (*model.ParentAuth, error)

	// 비밀번호 변경 메서드
	ChangeStudentPW(uuid string, studentPW string) error
	ChangeTeacherPW(uuid string, teacherPW string) error
	ChangeParentPW(uuid string, parentPW string) error

	// 계성 삭제 메서드 (Soft Delete)
	DeleteStudentAuth(uuid string) error
	DeleteTeacherAuth(uuid string) error
	DeleteParentAuth(uuid string) error

	// ---

	// 사용자 정보 추가 메서드
	CreateStudentInform(inform *model.StudentInform) (resultInform *model.StudentInform, err error)
	CreateTeacherInform(inform *model.TeacherInform) (resultInform *model.TeacherInform, err error)
	CreateParentInform(inform *model.ParentInform) (resultInform *model.ParentInform, err error)

	// 사용자 정보로 uuid 조회 메서드 (계정 삭제 시 사용)
	GetStudentUUIDsWithInform(*model.StudentInform) (uuidArr []string, err error)
	GetTeacherUUIDsWithInform(*model.TeacherInform) (uuidArr []string, err error)
	GetParentUUIDsWithInform(*model.ParentInform) (uuidArr []string, err error)

	// 계정 UUID로 정보 조회 메서드
	GetStudentInformWithUUID(uuid string) (*model.StudentInform, error)
	GetTeacherInformWithUUID(uuid string) (*model.TeacherInform, error)
	GetParentInformWithUUID(uuid string) (*model.ParentInform, error)

	// 사용자 정보 수정 메서드
	ModifyStudentInform(uuid string, revisionInform *model.StudentInform) (err error)
	ModifyTeacherInform(uuid string, revisionInform *model.TeacherInform) (err error)
	ModifyParentInform(uuid string, revisionInform *model.ParentInform) (err error)

	// ---

	// 트랜잭션 관련 메서드
	BeginTx()
	Commit() *gorm.DB
	Rollback() *gorm.DB
}
