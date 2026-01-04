// Package user 提供用户相关的数据访问功能
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

// IsExistUser 判断用户是否存在（系统仅支持通过账号进行登录）
//
// 参数:
//   - username: 用户名
//
// 返回值:
//   - bool: 用户存在返回true，否则返回false
//   - *models.User: 用户存在时返回用户信息，否则返回nil
func IsExistUser(username string) (bool, *models.User) {
	user, err := mysql.GetUserByUsername(username)

	if err == gorm.ErrRecordNotFound || user == nil {
		return false, nil
	}
	return true, user
}

// Register 用户注册
//
// 参数:
//   - username: 用户名
//   - email: 邮箱地址
//   - password: 密码
//
// 返回值:
//   - *models.User: 注册成功时返回用户信息，失败时返回nil
//   - bool: 注册成功返回true，失败返回false
func Register(username, email, password string) (*models.User, bool) {
	user, err := mysql.InsertUser(&models.User{
		Email:    email,
		Name:     username,
		Username: username,
		Password: utils.MD5(password),
	})
	if err != nil {
		return nil, false
	}
	return user, true
}
