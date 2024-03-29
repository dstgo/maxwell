package types

import (
	"github.com/dstgo/maxwell/app/conf"
	"gorm.io/gorm"
)

type Env struct {
	AppConf *conf.AppConf
	DB      *gorm.DB
}
