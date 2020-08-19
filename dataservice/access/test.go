package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type Test struct {}

func (t Test) CreateStudentAuth(*model.StudentAuth) (result *model.StudentAuth, err error) { return }
func (t Test) CreateTeacherAuth(*model.TeacherAuth) (result *model.TeacherAuth, err error) { return }
func (t Test) CreateParentAuth(*model.ParentAuth) (result *model.ParentAuth, err error) { return }

// 사용자 정보 추가 메서드
func (t Test) CreateStudentInform(*model.StudentInform) (result *model.StudentInform, err error) { return }
func (t Test) CreateTeacherInform(*model.TeacherInform) (result *model.TeacherInform, err error) { return }
func (t Test) CreateParentInform(*model.ParentInform) (result *model.ParentInform, err error) { return }

// 정보 조회 메서드
func (t Test) GetStudentInform(sid string) (*model.StudentInform, error) { return nil, nil}
func (t Test) GetTeacherInform(tid string) (*model.TeacherInform, error) { return nil, nil }
func (t Test) GetParentInform(pid string) (*model.ParentInform, error) { return nil, nil }

// 비밀번호 변경 메서드
func (t Test) ChangeStudentAuthPw(sid string, pw string) error { return nil }
func (t Test) ChangeTeacherAuthPw(tid string, pw string) error { return nil }
func (t Test) ChangeParentAuthPw(pid string, pw string) error { return nil }

func (t Test) Begin(db *gorm.DB) {}
func (t Test) Commit() *gorm.DB { return nil }
func (t Test) Rollback() *gorm.DB { return nil }