package access

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

const (
	emptyString = ""
	emptyInt = 0
)

func (d *_default) GetStudentAuthWithID(studentID string) (auth *model.StudentAuth, err error) {
	auth = new(model.StudentAuth)
	err = d.tx.Where("student_id = ?", studentID).Find(&auth).Error
	return
}

func (d *_default) GetTeacherAuthWithID(teacherID string) (auth *model.TeacherAuth, err error) {
	auth = new(model.TeacherAuth)
	err = d.tx.Where("teacher_id = ?", teacherID).Find(&auth).Error
	return
}

func (d *_default) GetParentAuthWithID(parentID string) (auth *model.ParentAuth, err error) {
	auth = new(model.ParentAuth)
	err = d.tx.Where("parent_id = ?", parentID).Find(&auth).Error
	return
}

func (d *_default) GetAdminAuthWithID(adminID string) (auth *model.AdminAuth, err error) {
	auth = new(model.AdminAuth)
	err = d.tx.Where("admin_id = ?", adminID).Find(&auth).Error
	return
}

func (d *_default) GetStudentAuthWithUUID(uuid string) (auth *model.StudentAuth, err error) {
	auth = new(model.StudentAuth)
	err = d.tx.Where("uuid = ?", uuid).Find(auth).Error
	return
}

func (d *_default) GetTeacherAuthWithUUID(uuid string) (auth *model.TeacherAuth, err error) {
	auth = new(model.TeacherAuth)
	err = d.tx.Where("uuid = ?", uuid).Find(auth).Error
	return
}

func (d *_default) GetParentAuthWithUUID(uuid string) (auth *model.ParentAuth, err error) {
	auth = new(model.ParentAuth)
	err = d.tx.Where("uuid = ?", uuid).Find(auth).Error
	return
}

func (d *_default) GetStudentUUIDsWithInform(inform *model.StudentInform) (uuidArr []string, err error) {
	cascadeTx := d.tx.New()

	if inform.StudentUUID != emptyString { cascadeTx = cascadeTx.Where("student_uuid LIKE ?", "%"+inform.StudentUUID+"%") }
	if inform.Grade != emptyInt          { cascadeTx = cascadeTx.Where("grade = ?", inform.Grade) }
	if inform.Class != emptyInt          { cascadeTx = cascadeTx.Where("class = ?", inform.Class) }
	if inform.StudentNumber != emptyInt  { cascadeTx = cascadeTx.Where("student_number = ?", inform.StudentNumber) }
	if inform.Name != emptyString        { cascadeTx = cascadeTx.Where("name LIKE ?", "%"+inform.Name+"%") }
	if inform.PhoneNumber != emptyString { cascadeTx = cascadeTx.Where("phone_number LIKE ?", "%"+inform.PhoneNumber+"%") }
	if inform.ProfileURI != emptyString  { cascadeTx = cascadeTx.Where("profile_uri LIKE ?", "%"+inform.ProfileURI+"%") }

	informs := make([]*model.StudentInform, 1, 3)
	err = cascadeTx.Find(&informs).Error

	if len(informs) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	for _, inform := range informs {
		uuidArr = append(uuidArr, string(inform.StudentUUID))
	}
	return
}

func (d *_default) GetTeacherUUIDsWithInform(inform *model.TeacherInform) (uuidArr []string, err error) {
	cascadeTx := d.tx.New()

	if inform.TeacherUUID != emptyString { cascadeTx = cascadeTx.Where("teacher_uuid LIKE ?", "%"+inform.TeacherUUID+"%") }
	if inform.Grade != emptyInt          { cascadeTx = cascadeTx.Where("grade = ?", inform.Grade) }
	if inform.Class != emptyInt          { cascadeTx = cascadeTx.Where("class = ?", inform.Class) }
	if inform.Name != emptyString        { cascadeTx = cascadeTx.Where("name LIKE ?", "%"+inform.Name+"%") }
	if inform.PhoneNumber != emptyString { cascadeTx = cascadeTx.Where("phone_number LIKE ?", "%"+inform.PhoneNumber+"%") }

	if inform.Grade != emptyInt {
		if int64(inform.Grade) == model.TeacherInformInstance.Grade.NullReplaceValue() {
			cascadeTx = cascadeTx.Where("grade IS NULL")
		} else {
			cascadeTx = cascadeTx.Where("grade = ?", inform.Grade)
		}
	}

	if inform.Class != emptyInt {
		if int64(inform.Class) == model.TeacherInformInstance.Class.NullReplaceValue() {
			cascadeTx = cascadeTx.Where("class IS NULL")
		} else {
			cascadeTx = cascadeTx.Where("class = ?", inform.Class)
		}
	}


	informs := make([]*model.TeacherInform, 1, 3)
	err = cascadeTx.Find(&informs).Error

	if len(informs) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	for _, inform := range informs {
		uuidArr = append(uuidArr, string(inform.TeacherUUID))
	}
	return
}

