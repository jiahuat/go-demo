package main

import (
	"log"

	_ "westone.com/wscp-restful/api"
)

// 添加注释以描述 server 信息
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// parse flag
	// parse config
	// new server
	// start servershou11
	log.Println("start server")
	execute()
	log.Println("end of server")
}
