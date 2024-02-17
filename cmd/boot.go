package cmd

import (
	appconfig "github.com/Ganasa18/simple-crud-builder-go/config"
	server "github.com/Ganasa18/simple-crud-builder-go/http"
	modelDynamicController "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/controller"
	modelDynamicRepository "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/repository"
	modelDynamicService "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/service"
	tableDynamicController "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_table/controller"
	tableDynamicRepository "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_table/repository"
	tableDynamicService "github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_table/service"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func initHTTP() error {

	appConf := appconfig.InitAppConfig()

	var gConfig *gorm.Config = &gorm.Config{}
	db, err := appconfig.NewDatabase(appConf, gConfig)

	validate := validator.New()

	if err != nil {
		utils.IsErrorDoPanic(err)
	}

	// define module model dynamic
	modelDynamicRepo := modelDynamicRepository.NewDynamicModelRepository(db)
	modelDynamicSvc := modelDynamicService.NewDynamicModelService(modelDynamicRepo, validate)
	modelDynamicCtrl := modelDynamicController.NewDynamicModelController(modelDynamicSvc)

	// defin module model dynamic table
	tableDynamicRepo := tableDynamicRepository.NewDynamicTableRepository(db)
	tableDynamicSvc := tableDynamicService.NewDynamicTableService(tableDynamicRepo, validate)
	tableDynamicCtrl := tableDynamicController.NewDynamicTableController(tableDynamicSvc)

	// run server
	err = server.RunHttpServer(appConf, modelDynamicCtrl, tableDynamicCtrl)
	if err != nil {
		utils.IsErrorDoPanic(err)
	}

	return nil
}
