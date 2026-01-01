package user

import (
	"context"

	"gorm.io/gorm"

	"server/models"
)

const (
	CodeMsg     = "AgentGo验证码如下（验证码仅限于2分钟内有效）: "
	UserNameMsg = "AgentGo的账号如下：请保留好，后续可以用账号进行登录 "
)

var ctx = context.Background()

// 只能通过账号进行登录
func IsExistUser(username string) (bool, *models.User) {
	user, err := mysql.GetUserByUsername(username)

	if err == gorm.ErrRecordNotFound || user == nil {
		return false, nil
	}
	return true, user
}

func Register(username, email, password string) (*models.User, bool) {
	if user, err := mysql.InsertUser(&models.User{
		Email:    email,
		Name:     username,
		Username: username,
		Password: utils.MD5(password),
	}); err != nil {
		return nil, false
	} else {
		return user, true
	}
}
