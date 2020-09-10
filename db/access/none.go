package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type None struct {}
// 계정 생성 메서드
func (t None) CreateStudentAuth(auth *model.StudentAuth) (resultAuth *model.StudentAuth, err error) { return }
func (t None) CreateTeacherAuth(auth *model.TeacherAuth) (resultAuth *model.TeacherAuth, err error) { return }
func (t None) CreateParentAuth(auth *model.ParentAuth) (resultAuth *model.ParentAuth, err error) { return }

// 계정 ID로 계정 정보 조회 메서드
func (t None) GetStudentAuthWithID(studentID string) (*model.StudentAuth, error) { return nil, nil }
func (t None) GetTeacherAuthWithID(teacherID string) (*model.TeacherAuth, error) { return nil, nil }
func (t None) GetParentAuthWithID(parentID string) (*model.ParentAuth, error) { return nil, nil }

// 비밀번호 변경 메서드
func (t None) ChangeStudentAuthPw(uuid string, studentPW string) error { return nil }
func (t None) ChangeTeacherAuthPw(uuid string, teacherPW string) error { return nil }
func (t None) ChangeParentAuthPw(uuid string, parentPW string) error { return nil }

// 계성 삭제 메서드 (Soft Delete)
func (t None) DeleteStudentAuth(uuid string) error { return nil }
func (t None) DeleteTeacherAuth(uuid string) error { return nil }
func (t None) DeleteParentAuth(uuid string) error { return nil }

// 사용자 정보 추가 메서드
func (t None) CreateStudentInform(inform *model.StudentInform) (resultInform *model.StudentInform, err error) { return }
func (t None) CreateTeacherInform(inform *model.TeacherInform) (resultInform *model.TeacherInform, err error) { return }
func (t None) CreateParentInform(inform *model.ParentInform) (resultInform *model.ParentInform, err error) { return }

// 사용자 정보로 uuid 조회 메서드
func (t None) GetStudentUUIDsWithInform(*model.StudentInform) (uuidArr []string, err error) { return }
func (t None) GetTeacherUUIDsWithInform(*model.TeacherInform) (uuidArr []string, err error) { return }
func (t None) GetParentUUIDsWithInform(*model.ParentInform) (uuidArr []string, err error) { return }

// 정보 조회 메서드 (계정 삭제 시 사용)
func (t None) GetStudentInformWithUUID(uuid string) (inform *model.StudentInform, err error) { return }
func (t None) GetTeacherInformWithUUID(uuid string) (inform *model.TeacherInform, err error) { return }
func (t None) GetParentInformWithUUID(uuid string) (inform *model.ParentInform, err error) { return }

// 사용자 정보 수정 메서드
func (t None) ModifyStudentInform(uuid string, revisionInform *model.StudentInform) (err error) { return }
func (t None) ModifyTeacherInform(uuid string, revisionInform *model.TeacherInform) (err error) { return }
func (t None) ModifyParentInform(uuid string, revisionInform *model.ParentInform) (err error) { return }

// 트랜잭션 관련 메서드
func (t None) Begin(db *gorm.DB) {}
func (t None) Commit() *gorm.DB { return nil }
func (t None) Rollback() *gorm.DB { return nil }