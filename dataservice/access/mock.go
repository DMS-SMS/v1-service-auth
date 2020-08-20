package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock *mock.Mock
}

func NewMock(mock *mock.Mock) Mock {
	return Mock{mock: mock}
}

// 계정 생성 메서드
func (m Mock) CreateStudentAuth(auth *model.StudentAuth) (result *model.StudentAuth, err error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.StudentAuth), args.Error(1)
}

func (m Mock) CreateTeacherAuth(auth *model.TeacherAuth) (result *model.TeacherAuth, err error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.TeacherAuth), args.Error(1)
}

func (m Mock) CreateParentAuth(auth *model.ParentAuth) (result *model.ParentAuth, err error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.ParentAuth), args.Error(1)
}

// 계정 ID로 계정 정보 조회 메서드
func (m Mock) GetStudentAuthWithID(id string) (*model.StudentAuth, error) {
	args := m.mock.Called(id)
	return args.Get(0).(*model.StudentAuth), args.Error(1)
}

func (m Mock) GetTeacherAuthWithID(id string) (*model.TeacherAuth, error) {
	args := m.mock.Called(id)
	return args.Get(0).(*model.TeacherAuth), args.Error(1)
}

func (m Mock) GetParentAuthWithID(id string) (*model.ParentAuth, error) {
	args := m.mock.Called(id)
	return args.Get(0).(*model.ParentAuth), args.Error(1)
}

// 비밀번호 변경 메서드
func (m Mock) ChangeStudentAuthPw(sid string, pw string) error {
	return m.mock.Called(sid, pw).Error(0)
}

func (m Mock) ChangeTeacherAuthPw(tid string, pw string) error {
	return m.mock.Called(tid, pw).Error(0)
}

func (m Mock) ChangeParentAuthPw(pid string, pw string) error {
	return m.mock.Called(pid, pw).Error(0)
}

// 계성 삭제 메서드 (Soft Delete)
func (m Mock) DeleteStudentAuth(sid uint) error {
	return m.mock.Called(sid).Error(0)
}

func (m Mock) DeleteTeacherAuth(tid string) error {
	return m.mock.Called(tid).Error(0)
}

func (m Mock) DeleteParentAuth(pid string) error {
	return m.mock.Called(pid).Error(0)
}

// ---

// 사용자 정보 추가 메서드
func (m Mock) CreateStudentInform(inform *model.StudentInform) (result *model.StudentInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m Mock) CreateTeacherInform(inform *model.TeacherInform) (result *model.TeacherInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m Mock) CreateParentInform(inform *model.ParentInform) (result *model.ParentInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.ParentInform), args.Error(1)
}

// 사용자 정보로 uuid 조회 메서드 (계정 삭제 시 사용)
func (m Mock) GetStudentUUIDWithInform(inform *model.StudentInform) (sid []string, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).([]string), args.Error(1)
}

func (m Mock) GetTeacherUUIDWithInform(inform *model.TeacherInform) (tid []string, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).([]string), args.Error(1)
}

func (m Mock) GetParentUUIDWithInform(inform *model.ParentInform) (pid []string, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).([]string), args.Error(1)
}

// 정보 조회 메서드
func (m Mock) GetStudentInformWithUUID(sid string) (*model.StudentInform, error) {
	args := m.mock.Called(sid)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m Mock) GetTeacherInformWithUUID(tid string) (*model.TeacherInform, error) {
	args := m.mock.Called(tid)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m Mock) GetParentInformWithUUID(pid string) (*model.ParentInform, error) {
	args := m.mock.Called(pid)
	return args.Get(0).(*model.ParentInform), args.Error(1)
}

// 사용자 정보 수정 메서드
func (m Mock) ModifyStudentInform(sid string, modify *model.StudentInform) (result *model.StudentInform, err error) {
	args := m.mock.Called(sid, modify)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m Mock) ModifyTeacherInform(tid string, modify *model.TeacherInform) (result *model.TeacherInform, err error) {
	args := m.mock.Called(tid, modify)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m Mock) ModifyParentInform(pid string, modify *model.ParentInform) (result *model.ParentInform, err error) {
	args := m.mock.Called(pid, modify)
	return args.Get(0).(*model.ParentInform), args.Error(1)
}

// ---

// 트랜잭션 관련 메서드
func (m Mock) Begin(db *gorm.DB) {
	m.mock.Called(db)
}

func (m Mock) Commit() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}

func (m Mock) Rollback() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}