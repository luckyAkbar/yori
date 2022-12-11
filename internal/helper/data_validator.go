package helper

import (
	"strconv"
	"strings"
	"time"

	"github.com/luckyAkbar/yori/internal/config"
	"github.com/luckyAkbar/yori/internal/model"
	"github.com/sirupsen/logrus"
)

func ValidateCSVHeader(header []string) error {

	validHeader := config.CSVHeader()

	if len(validHeader) != len(header) {
		return ErrInvalidHeaderErr
	}

	for i, h := range validHeader {
		if h != header[i] {
			return ErrInvalidHeaderErr
		}
	}

	return nil
}

func GenerateAndValidateRecords(data [][]string) ([]model.Record, error) {
	var records []model.Record
	validHeader := config.CSVHeader()

	for i, record := range data {
		if i == 0 { // skip the csv header
			continue
		}

		if len(record) != len(validHeader) {
			return records, ErrInvalidHeaderErr
		}

		no, err := strconv.ParseInt(trim(record[0]), 10, 64)
		if err != nil {
			logrus.Error("invalid no value: ", err)
			return records, ErrInvalidValueToRecords
		}

		date, err := parseDate(trim(record[3]))
		if err != nil {
			logrus.Info(record[3])
			logrus.Error("invalid value for tgl_mohon_faktur. It mush be date-able")

			return records, ErrInvalidValueToRecords
		}

		bulan, err := strconv.Atoi(trim(record[9]))
		if err != nil {
			logrus.Error("invalid bulan value: ", err)
			return records, ErrInvalidValueToRecords
		}

		tahun, err := strconv.Atoi(trim(record[10]))
		if err != nil {
			logrus.Error("invalid tahun value: ", err)
			return records, ErrInvalidValueToRecords
		}

		ktp := trim(record[6])
		if ktp == "" {
			logrus.Error("ktp tidak boleh kosong")
			return records, ErrInvalidValueToRecords
		}

		records = append(records, model.Record{
			No:             no,
			Nama:           trim(record[1]),
			NoEngine:       trim(record[2]),
			TglMohonFaktur: date,
			Fincoy:         trim(record[4]),
			Type:           trim(record[5]),
			Ktp:            ktp,
			Kk:             trim(record[7]),
			Dealer:         strings.ToLower(trim(record[8])),
			Bulan:          bulan,
			Tahun:          tahun,
		})
	}

	return records, nil
}

func RemoveDuplicateKTP(records []model.Record) []model.Record {
	var res []model.Record

	for _, record := range records {
		found := false
		for _, rec := range res {
			if record.Ktp == rec.Ktp {
				found = true
			}
		}

		if !found {
			res = append(res, record)
		}
	}

	return res
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func parseDate(s string) (time.Time, error) {
	parts := strings.Split(s, "-")
	y := parts[0]
	m := parts[1]
	d := parts[2]

	year, err := strconv.Atoi(y)
	if err != nil {
		logrus.Error(err)
		return time.Now(), err
	}

	month, err := strconv.Atoi(m)
	if err != nil {
		logrus.Error(err)
		return time.Now(), err
	}

	day, err := strconv.Atoi(d)
	if err != nil {
		logrus.Error(err)
		return time.Now(), err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}
