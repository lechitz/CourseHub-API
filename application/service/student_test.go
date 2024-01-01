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

func TestStudentService_Create(t *testing.T) {

	tests := []struct {
		Name                            string
		Student                         domain.StudentDomain
		StudentDomainDataBaseRepository output.IStudentDomainDataBaseRepository
		StudentDomainCacheRepository    output.IStudentDomainCacheRepository
		ExpectedResult                  domain.StudentDomain
		ExpectedError                   error
	}{
		{
			Name: "success to save student",
			Student: domain.StudentDomain{
				Name: "Silvana",
			},
			StudentDomainDataBaseRepository: output.StudentDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error) {
					return domain.StudentDomain{
						ID:   1,
						Name: "Silvana",
					}, nil
				},
			},
			StudentDomainCacheRepository: output.StudentDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.StudentDomain{
				ID:   1,
				Name: "Silvana",
			},
			ExpectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			studentService := StudentService{
				LoggerSugar:                     loggerSugar,
				StudentDomainCacheRepository:    test.StudentDomainCacheRepository,
				StudentDomainDataBaseRepository: test.StudentDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			student, err := studentService.Create(contextControl, test.Student)
			assert.Equal(t, test.ExpectedResult, student)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}

func TestStudentService_GetById(t *testing.T) {

	tests := []struct {
		Name                            string
		Student                         domain.StudentDomain
		StudentDomainDataBaseRepository output.IStudentDomainDataBaseRepository
		StudentDomainCacheRepository    output.IStudentDomainCacheRepository
		ExpectedResult                  domain.StudentDomain
		ExpectedExists                  bool
		ExpectedError                   error
	}{
		{
			Name: "success to get a student by id",
			Student: domain.StudentDomain{
				Name: "Silvana",
			},
			StudentDomainDataBaseRepository: output.StudentDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error) {
					return domain.StudentDomain{
						ID:   1,
						Name: "Silvana",
					}, true, nil
				},
			},
			StudentDomainCacheRepository: output.StudentDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.StudentDomain{
				ID:   1,
				Name: "Silvana",
			},
			ExpectedExists: true,
			ExpectedError:  nil,
		},
		{
			Name: "Student not found",
			Student: domain.StudentDomain{
				Name: "Silvana",
			},
			StudentDomainDataBaseRepository: output.StudentDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error) {
					return domain.StudentDomain{}, false, nil
				},
			},
			StudentDomainCacheRepository: output.StudentDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.StudentDomain{},
			ExpectedExists: false,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			studentService := StudentService{
				LoggerSugar:                     loggerSugar,
				StudentDomainCacheRepository:    test.StudentDomainCacheRepository,
				StudentDomainDataBaseRepository: test.StudentDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			student, exists, err := studentService.GetByID(contextControl, 1)
			assert.Equal(t, test.ExpectedResult, student)
			assert.Equal(t, test.ExpectedExists, exists)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}

func TestStudentService_GetStudents(t *testing.T) {

	tests := []struct {
		Name                            string
		Students                        []domain.StudentDomain
		StudentDomainDataBaseRepository output.IStudentDomainDataBaseRepository
		StudentDomainCacheRepository    output.IStudentDomainCacheRepository
		ExpectedResult                  []domain.StudentDomain
		ExpectedError                   error
	}{
		{
			Name: "success to get students",
			Students: []domain.StudentDomain{
				{
					ID:   1,
					Name: "Silvana",
				},
				{
					ID:   2,
					Name: "Silvana",
				},
			},
			StudentDomainDataBaseRepository: output.StudentDomainDataBaseRepositoryMock{
				GetStudentsMock: func(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error) {
					return []domain.StudentDomain{
						{
							ID:   1,
							Name: "Silvana",
						},
						{
							ID:   2,
							Name: "Silvana",
						},
					}, nil
				},
			},
			StudentDomainCacheRepository: output.StudentDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: []domain.StudentDomain{
				{
					ID:   1,
					Name: "Silvana",
				},
				{
					ID:   2,
					Name: "Silvana",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:     "error to get students",
			Students: []domain.StudentDomain{},
			StudentDomainDataBaseRepository: output.StudentDomainDataBaseRepositoryMock{
				GetStudentsMock: func(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error) {
					return []domain.StudentDomain{}, fmt.Errorf("error to get students")
				},
			},
			StudentDomainCacheRepository: output.StudentDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: []domain.StudentDomain{},
			ExpectedError:  fmt.Errorf("error to get students"),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			studentService := StudentService{
				LoggerSugar:                     loggerSugar,
				StudentDomainDataBaseRepository: test.StudentDomainDataBaseRepository,
				StudentDomainCacheRepository:    test.StudentDomainCacheRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			students, err := studentService.GetStudents(contextControl, test.Students)
			assert.Equal(t, test.ExpectedResult, students)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}
