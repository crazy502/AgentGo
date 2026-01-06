// Package 负责用户的路由处理
package router

import (
	"github.com/gin-gonic/gin"

	"server/controller/user"
)

func RegisterUserRouter(r *gin.RouterGroup) {
	{
		r.POST("/register", user.Register)
		r.POST("/login", user.Login)
		r.POST("/captcha", user.HandleCaptcha)
	}

}
