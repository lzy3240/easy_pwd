package global

import "gorm.io/gorm"

type PassWord struct {
	gorm.Model
	Title    string `gorm:"type:varchar(255);not null" json:"title"`
	Username string `gorm:"type:varchar(255);not null" json:"username"`
	Password string `gorm:"type:varchar(4096);not null" json:"password"`
	Category string `gorm:"type:varchar(255);" json:"category"`
	Url      string `gorm:"type:varchar(255);" json:"url"`
	Remark   string `gorm:"type:varchar(255);" json:"remark"`
}

func (PassWord) TableName() string {
	return "password"
}

type Setting struct {
	gorm.Model
	Key   string `gorm:"type:varchar(255);not null" json:"key"`
	Value string `gorm:"type:varchar(4096);not null" json:"value"`
}

func (Setting) TableName() string {
	return "setting"
}
