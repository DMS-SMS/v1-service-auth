package dataservice

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"reflect"
)

type AuthDataAccessor interface {
	// 계정 생성 메서드
	CreateStudentAuth(*model.StudentAuth) (result *model.StudentAuth, err error)
	CreateTeacherAuth(*model.TeacherAuth) (result *model.TeacherAuth, err error)
	CreateParentAuth(*model.ParentAuth) (result *model.ParentAuth, err error)

	// 사용자 정보 추가 메서드
	CreateStudentInform(*model.StudentInform) (result *model.StudentInform, err error)
	CreateTeacherInform(*model.TeacherInform) (result *model.TeacherInform, err error)
	CreateParentInform(*model.ParentInform) (result *model.ParentInform, err error)

	// 정보 조회 메서드
	GetStudentInform(sid string) (*model.StudentInform, error)
	GetTeacherInform(tid string) (*model.TeacherInform, error)
	GetParentInform(pid string) (*model.ParentInform, error)

	// 비밀번호 변경 메서드
	ChangeStudentAuthPw(sid string, pw string) error
	ChangeTeacherAuthPw(tid string, pw string) error
	ChangeParentAuthPw(pid string, pw string) error

	Begin(db *gorm.DB)
	Commit() *gorm.DB
	Rollback() *gorm.DB
}

type AuthDataTxManage struct {
	db *gorm.DB
	accessor AuthDataAccessor
}

func NewAuthDataTxManage(db *gorm.DB, accessor AuthDataAccessor) AuthDataTxManage {
	return AuthDataTxManage{
		db:       db,
		accessor: accessor,
	}
}

func (atm AuthDataTxManage) BeginTx() (accessor AuthDataAccessor) {
	t := reflect.TypeOf(atm.accessor).Elem()
	accessor = reflect.New(t).Elem().Interface().(AuthDataAccessor)
	//fmt.Println(reflect.TypeOf(accessor))
	accessor.Begin(atm.db)
	return
}