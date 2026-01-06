package server

import (
	"fmt"

	"server/router"
)

func StartServer(addr string, port string) error {
	r := router.InitRouter()
	return r.Run(fmt.Sprintf("%s:%s", addr, port))
}

//TODO 加载配置
