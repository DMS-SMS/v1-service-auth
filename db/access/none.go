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
func (t None) GetStudentAuthWithID(studentID string) (auth *model.StudentAuth, err error) { return }
func (t None) GetTeacherAuthWithID(teacherID string) (auth *model.TeacherAuth, err error) { return }
func (t None) GetParentAuthWithID(parentID string) (auth *model.ParentAuth, err error) { return }
func (t None) GetAdminAuthWithID(adminID string) (auth *model.AdminAuth, err error) { return }

// UUID로 계정 존재 여부 확인 메서드
func (t None) GetStudentAuthWithUUID(uuid string) (auth *model.StudentAuth, err error) { return }
func (t None) GetTeacherAuthWithUUID(uuid string) (auth *model.TeacherAuth, err error) { return }
func (t None) GetParentAuthWithUUID(uuid string) (auth *model.ParentAuth, err error) { return }

// 비밀번호 변경 메서드
func (t None) ChangeStudentPW(uuid string, studentPW string) error { return nil }
func (t None) ChangeTeacherPW(uuid string, teacherPW string) error { return nil }
func (t None) ChangeParentPW(uuid string, parentPW string) error { return nil }

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
func (t None) GetStudentInformsWithUUIDs(uuidArr []string) (informs []*model.StudentInform, err error) { return }
func (t None) GetTeacherInformWithUUID(uuid string) (inform *model.TeacherInform, err error) { return }
func (t None) GetParentInformWithUUID(uuid string) (inform *model.ParentInform, err error) { return }

// 사용자 정보 수정 메서드
func (t None) ModifyStudentInform(uuid string, revisionInform *model.StudentInform) (err error) { return }
func (t None) ModifyTeacherInform(uuid string, revisionInform *model.TeacherInform) (err error) { return }
func (t None) ModifyParentInform(uuid string, revisionInform *model.ParentInform) (err error) { return }

// 트랜잭션 관련 메서드
func (t None) BeginTx() {}
func (t None) Commit() *gorm.DB { return nil }
func (t None) Rollback() *gorm.DB { return nil }