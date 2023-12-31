package service

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/CourseHub-API/application/domain"
	"github.com/lechitz/CourseHub-API/application/port/output"
	"github.com/lechitz/CourseHub-API/configuration/environment"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"time"
)

var loggerSugar *zap.SugaredLogger

func init() {

	err := envconfig.Process("setting", &environment.Setting)
	if err != nil {
		panic(err.Error())
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	loggerSugar = logger.Sugar()
	loggerSugar.Infow("testing")

}

func TestCourseService_Create(t *testing.T) {

	tests := []struct {
		Name                           string
		Course                         domain.CourseDomain
		CourseDomainDataBaseRepository output.ICourseDomainDataBaseRepository
		CourseDomainCacheRepository    output.ICourseDomainCacheRepository
		ExpectedResult                 domain.CourseDomain
		ExpectedError                  error
	}{
		{
			Name: "success to save course",
			Course: domain.CourseDomain{
				Description: "Matematica",
				Outline:     "calculo 1",
			},
			CourseDomainDataBaseRepository: output.CourseDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error) {
					return domain.CourseDomain{
						ID:          1,
						Description: "Matematica",
						Outline:     "calculo 1",
					}, nil
				},
			},
			CourseDomainCacheRepository: output.CourseDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.CourseDomain{
				ID:          1,
				Description: "Matematica",
				Outline:     "calculo 1",
			},
			ExpectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			courseService := CourseService{
				LoggerSugar:                    loggerSugar,
				CourseDomainCacheRepository:    test.CourseDomainCacheRepository,
				CourseDomainDataBaseRepository: test.CourseDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			course, err := courseService.Create(contextControl, test.Course)
			assert.Equal(t, test.ExpectedResult, course)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}

func TestCourseService_GetById(t *testing.T) {

	tests := []struct {
		Name                           string
		Course                         domain.CourseDomain
		CourseDomainDataBaseRepository output.ICourseDomainDataBaseRepository
		CourseDomainCacheRepository    output.ICourseDomainCacheRepository
		ExpectedResult                 domain.CourseDomain
		ExpectedExists                 bool
		ExpectedError                  error
	}{
		{
			Name: "success to get a course by id",
			Course: domain.CourseDomain{
				Description: "Português",
				Outline:     "Tempos Verbais",
			},
			CourseDomainDataBaseRepository: output.CourseDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, bool, error) {
					return domain.CourseDomain{
						ID:          1,
						Description: "Português",
						Outline:     "Tempos Verbais",
					}, true, nil
				},
			},
			CourseDomainCacheRepository: output.CourseDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.CourseDomain{
				ID:          1,
				Description: "Português",
				Outline:     "Tempos Verbais",
			},
			ExpectedExists: true,
			ExpectedError:  nil,
		},
		{
			Name: "Course not found",
			Course: domain.CourseDomain{
				Description: "Portuguêss",
				Outline:     "Tempos Verbais",
			},
			CourseDomainDataBaseRepository: output.CourseDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, bool, error) {
					return domain.CourseDomain{}, false, nil
				},
			},
			CourseDomainCacheRepository: output.CourseDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.CourseDomain{},
			ExpectedExists: false,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			courseService := CourseService{
				LoggerSugar:                    loggerSugar,
				CourseDomainCacheRepository:    test.CourseDomainCacheRepository,
				CourseDomainDataBaseRepository: test.CourseDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			course, exists, err := courseService.GetByID(contextControl, 1)
			assert.Equal(t, test.ExpectedResult, course)
			assert.Equal(t, test.ExpectedExists, exists)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}

func TestCourseService_GetCourses(t *testing.T) {

	tests := []struct {
		Name                           string
		Courses                        []domain.CourseDomain
		CourseDomainDataBaseRepository output.ICourseDomainDataBaseRepository
		CourseDomainCacheRepository    output.ICourseDomainCacheRepository
		ExpectedResult                 []domain.CourseDomain
		ExpectedError                  error
	}{
		{
			Name: "success to get courses",
			Courses: []domain.CourseDomain{
				{
					ID:          1,
					Description: "Matematica",
					Outline:     "calculo 1",
				},
				{
					ID:          2,
					Description: "Português",
					Outline:     "Tempos Verbais",
				},
			},
			CourseDomainDataBaseRepository: output.CourseDomainDataBaseRepositoryMock{
				GetCoursesMock: func(contextControl domain.ContextControl, courses []domain.CourseDomain) ([]domain.CourseDomain, error) {
					return []domain.CourseDomain{
						{
							ID:          1,
							Description: "Matematica",
							Outline:     "calculo 1",
						},
						{
							ID:          2,
							Description: "Português",
							Outline:     "Tempos Verbais",
						},
					}, nil
				},
			},
			CourseDomainCacheRepository: output.CourseDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: []domain.CourseDomain{
				{
					ID:          1,
					Description: "Matematica",
					Outline:     "calculo 1",
				},
				{
					ID:          2,
					Description: "Português",
					Outline:     "Tempos Verbais",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:    "error to get courses",
			Courses: []domain.CourseDomain{},
			CourseDomainDataBaseRepository: output.CourseDomainDataBaseRepositoryMock{
				GetCoursesMock: func(contextControl domain.ContextControl, courses []domain.CourseDomain) ([]domain.CourseDomain, error) {
					return []domain.CourseDomain{}, fmt.Errorf("error to get courses")
				},
			},
			CourseDomainCacheRepository: output.CourseDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: []domain.CourseDomain{},
			ExpectedError:  fmt.Errorf("error to get courses"),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			courseService := CourseService{
				LoggerSugar:                    loggerSugar,
				CourseDomainDataBaseRepository: test.CourseDomainDataBaseRepository,
				CourseDomainCacheRepository:    test.CourseDomainCacheRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			courses, err := courseService.GetCourses(contextControl, test.Courses)
			assert.Equal(t, test.ExpectedResult, courses)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}
