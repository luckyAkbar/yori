package console

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luckyAkbar/yori/internal/config"
	"github.com/luckyAkbar/yori/internal/db"
	"github.com/luckyAkbar/yori/internal/delivery"
	"github.com/luckyAkbar/yori/internal/repository"
	"github.com/luckyAkbar/yori/internal/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runServer = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  "Use this command to start Himatro API HTTP server",
	Run:   InitServer,
}

func init() {
	RootCmd.AddCommand(runServer)
}

func InitServer(_ *cobra.Command, _ []string) {
	db.InitializePostgresConn()

	sqlDB, err := db.PostgresDB.DB()
	if err != nil {
		logrus.Fatal("unable to start server. reason: ", err.Error())
	}

	defer sqlDB.Close()

	recordRepo := repository.NewRecordRepository(db.PostgresDB)
	recordUsecase := usecase.NewRecordUsecase(recordRepo)

	HTTPServer := echo.New()

	HTTPServer.Pre(middleware.AddTrailingSlash())
	HTTPServer.Use(middleware.Logger())

	RESTGroup := HTTPServer.Group("rest")

	delivery.InitService(recordUsecase, RESTGroup)

	if err := HTTPServer.Start(fmt.Sprintf(":%s", config.ServerPort())); err != nil {
		logrus.Fatal("unable to start server. reason: ", err.Error())
	}
}
