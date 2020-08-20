package access

import (
	"auth/model"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock *mock.Mock
}

func NewMock(mock *mock.Mock) Mock {
	return Mock{
		mock: mock,
	}
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
	args := m.mock.Called(sid, pw)
	return args.Error(0)
}

func (m Mock) ChangeTeacherAuthPw(tid string, pw string) error {
	args := m.mock.Called(tid, pw)
	return args.Error(0)
}

func (m Mock) ChangeParentAuthPw(pid string, pw string) error {
	args := m.mock.Called(pid, pw)
	return args.Error(0)
}

// 계성 삭제 메서드 (Soft Delete)
func (m Mock) DeleteStudentAuth(sid uint) error {
	args := m.mock.Called(sid)
	return args.Error(0)
}

func (m Mock) DeleteTeacherAuth(tid string) error {
	args := m.mock.Called(tid)
	return args.Error(0)
}

func (m Mock) DeleteParentAuth(pid string) error {
	args := m.mock.Called(pid)
	return args.Error(0)
}
