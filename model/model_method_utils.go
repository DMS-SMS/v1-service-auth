package model

import (
	"reflect"
	"time"
)

// ExceptGormModel 메서드 -> 리시버 변수로부터 gorm.Model(임베딩 객체)에 포함되어있는 필드 값 초기화 후 반환 메서드
func (sa *StudentAuth) ExceptGormModel() *StudentAuth { return exceptGormModel(sa).(*StudentAuth) }

func (ta *TeacherAuth) ExceptGormModel() *TeacherAuth { return exceptGormModel(ta).(*TeacherAuth) }

func (pa *ParentAuth) ExceptGormModel() *ParentAuth  { return exceptGormModel(pa).(*ParentAuth) }

func (si *StudentInform) ExceptGormModel() *StudentInform { return exceptGormModel(si).(*StudentInform) }

func (ti *TeacherInform) ExceptGormModel() *TeacherInform { return exceptGormModel(ti).(*TeacherInform) }

func (pi *ParentInform) ExceptGormModel() *ParentInform { return exceptGormModel(pi).(*ParentInform) }

// 매개변수로 전달한 변수로부터 해당 타입의 새로운 변수를 선언하여 gorm.Model 필드 데이터 초기화 후 해당 변수 반환 함수
// 제약조건 -> 매개변수로 넘길 변수는 포인터 변수여야 하고 구조체인 동시에 gorm.Model 객체의 필드들을 가지고 있어야 함!! (X -> panic 발생)
func exceptGormModel(table interface{}) (gormModelExceptTable interface{}) {
	gormModelExceptTable = reflect.New(reflect.ValueOf(table).Elem().Type()).Interface()
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("ID").Set(reflect.ValueOf(uint(0)))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(time.Time{}))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Time{}))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("DeletedAt").Set(reflect.ValueOf((*time.Time)(nil)))
	return
}