func (d *_default) GetParentUUIDsWithInform(inform *model.ParentInform) (uuidArr []string, err error) {
	cascadeTx := d.tx.New()

	if inform.ParentUUID != emptyString  { cascadeTx = cascadeTx.Where("parent_uuid LIKE ?", inform.ParentUUID) }
	if inform.Name != emptyString        { cascadeTx = cascadeTx.Where("name LIKE ?", inform.Name) }
	if inform.PhoneNumber != emptyString { cascadeTx = cascadeTx.Where("phone_number LIKE ?", inform.PhoneNumber) }

	informs := make([]*model.ParentInform, 1, 3)
	err = cascadeTx.Find(&informs).Error

	if len(informs) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}

	for _, inform := range informs {
		uuidArr = append(uuidArr, string(inform.ParentUUID))
	}
	return
}

func (d *_default) GetStudentInformWithUUID(uuid string) (inform *model.StudentInform, err error) {
	inform = new(model.StudentInform)
	err = d.tx.Where("student_uuid = ?", uuid).Find(inform).Error
	return
}

func (d *_default) GetStudentInformsWithUUIDs(uuidArr []string) (informs []*model.StudentInform, err error) {
	for _, uuid := range uuidArr {
		inform := model.StudentInform{}
		selectErr := d.tx.Where("student_uuid = ?", uuid).Find(&inform).Error
		informs = append(informs, &inform)
		if selectErr != nil {
			err = selectErr
		}
	}
	return
}

func (d *_default) GetStudentInformsWithParentUUID(parentUUID string) (informs []*model.StudentInform, err error) {
	informs = make([]*model.StudentInform, 1, 3)
	err = d.tx.Raw("SELECT student_informs.* FROM student_auths, student_informs " +
		"WHERE student_auths.uuid = student_informs.student_uuid AND student_auths.parent_uuid = ?", parentUUID).Scan(&informs).Error
	return 
}

func (d *_default) GetTeacherInformWithUUID(uuid string) (inform *model.TeacherInform, err error) {
	inform = new(model.TeacherInform)
	err = d.tx.Where("teacher_uuid = ?", uuid).Find(inform).Error
	return
}

func (d *_default) GetParentInformWithUUID(uuid string) (inform *model.ParentInform, err error) {
	inform = new(model.ParentInform)
	err = d.tx.Where("parent_uuid = ?", uuid).Find(inform).Error
	return
}

func (d *_default) GetUnsignedStudents(targetGrade, targetGroup int64) (students []*model.UnsignedStudent, err error) {
	cascadeTx := d.tx.New()

	if targetGrade != emptyInt  { cascadeTx = cascadeTx.Where("grade = ?", model.Grade(targetGrade)) }
	if targetGroup != emptyInt  { cascadeTx = cascadeTx.Where("class = ?", model.Class(targetGroup)) }

	students = []*model.UnsignedStudent{}
	err = cascadeTx.Find(&students).Error
	
	return 
}

func (d *_default) GetUnsignedStudentWithAuthCode(authCode int64) (student *model.UnsignedStudent, err error) {
	student = new(model.UnsignedStudent)
	err = d.tx.Where("auth_code = ?", authCode).Find(student).Error
	return
}
