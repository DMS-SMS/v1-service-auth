package model

import (
	"reflect"
	"time"
)

// 매개변수로 전달 받은 변수의 DeepCopy 본사본을 생성하여 반환하는 함수
// 제약조건 -> 매개변수로 넘길 변수는 포인터 변수여야 함!! (X -> panic 발생)
func deepCopyModel(model interface{}) interface{} {
	duplicateModel := reflect.New(reflect.ValueOf(model).Elem().Type())
	duplicateModel.Elem().Set(reflect.ValueOf(model).Elem())
	return duplicateModel.Interface()
}

// deepCopyModel 함수를 이용하여 본사본 생성 후 gorm.Model 필드 값 초기화 후 해당 변수 반환 함수
// 제약조건 -> 매개변수로 넘길 변수는 구조체인 동시에 gorm.Model 객체의 필드들을 가지고 있어야 함!! (X -> panic 발생)
func exceptGormModel(model interface{}) (gormModelExceptTable interface{}) {
	gormModelExceptTable = deepCopyModel(model)

	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("ID").Set(reflect.ValueOf(uint(0)))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(time.Time{}))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Time{}))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("DeletedAt").Set(reflect.ValueOf((*time.Time)(nil)))
	return
}

// DeepCopy 메서드 -> 리시버 변수에 대한 DeepCopy 본사본 생성 및 반환 메서드
func (sa *StudentAuth)   DeepCopy() *StudentAuth   { return deepCopyModel(sa).(*StudentAuth) }
func (ta *TeacherAuth)   DeepCopy() *TeacherAuth   { return deepCopyModel(ta).(*TeacherAuth) }
func (pa *ParentAuth)    DeepCopy() *ParentAuth    { return deepCopyModel(pa).(*ParentAuth) }
func (si *StudentInform) DeepCopy() *StudentInform { return deepCopyModel(si).(*StudentInform) }
func (ti *TeacherInform) DeepCopy() *TeacherInform { return deepCopyModel(ti).(*TeacherInform) }
func (pi *ParentInform)  DeepCopy() *ParentInform  { return deepCopyModel(pi).(*ParentInform) }

// ExceptGormModel 메서드 -> 리시버 변수로부터 gorm.Model(임베딩 객체)에 포함되어있는 필드 값 초기화 후 반환 메서드
func (sa *StudentAuth)   ExceptGormModel() *StudentAuth   { return exceptGormModel(sa).(*StudentAuth) }
func (ta *TeacherAuth)   ExceptGormModel() *TeacherAuth   { return exceptGormModel(ta).(*TeacherAuth) }
func (pa *ParentAuth)    ExceptGormModel() *ParentAuth    { return exceptGormModel(pa).(*ParentAuth) }
func (si *StudentInform) ExceptGormModel() *StudentInform { return exceptGormModel(si).(*StudentInform) }
func (ti *TeacherInform) ExceptGormModel() *TeacherInform { return exceptGormModel(ti).(*TeacherInform) }
func (pi *ParentInform)  ExceptGormModel() *ParentInform  { return exceptGormModel(pi).(*ParentInform) }

// XXXConstraintName 메서드 -> XXX PK의 Constraint Name 값 반환 메서드
func (sa *StudentAuth)   ParentUUIDConstraintName()  string { return "student_auths_parent_uuid_parent_auths_uuid_foreign" }
func (si *StudentInform) StudentUUIDConstraintName() string { return "student_informs_student_uuid_student_auths_uuid_foreign" }
func (ti *TeacherInform) TeacherUUIDConstraintName() string { return "teacher_informs_teacher_uuid_teacher_auths_uuid_foreign" }
func (pi *ParentInform)  ParentUUIDConstraintName()  string { return "parent_informs_parent_uuid_parent_auths_uuid_foreign" }

// TableName 메서드 -> 리시버 변수에 해당되는 테이블의 이름 반환 메서드
func (sa *StudentAuth)   TableName() string { return "student_auths" }
func (ta *TeacherAuth)   TableName() string { return "teacher_auths" }
func (pa *ParentAuth)    TableName() string { return "parent_auths" }
func (si *StudentInform) TableName() string { return "student_informs" }
func (ti *TeacherInform) TableName() string { return "teacher_informs" }
func (pi *ParentInform)  TableName() string { return "parent_informs" }