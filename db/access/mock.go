package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type _mock struct {
	mock *mock.Mock
}

func Mock(mock *mock.Mock) _mock {
	return _mock{mock: mock}
}

// 계정 생성 메서드
func (m _mock) CreateStudentAuth(auth *model.StudentAuth) (resultAuth *model.StudentAuth, err error) {
	auth.StudentPW = ""
	args := m.mock.Called(auth)
	return args.Get(0).(*model.StudentAuth), args.Error(1)
}

func (m _mock) CreateTeacherAuth(auth *model.TeacherAuth) (resultAuth *model.TeacherAuth, err error) {
	auth.TeacherPW = ""
	args := m.mock.Called(auth)
	return args.Get(0).(*model.TeacherAuth), args.Error(1)
}

func (m _mock) CreateParentAuth(auth *model.ParentAuth) (resultAuth *model.ParentAuth, err error) {
	auth.ParentPW = ""
	args := m.mock.Called(auth)
	return args.Get(0).(*model.ParentAuth), args.Error(1)
}

func (m _mock) CreateParentChildren(auth *model.ParentChildren) (*model.ParentChildren, error) {
	args := m.mock.Called(auth)
	return args.Get(0).(*model.ParentChildren), args.Error(1)
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

func (m _mock) GetAdminAuthWithID(adminID string) (*model.AdminAuth, error) {
	args := m.mock.Called(adminID)
	return args.Get(0).(*model.AdminAuth), args.Error(1)
}

// UUID로 계정 존재 여부 확인 메서드
func (m _mock) GetStudentAuthWithUUID(uuid string) (*model.StudentAuth, error) {
	args := m.mock.Called(uuid)
	return args.Get(0).(*model.StudentAuth), args.Error(1)
}

func (m _mock) GetTeacherAuthWithUUID(uuid string) (*model.TeacherAuth, error) {
	args := m.mock.Called(uuid)
	return args.Get(0).(*model.TeacherAuth), args.Error(1)
}

func (m _mock) GetParentAuthWithUUID(uuid string) (*model.ParentAuth, error) {
	args := m.mock.Called(uuid)
	return args.Get(0).(*model.ParentAuth), args.Error(1)
}

// 비밀번호 변경 메서드
func (m _mock) ChangeStudentPW(uuid string, studentPW string) error {
	studentPW = ""
	return m.mock.Called(uuid, studentPW).Error(0)
}

func (m _mock) ChangeTeacherPW(uuid string, teacherPW string) error {
	teacherPW = ""
	return m.mock.Called(uuid, teacherPW).Error(0)
}

func (m _mock) ChangeParentPW(uuid string, parentPW string) error {
	parentPW = ""
	return m.mock.Called(uuid, parentPW).Error(0)
}

func (m _mock) ChangeParentUUID(studentUUID string, parentUUID string) error {
	return m.mock.Called(studentUUID, parentUUID).Error(0)
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
func (m _mock) CreateStudentInform(inform *model.StudentInform) (resultInform *model.StudentInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.StudentInform), args.Error(1)
}

func (m _mock) CreateTeacherInform(inform *model.TeacherInform) (resultInform *model.TeacherInform, err error) {
	args := m.mock.Called(inform)
	return args.Get(0).(*model.TeacherInform), args.Error(1)
}

func (m _mock) CreateParentInform(inform *model.ParentInform) (resultInform *model.ParentInform, err error) {
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

func (m _mock) GetStudentInformsWithUUIDs(uuidArr []string) ([]*model.StudentInform, error) {
	args := m.mock.Called(uuidArr)
	return args.Get(0).([]*model.StudentInform), args.Error(1)
}

func (m _mock) GetStudentInformsWithParentUUID(parentUUID string) ([]*model.StudentInform, error) {
	args := m.mock.Called(parentUUID)
	return args.Get(0).([]*model.StudentInform), args.Error(1)
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
func (m _mock) ModifyStudentInform(uuid string, revisionInform *model.StudentInform) (err error) {
	args := m.mock.Called(uuid, revisionInform)
	return args.Error(0)
}

func (m _mock) ModifyTeacherInform(uuid string, revisionInform *model.TeacherInform) (err error) {
	args := m.mock.Called(uuid, revisionInform)
	return args.Error(0)
}

func (m _mock) ModifyParentInform(uuid string, revisionInform *model.ParentInform) (err error) {
	args := m.mock.Called(uuid, revisionInform)
	return args.Error(0)
}

// 사용자 정보 삭제 메서드 (Soft Delete)
func (m _mock) DeleteStudentInform(studentUUID string) error {
	return m.mock.Called(studentUUID).Error(0)
}

func (m _mock) DeleteTeacherInform(teacherUUID string) error {
	return m.mock.Called(teacherUUID).Error(0)
}

func (m _mock) DeleteParentInform(parentUUID string) error {
	return m.mock.Called(parentUUID).Error(0)
}

// ---

// 예비 계정 생성 메서드
func (m _mock) AddUnsignedStudent(student *model.UnsignedStudent) (*model.UnsignedStudent, error) {
	args := m.mock.Called(student)
	return args.Get(0).(*model.UnsignedStudent), args.Error(1)
}

func (m _mock) 	GetUnsignedStudents(targetGrade, targetGroup, targetNumber int64) ([]*model.UnsignedStudent, error) {
	args := m.mock.Called(targetGrade, targetGroup)
	return args.Get(0).([]*model.UnsignedStudent), args.Error(1)
}

func (m _mock) GetUnsignedStudentWithAuthCode(authCode int64) (*model.UnsignedStudent, error) {
	args := m.mock.Called(authCode)
	return args.Get(0).(*model.UnsignedStudent), args.Error(1)
}

func (m _mock) GetParentChildWithInform(grade, group, number int64, name string) (*model.ParentChildren, error) {
	args := m.mock.Called(grade, group, number, name)
	return args.Get(0).(*model.ParentChildren), args.Error(1)
}

func (m _mock) ModifyParentChildren(current *model.ParentChildren, revision *model.ParentChildren) error {
	return m.mock.Called(current, revision).Error(0)
}

func (m _mock) DeleteUnsignedStudent(authCode int64) error {
	return m.mock.Called(authCode).Error(0)
}

// ---

// 트랜잭션 관련 메서드
func (m _mock) BeginTx() {
	m.mock.Called()
}

func (m _mock) Commit() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}

func (m _mock) Rollback() *gorm.DB {
	return m.mock.Called().Get(0).(*gorm.DB)
}
