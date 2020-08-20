package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type None struct {}
// 계정 생성 메서드
func (t None) CreateStudentAuth(*model.StudentAuth) (result *model.StudentAuth, err error) { return }
func (t None) CreateTeacherAuth(*model.TeacherAuth) (result *model.TeacherAuth, err error) { return }
func (t None) CreateParentAuth(*model.ParentAuth) (result *model.ParentAuth, err error) { return }

// 계정 ID로 계정 정보 조회 메서드
func (t None) GetStudentAuthWithID(id string) (*model.StudentAuth, error) { return nil, nil }
func (t None) GetTeacherAuthWithID(id string) (*model.TeacherAuth, error) { return nil, nil }
func (t None) GetParentAuthWithID(id string) (*model.ParentAuth, error) { return nil, nil }

// 비밀번호 변경 메서드
func (t None) ChangeStudentAuthPw(sid string, pw string) error { return nil }
func (t None) ChangeTeacherAuthPw(tid string, pw string) error { return nil }
func (t None) ChangeParentAuthPw(pid string, pw string) error { return nil }

// 계성 삭제 메서드 (Soft Delete)
func (t None) DeleteStudentAuth(sid uint) error { return nil }
func (t None) DeleteTeacherAuth(tid string) error { return nil }
func (t None) DeleteParentAuth(pid string) error { return nil }

// 사용자 정보 추가 메서드
func (t None) CreateStudentInform(*model.StudentInform) (result *model.StudentInform, err error) { return }
func (t None) CreateTeacherInform(*model.TeacherInform) (result *model.TeacherInform, err error) { return }
func (t None) CreateParentInform(*model.ParentInform) (result *model.ParentInform, err error) { return }

// 사용자 정보로 uuid 조회 메서드
func (t None) GetStudentUUIDWithInform(*model.StudentInform) (sid []string, err error) { return nil, nil }
func (t None) GetTeacherUUIDWithInform(*model.TeacherInform) (tid []string, err error) { return nil, nil }
func (t None) GetParentUUIDWithInform(*model.ParentInform) (pid []string, err error) { return nil, nil }

// 정보 조회 메서드 (계정 삭제 시 사용)
func (t None) GetStudentInformWithUUID(sid string) ([]*model.StudentInform, error) { return nil, nil }
func (t None) GetTeacherInformWithUUID(tid string) ([]*model.TeacherInform, error) { return nil, nil }
func (t None) GetParentInformWithUUID(pid string) ([]*model.ParentInform, error) { return nil, nil }

// 사용자 정보 수정 메서드
func (t None) ModifyStudentInform(sid string, modify *model.StudentInform) (result *model.StudentInform, err error) { return }
func (t None) ModifyTeacherInform(tid string, modify *model.TeacherInform) (result *model.TeacherInform, err error) { return }
func (t None) ModifyParentInform(pid string, modify *model.StudentInform) (result *model.ParentInform, err error) { return }

// 트랜잭션 관련 메서드
func (t None) Begin(db *gorm.DB) {}
func (t None) Commit() *gorm.DB { return nil }
func (t None) Rollback() *gorm.DB { return nil }