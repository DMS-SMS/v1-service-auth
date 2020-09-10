package access

import (
	"auth/db/access/errors"
	"auth/model"
)

func (d *_default) ModifyStudentInform(uuid string, revisionInform *model.StudentInform) (err error) {
	contextForUpdate := make(map[string]interface{})

	if revisionInform.StudentUUID != emptyString {
		err = errors.StudentUUIDCannotBeChanged
		return
	}

	if revisionInform.Grade != emptyInt {
		contextForUpdate[model.StudentInformInstance.Grade.KeyName()] = revisionInform.Grade
	}

	if revisionInform.Class != emptyInt {
		contextForUpdate[model.StudentInformInstance.Class.KeyName()] = revisionInform.Class
	}

	if revisionInform.StudentNumber != emptyInt {
		contextForUpdate[model.StudentInformInstance.StudentNumber.KeyName()] = revisionInform.StudentNumber
	}

	if revisionInform.Name != emptyString {
		contextForUpdate[model.StudentInformInstance.Name.KeyName()] = revisionInform.Name
	}

	if revisionInform.PhoneNumber != emptyString {
		contextForUpdate[model.StudentInformInstance.PhoneNumber.KeyName()] = revisionInform.PhoneNumber
	}

	if revisionInform.ProfileURI != emptyString {
		contextForUpdate[model.StudentInformInstance.ProfileURI.KeyName()] = revisionInform.ProfileURI
	}

	err = d.tx.Model(&model.StudentInform{}).Where("student_uuid = ?", uuid).Updates(contextForUpdate).Error
	return
}