package console

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/luckyAkbar/yori/internal/db"
	"github.com/luckyAkbar/yori/internal/helper"
	"github.com/luckyAkbar/yori/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var seederCmd = &cobra.Command{
	Use:  "seeder",
	Long: "seed the database using the data.csv file in the root dir",
	Run:  seeder,
}

func init() {
	RootCmd.AddCommand(seederCmd)
}

func seeder(cmd *cobra.Command, args []string) {
	file, err := os.Open("data.csv")
	if err != nil {
		logrus.Error("failed to open data.csv: ", err)
		os.Exit(1)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		logrus.Error("failed to read data from csv file: ", err)
		os.Exit(1)
	}

	if err := helper.ValidateCSVHeader(data[0]); err != nil {
		logrus.Error("invalid header: ", err)
		os.Exit(1)
	}

	records, err := helper.GenerateAndValidateRecords(data)
	if err != nil {
		logrus.Error("failed to generate records: ", err)
		os.Exit(1)
	}

	db.InitializePostgresConn()

	recordRepo := repository.NewRecordRepository(db.PostgresDB)

	if err := recordRepo.SaveBulk(context.TODO(), records); err != nil {
		logrus.Error("failed to save: ", err)
		os.Exit(1)
	}

	logrus.Info("finished saving")
	os.Exit(0)
}
