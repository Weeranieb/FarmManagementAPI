package mock

import "gorm.io/gorm"

type ContextDB struct {
	Postgresql *gorm.DB
}

var Context ContextDB
