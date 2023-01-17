package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisSwagger(engine *gin.Engine) {

	engine.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
