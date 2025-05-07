package cmd

import (
	"context"
	"log"
	"smart-scene-app-api/common"
	"smart-scene-app-api/server"
	logger2 "smart-scene-app-api/services/logger"
	postgres3 "smart-scene-app-api/services/postgres"
	"smart-scene-app-api/services/rest_api_service"

	"github.com/spf13/cobra"
)

var restApiServiceCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Your Server",
	Long:  "Let start a server with your opinion",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		common.FetchMasterErrData()

		loggerPkg := logger2.NewLogger("Logger")
		if err := loggerPkg.Run(); err != nil {
			log.Panic(err)
		}

		logger := loggerPkg.Get()

		start, _ := cmd.Flags().GetBool("start")

		if start {
			svr := server.NewServer("SupplierLoyaltyService", 8080)
			restHdl := rest_api_service.RestHandler(svr)
			err, postgres := postgres3.NewMainPostgres(common.PREFIX_MAIN_POSTGRES)
			if err != nil {
				logger.Error().Println("NewMainPostgres", err)
				return
			}

			svr.AddLogger(logger)
			svr.InitContext(ctx)
			svr.InitService(postgres)
			svr.AddHandler(restHdl)
			if err := svr.Run(); err != nil {
				logger.Error().Printf("Server is stopped by %v", err.Error())
			}
		}
	},
}
