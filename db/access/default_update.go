package access

import (
	"auth/db/access/errors"
	"auth/model"
)

func (d *_default) ModifyStudentInform(uuid string, revisionInform *model.StudentInform) (err error) {
	contextForUpdate := make(map[string]interface{}, 6)

	if revisionInform.StudentUUID != emptyString {
		err = errors.StudentUUIDCannotBeChanged
		return
	}

	if revisionInform.Grade != emptyInt          { contextForUpdate[revisionInform.Grade.KeyName()] = revisionInform.Grade }
	if revisionInform.Class != emptyInt          { contextForUpdate[revisionInform.Class.KeyName()] = revisionInform.Class }
	if revisionInform.StudentNumber != emptyInt  { contextForUpdate[revisionInform.StudentNumber.KeyName()] = revisionInform.StudentNumber }
	if revisionInform.Name != emptyString        { contextForUpdate[revisionInform.Name.KeyName()] = revisionInform.Name }
	if revisionInform.PhoneNumber != emptyString { contextForUpdate[revisionInform.PhoneNumber.KeyName()] = revisionInform.PhoneNumber }
	if revisionInform.ProfileURI != emptyString  { contextForUpdate[revisionInform.ProfileURI.KeyName()] = revisionInform.ProfileURI }

	err = d.tx.Model(&model.StudentInform{}).Where("student_uuid = ?", uuid).Updates(contextForUpdate).Error
	return
}

func (d *_default) ModifyTeacherInform(uuid string, revisionInform *model.TeacherInform) (err error) {
	contextForUpdate := make(map[string]interface{}, 4)

	if revisionInform.TeacherUUID != emptyString {
		err = errors.TeacherUUIDCannotBeChanged
		return
	}

	if revisionInform.Name != emptyString        { contextForUpdate[revisionInform.Name.KeyName()] = revisionInform.Name }
	if revisionInform.PhoneNumber != emptyString { contextForUpdate[revisionInform.PhoneNumber.KeyName()] = revisionInform.PhoneNumber }

	if revisionInform.Grade != emptyInt {
		if int64(revisionInform.Grade) == model.TeacherInformInstance.Grade.NullReplaceValue() {
			contextForUpdate[revisionInform.Grade.KeyName()] = model.Grade(0)
		} else {
			contextForUpdate[revisionInform.Grade.KeyName()] = revisionInform.Grade
		}
	}

	if revisionInform.Class != emptyInt {
		if int64(revisionInform.Class) == model.TeacherInformInstance.Class.NullReplaceValue() {
			contextForUpdate[revisionInform.Class.KeyName()] = model.Class(0)
		} else {
			contextForUpdate[revisionInform.Class.KeyName()] = revisionInform.Class
		}
	}

	err = d.tx.Model(&model.TeacherInform{}).Where("teacher_uuid = ?", uuid).Updates(contextForUpdate).Error
	return
}

func (d *_default) ModifyParentInform(uuid string, revisionInform *model.ParentInform) (err error) {
	contextForUpdate := make(map[string]interface{}, 6)

	if revisionInform.ParentUUID != emptyString {
		err = errors.ParentUUIDCannotBeChanged
		return
	}

	if revisionInform.Name != emptyString        { contextForUpdate[revisionInform.Name.KeyName()] = revisionInform.Name }
	if revisionInform.PhoneNumber != emptyString { contextForUpdate[revisionInform.PhoneNumber.KeyName()] = revisionInform.PhoneNumber }

	err = d.tx.Model(&model.ParentInform{}).Where("parent_uuid = ?", uuid).Updates(contextForUpdate).Error
	return
}

func (d *_default) ChangeStudentPW(uuid string, studentPW string) (err error) {
	err = d.tx.Model(&model.StudentAuth{}).Where("uuid = ?", uuid).Update("student_pw", studentPW).Error
	return
}