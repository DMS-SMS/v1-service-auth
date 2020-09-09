package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type _mock struct {
	mock *mock.Mock
}

func NewMock(mock *mock.Mock) _mock {
	return _mock{mock: mock}
}

// 계정 생성 메서드
func (m _mock) CreateStudentAuth(auth *model.StudentAuth) (result *model.StudentAuth, err error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.StudentAuth), args.Error(1)
}

func (m _mock) CreateTeacherAuth(auth *model.TeacherAuth) (result *model.TeacherAuth, err error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.TeacherAuth), args.Error(1)
}

func (m _mock) CreateParentAuth(auth *model.ParentAuth) (result *model.ParentAuth, err error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.ParentAuth), args.Error(1)
}

// 계정 ID로 계정 정보 조회 메서드
func (m _mock) GetStudentAuthWithID(studentID string) (*model.StudentAuth, error) {
	args := m.mock.Called(studentID)
	return args.Get(0).(*model.StudentAuth), args.Error(1)
}

func (m _mock) GetTeacherAuthWithID(teacherID string) (*model.TeacherAuth, error) {
	args := m.mock.Called(teacherID)
	return args.Get(0).(*model.TeacherAuth), args.Error(1)
}

func (m _mock) GetParentAuthWithID(parentID string) (*model.ParentAuth, error) {
	args := m.mock.Called(parentID)
	return args.Get(0).(*model.ParentAuth), args.Error(1)
}

// 비밀번호 변경 메서드
func (m _mock) ChangeStudentAuthPw(uuid string, studentPW string) error {
	return m.mock.Called(uuid, studentPW).Error(0)
}

func (m _mock) ChangeTeacherAuthPw(uuid string, teacherPW string) error {
	return m.mock.Called(uuid, teacherPW).Error(0)
}

func (m _mock) ChangeParentAuthPw(uuid string, parentPW string) error {
	return m.mock.Called(uuid, parentPW).Error(0)
}

// 계성 삭제 메서드 (Soft Delete)
func (m _mock) DeleteStudentAuth(uuid string) error {
	return m.mock.Called(uuid).Error(0)
}

func (m _mock) DeleteTeacherAuth(uuid string) error {
	return m.mock.Called(uuid).Error(0)
}

func (m _mock) DeleteParentAuth(uuid string) error {
	return m.mock.Called(uuid).Error(0)
}

// ---

// 사용자 정보 추가 메서드
func (m _mock) CreateStudentInform(inform *model.StudentInform) (result *model.StudentInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m _mock) CreateTeacherInform(inform *model.TeacherInform) (result *model.TeacherInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m _mock) CreateParentInform(inform *model.ParentInform) (result *model.ParentInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.ParentInform), args.Error(1)
}

// 사용자 정보로 uuid 조회 메서드 (계정 삭제 시 사용)
func (m _mock) GetStudentUUIDsWithInform(inform *model.StudentInform) (uuidArr []string, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).([]string), args.Error(1)
}

func (m _mock) GetTeacherUUIDsWithInform(inform *model.TeacherInform) (uuidArr []string, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).([]string), args.Error(1)
}

func (m _mock) GetParentUUIDsWithInform(inform *model.ParentInform) (uuidArr []string, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).([]string), args.Error(1)
}

// 정보 조회 메서드
func (m _mock) GetStudentInformWithUUID(uuid string) (*model.StudentInform, error) {
	args := m.mock.Called(uuid)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m _mock) GetTeacherInformWithUUID(uuid string) (*model.TeacherInform, error) {
	args := m.mock.Called(uuid)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m _mock) GetParentInformWithUUID(uuid string) (*model.ParentInform, error) {
	args := m.mock.Called(uuid)
	return args.Get(0).(*model.ParentInform), args.Error(1)
}

// 사용자 정보 수정 메서드
func (m _mock) ModifyStudentInform(uuid string, modify *model.StudentInform) (result *model.StudentInform, err error) {
	args := m.mock.Called(uuid, modify)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m _mock) ModifyTeacherInform(uuid string, modify *model.TeacherInform) (result *model.TeacherInform, err error) {
	args := m.mock.Called(uuid, modify)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m _mock) ModifyParentInform(uuid string, modify *model.ParentInform) (result *model.ParentInform, err error) {
	args := m.mock.Called(uuid, modify)
	return args.Get(0).(*model.ParentInform), args.Error(1)
}

// ---

// 트랜잭션 관련 메서드
func (m _mock) Begin(db *gorm.DB) {
	m.mock.Called(db)
}

func (m _mock) Commit() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}

func (m _mock) Rollback() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}