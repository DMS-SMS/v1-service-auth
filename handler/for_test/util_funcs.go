package test

import (
	"auth/tool/random"
	"github.com/jinzhu/gorm"
	"time"
)

func createGormModelOnCurrentTime() gorm.Model {
	currentTime := time.Now()
	return gorm.Model{
		ID:        uint(random.Int64WithLength(3)),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		DeletedAt: nil,
	}
}