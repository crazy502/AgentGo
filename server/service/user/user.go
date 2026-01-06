// Package user service层提供用户领域的业务服务，实现用户登录、注册及验证码处理等核心业务逻辑
package user

import (
	"gorm.io/gorm/utils"

	"server/dao/user"
	"server/models"
)

// Login 用户登录
//
// Login 校验用户名和密码是否正确
// 登录成功后生成并返回 JWT Token
//
// 参数
//   - username：用户名
//   - password：明文密码
//
// 返回值
//   - string：登录成功时返回JWT Token，失败时为空字符串
//   - code.Code：业务状态码
func Login(username, password string) (string, code.Code) {
	var userInfo *models.User
	var ok bool
	//1:判断用户是否存在
	if ok, userInfo = user.IsExistUser(username); !ok {
		return "", code.CodeUserNotExist
	}
	//2:判断用户是否密码账号正确
	if userInfo.Password != utils.MD5(password) {
		return "", code.CodeInvalidPassword
	}
	//3:返回一个Token
	token, err := myjwt.GenerateToken(userInfo.ID, userInfo.Username)

	if err != nil {
		return "", code.CodeServerBusy
	}
	return token, code.CodeSuccess
}

// Register 用户注册
//
// Register 验证邮箱与验证码是否合法，
// 生成系统账号并写入数据库，
// 同时发送账号到用户邮箱，
// 最后返回登录Token。
//
// 参数
//   - email：用户邮箱
//   - password：明文密码
//   - captcha：邮箱验证码
//
// 返回值
//   - string：登录成功时返回JWT Token，失败时为空字符串
//   - code.Code：业务状态码
func Register(email, password, captcha string) (string, code.Code) {
	var ok bool
	var userInfo *models.User

	//1:先判断用户是否已经存在了
	if ok, _ := user.IsExistUser(email); ok {
		return "", code.CodeUserExist
	}

	//2:从redis中验证验证码是否有效
	if ok, _ := myredis.CheckCaptchaForEmail(email, captcha); ok {
		return "", code.CodeInvalidCaptcha
	}

	//3:生成11位账号
	username := utils.GetRandomNumbers(11)

	//4:注册到数据库中（调用DAO）
	if userInfo, ok = user.Register(username, email, password); !ok {
		return "", code.CodeServerBusy
	}

	//5:将账号一并发送到对应邮箱上，后续需要账号登录（业务通知）
	if err := myemail.SendCaptcha(email, username, user.UserNameMsg); err != nil {
		return "", code.CodeServerBusy
	}

	//6:生成Token（登录凭证）
	token, err := myjwt.GenerateToken(userInfo.ID, userInfo.Username)

	if err != nil {
		return "", code.CodeServerBusy
	}
	return token, code.CodeSuccess
}

// SendCaptcha 往指定邮箱发送验证码
//
// SendCaptcha 生成6位随机验证码，
// 先写入Redis（2分钟有效），
// 再通过邮件服务发送给用户。
//
// 参数
//   - email：用户邮箱
//
// 返回值
//   - code.Code：业务状态码
func SendCaptcha(email_ string) code.Code {
	send_code := utils.GetRandomNumbers(6)
	//1:先存放到redis中
	if err := myredis.SetCaptchaForEmail(email_, send_code); err != nil {
		return code.CodeServerBusy
	}
	//2:再进行远程发送
	if err := myemail.SendCaptcha(email_, send_code, myemail.CodeMsg); err != nil {
		return code.CodeServerBusy
	}
	return code.CodeSuccess
}
