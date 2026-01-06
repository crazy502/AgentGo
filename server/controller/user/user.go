// Package user controller层 实现用户相关的HTTP控制器
// 负责
//   - 解析和校验HTTP请求参数
//     -调用service层执行业务逻辑
//   - 构造并返回HTTP响应
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/service/user"
)

type (
	//目前username只能是账号登录（后期考虑添加邮箱账号均可）
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// omitempty当字段为空的时候，不返回这个东西
	LoginResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
	//验证码由后端生成，存放到redis中，固然需要先发送一次请求CaptchaRequest,然后用返回的验证码
	//邮箱以及密码进行注册，后续再将账号进行返回
	RegisterRequest struct {
		Email    string `json:"email" binding:"required"`
		Captcha  string `json:"captcha"`
		Password string `json:"password"`
	}
	//注册成功之后，直接让其进行登录状态
	RegisterResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}

	CaptchaRequest struct {
		Email string `json:"email" binding:"required"`
	}

	CaptchaResponse struct {
		controller.Response
	}
)

// Login 用户登录
//
// 处理流程：
//  1. 解析请求参数
//  2. 调用service层执行登录逻辑
//  3. 返回登录结果与JWT令牌
func Login(c *gin.Context) {

	req := new(LoginRequest)
	res := new(LoginResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Login(req.Username, req.Password)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)

}

// Register 用户注册
//
// 处理流程：
//  1. 解析请求参数
//  2. 调用service层执行注册逻辑
//  3. 构造并返回注册结果与JWT令牌
func Register(c *gin.Context) {

	req := new(RegisterRequest)
	res := new(RegisterResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Register(req.Email, req.Password, req.Captcha)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// HandleCaptcha 向指定邮箱发送验证码
//
// 处理流程：
//  1. 解析并校验邮箱参数
//  2. 生成验证码并存入Redis
//  3. 通过邮件系统发送验证码
func HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)
	//解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	//给service层进行处理
	code_ := user.SendCaptcha(req.Email)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}
	//匿名字段，其实本身res.Success()调用就是res.Response.Success()
	//res.Response.Success()
	res.Success()
	c.JSON(http.StatusOK, res)
}
