package delivery

import (
	"net/http"

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
	s.group.GET("/find/ktp/:ktp/", s.handleFindByKtp())
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
			return c.JSON(http.StatusOK, record)
		}
	}
}
