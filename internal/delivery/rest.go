package delivery

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/yori/internal/model"
	"github.com/luckyAkbar/yori/internal/usecase"
	"github.com/sirupsen/logrus"
)

type Service struct {
	recordUsecase model.RecordUsecase
	fileUsecase   model.FileUsecase
	roUsecase     model.ROUsecase

	group *echo.Group
}

func InitService(recordUsecase model.RecordUsecase, fileUsecase model.FileUsecase, roUsecase model.ROUsecase, group *echo.Group) {
	s := &Service{
		recordUsecase,
		fileUsecase,
		roUsecase,
		group,
	}

	s.initRoutes()
}

func (s *Service) initRoutes() {
	s.group.GET("/", s.handleIndex())
	s.group.GET("/ro/", s.handleROChecking())
	s.group.POST("/ro/check/", s.handleAdvanceROChecking())
	s.group.GET("/ro/result/", s.handleGetResultPage())
	s.group.GET("/ro/result/:id/", s.handleGetResult())
	s.group.GET("/ro/result/:id/csv/", s.handleGetCSVResult())
	s.group.GET("/dealer/", s.handleCheckDealer())
	s.group.GET("/upload/", s.handleUpload())
	s.group.POST("/upload/", s.handleUploadNewCSVData())
	s.group.GET("/find/ktp/:ktp/", s.handleFindByKtp())
	s.group.GET("/find/kk/:kk/", s.handleFindByKK())
}

func (s *Service) handleIndex() echo.HandlerFunc {
	return func(c echo.Context) error {
		filepath := path.Join("internal/views", "index.html")
		tmpl, err := template.ParseFiles(filepath)
		if err != nil {
			return ErrInternal
		}

		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return ErrInternal
		}

		return nil
	}
}

func (s *Service) handleGetResultPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		notFoundPage := path.Join("internal/views", "result.html")
		tmpl, err := template.ParseFiles(notFoundPage)

		if err != nil {
			return ErrInternal
		}

		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return ErrInternal
		}

		return nil
	}
}

func (s *Service) handleGetResult() echo.HandlerFunc {
	type RES struct {
		Result    []model.AdvanceROCheckingResult `json:"result"`
		TotalBase int64                           `json:"total_base"`
		TotalRO   int64                           `json:"total_ro"`
	}
	return func(c echo.Context) error {
		param := c.Param("id")
		if param == "" {
			return ErrBadRequest
		}

		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrBadRequest
		}

		result, res, err := s.roUsecase.GetResult(c.Request().Context(), id)
		switch err {
		default:
			logrus.Error(err)
			return ErrInternal
		case usecase.ErrNotFound:
			return ErrNotFound
		case usecase.ErrInProgress:
			return c.NoContent(http.StatusServiceUnavailable)
		case nil:
			return c.JSON(http.StatusOK, &RES{
				Result:    result,
				TotalBase: res.TotalBase,
				TotalRO:   res.TotalRO,
			})
		}
	}
}

func (s *Service) handleUpload() echo.HandlerFunc {
	return func(c echo.Context) error {
		filepath := path.Join("internal/views", "upload.html")
		tmpl, err := template.ParseFiles(filepath)
		if err != nil {
			return ErrInternal
		}

		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return ErrInternal
		}

		return nil
	}
}

func (s *Service) handleROChecking() echo.HandlerFunc {
	return func(c echo.Context) error {
		filepath := path.Join("internal/views", "ro.html")
		tmpl, err := template.ParseFiles(filepath)
		if err != nil {
			return ErrInternal
		}

		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return ErrInternal
		}

		return nil
	}
}

func (s *Service) handleFindByKtp() echo.HandlerFunc {
	return func(c echo.Context) error {
		ktp := c.Param("ktp")
		if ktp == "" {
			return ErrBadRequest
		}

		record, err := s.recordUsecase.FindByKTP(c.Request().Context(), ktp)
		switch err {
		default:
			return ErrInternal
		case usecase.ErrNotFound:
			return ErrNotFound
		case nil:
			if len(record) == 0 {
				return ErrNotFound
			}

			return c.JSON(http.StatusOK, record)
		}
	}
}

func (s *Service) handleFindByKK() echo.HandlerFunc {
	return func(c echo.Context) error {
		kk := c.Param("kk")
		if kk == "" {
			return ErrBadRequest
		}

		record, err := s.recordUsecase.FindByKK(c.Request().Context(), kk)
		switch err {
		default:
			return ErrInternal
		case usecase.ErrNotFound:
			return ErrNotFound
		case nil:
			if len(record) == 0 {
				return ErrNotFound
			}

			return c.JSON(http.StatusOK, record)
		}
	}
}

func (s *Service) handleUploadNewCSVData() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			logrus.Info(err)
			return ErrBadRequest
		}

		err = s.fileUsecase.Upload(c.Request().Context(), file)
		switch err {
		default:
			logrus.Error(err)
			return ErrInternal
		case usecase.ErrBadRequest:
			return ErrBadRequest
		case nil:
			return c.NoContent(http.StatusOK)
		}
	}
}

func (s *Service) handleCheckDealer() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			return ErrBadRequest
		}

		err := s.roUsecase.CheckIsDealerExists(c.Request().Context(), strings.Trim(name, " "))
		switch err {
		default:
			logrus.Error(err)
			return ErrInternal
		case usecase.ErrNotFound:
			return ErrNotFound
		case nil:
			return c.NoContent(http.StatusOK)
		}
	}
}

func (s *Service) handleAdvanceROChecking() echo.HandlerFunc {
	type Response struct {
		ID string `json:"id"`
	}
	return func(c echo.Context) error {
		input := &model.AdvanceROCheckingInput{}
		if err := c.Bind(input); err != nil {
			return ErrBadRequest
		}

		id := utils.GenerateID()
		logrus.Info("ini id awal", id)

		go func(id int64) {
			ctx := context.TODO()
			logrus.Info("ini id di dalem go fund", id)
			err := s.roUsecase.HandleAdvanceChecking(ctx, input, id)
			if err != nil {
				logrus.Error(err)
			}
		}(id)

		return c.JSON(http.StatusOK, &Response{
			ID: fmt.Sprintf("%d", id),
		})
	}
}

func (s *Service) handleGetCSVResult() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := c.Param("id")
		if param == "" {
			return ErrBadRequest
		}

		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrBadRequest
		}

		ro, res, err := s.roUsecase.GetResult(c.Request().Context(), id)
		switch err {
		default:
			logrus.Error(err)
			return ErrInternal
		case usecase.ErrNotFound:
			return ErrNotFound
		case usecase.ErrInProgress:
			return c.NoContent(http.StatusServiceUnavailable)
		case nil:
			break
		}

		b := &bytes.Buffer{}
		wr := csv.NewWriter(b)

		wr.Write([]string{"no", "ktp", "pembayaran", "nama_dealer", "jumlah_order"})

		logrus.Info(utils.Dump(ro[10]))

		for i, r := range ro {
			wr.Write([]string{fmt.Sprintf("%d", i+1), r.KTP, strings.Join(r.Payments, "-"), strings.Join(r.DealerNames, "-"), fmt.Sprintf("%d", r.NumOrders)})
		}

		wr.Flush()

		if err := wr.Error(); err != nil {
			logrus.Error(err)
			return ErrInternal
		}

		logrus.Info(res.TotalBase)

		filename := fmt.Sprintf("%d-%d-%d", res.ID, res.TotalBase, res.TotalRO)

		c.Response().Header().Set("Content-Type", "text/csv")
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))
		c.Response().Write(b.Bytes())

		return c.NoContent(http.StatusLocked)
	}
}
