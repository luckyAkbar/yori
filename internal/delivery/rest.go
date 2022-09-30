package delivery

import (
	"html/template"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/yori/internal/model"
	"github.com/luckyAkbar/yori/internal/usecase"
)

type Service struct {
	recordUsecase model.RecordUsecase

	group *echo.Group
}

func InitService(recordUsecase model.RecordUsecase, group *echo.Group) {
	s := &Service{
		recordUsecase,
		group,
	}

	s.initRoutes()
}

func (s *Service) initRoutes() {
	s.group.GET("/ktp/", s.handleKTP())
	s.group.GET("/kk/", s.handleKK())
	s.group.GET("/find/ktp/:ktp/", s.handleFindByKtp())
	s.group.GET("/find/kk/:kk/", s.handleFindByKK())
}

func (s *Service) handleKTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		filepath := path.Join("internal/views", "ktp.html")
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

func (s *Service) handleKK() echo.HandlerFunc {
	return func(c echo.Context) error {
		filepath := path.Join("internal/views", "kk.html")
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
