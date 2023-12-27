package service

import (
	"github.com/lechitz/CourseHub-API/application/domain"
	"go.uber.org/zap"
)

type CourseService struct {
	LoggerSugar *zap.SugaredLogger
}

func (service CourseService) Create(contextControl domain.ContextControl, customer domain.Course) (domain.Course, error) {

	return domain.Course{}, nil
}
