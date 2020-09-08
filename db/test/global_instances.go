package test

import (
	"auth/db"
	"github.com/jinzhu/gorm"
)

var (
	manager db.AccessorManage
	dbc *gorm.DB
)

var passwords = map[string]string{
	"testPW1": "$2a$10$POwSnghOjkriuQ4w1Bj3zeHIGA7fXv8UI/UFXEhnnO5YrcwkUDcXq",
	"testPW2": "$2a$10$XxGXTboHZxhoqzKcBVqkJOiNSy6narAvIQ/ljfTJ4m93jAt8GyX.e",
	"testPW3": "$2a$10$sfZLOR8iVyhXI0y8nXcKIuKseahKu4NLSlocUWqoBdGrpLIZzxJ2S",
}