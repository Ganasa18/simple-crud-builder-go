package server

import (
	appconfig "github.com/Ganasa18/simple-crud-builder-go/config"
	modelDynamicController "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/controller"
	tableDynamicController "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_table/controller"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/exception"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpServe struct {
	router                 *gin.Engine
	modelDynamicController modelDynamicController.DynamicModelController
	tableDynamicController tableDynamicController.DynamicTableController
}

func RunHttpServer(
	appConf *appconfig.Config,
	modelDynamicController modelDynamicController.DynamicModelController,
	tableDynamicController tableDynamicController.DynamicTableController,

) error {
	var hs HttpServe

	hs.router = gin.New()
	gin.SetMode(appConf.GinMode)

	// Global Exception Error Handler
	hs.router.Use(exception.ExceptionRecoveryMiddleware)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	hs.router.Use(cors.New(corsConfig))

	hs.router.SetTrustedProxies([]string{appConf.AppUrl})

	hs.modelDynamicController = modelDynamicController
	hs.tableDynamicController = tableDynamicController

	hs.setupRouter()

	return hs.router.Run(appConf.AppUrl + ":" + appConf.AppPort)
}
