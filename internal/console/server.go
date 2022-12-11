package console

import (
	"fmt"
	"time"

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
	resultRepo := repository.NewResultRepo(db.PostgresDB)

	fileUsecase := usecase.NewFileUsecase(recordRepo)
	roUsecase := usecase.NewROUsecase(recordRepo, resultRepo)

	HTTPServer := echo.New()

	HTTPServer.Pre(middleware.AddTrailingSlash())
	HTTPServer.Use(middleware.Logger())
	HTTPServer.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      time.Hour * 5,
		ErrorMessage: "operation timed out",
	}))

	HTTPServer.Static("/assets/", "internal/assets")

	RESTGroup := HTTPServer.Group("index")

	delivery.InitService(recordUsecase, fileUsecase, roUsecase, RESTGroup)

	if err := HTTPServer.Start(fmt.Sprintf(":%s", config.ServerPort())); err != nil {
		logrus.Fatal("unable to start server. reason: ", err.Error())
	}
}
