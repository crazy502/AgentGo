package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey" json:"id"`                         //主键，自增
	Name      string         `gorm:"type:varchar(50)" json:"name"`                 //用户昵称
	Email     string         `gorm:"type:varchar(100);index" json:"email"`         //邮箱地址，支持索引，用于注册和验证码发送
	Username  string         `gorm:"type:varchar(50);uniqueIndex" json:"username"` //唯一索引。唯一用户名，用于登录（不同于邮箱）
	Password  string         `gorm:"type:varchar(255)" json:"-"`                   //MD5加密后的密码，不返回给前端
	CreatedAt time.Time      `json:"created_at"`                                   //自动时间戳
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"json:"-"` //支持软删除
}
