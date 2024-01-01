package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	adpterHttpInput "github.com/lechitz/CourseHub-API/adapter/input/http"
	"github.com/lechitz/CourseHub-API/adapter/input/http/handler"
	"github.com/lechitz/CourseHub-API/adapter/output/cache"
	"github.com/lechitz/CourseHub-API/adapter/output/database"
	"github.com/lechitz/CourseHub-API/application/service"
	"github.com/lechitz/CourseHub-API/configuration/environment"
	"github.com/lechitz/CourseHub-API/configuration/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"

	"github.com/go-chi/chi"
	"os"
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

}

func main() {

	redisCache := cache.NewRedis(loggerSugar)

	postgresConnectionDB := repository.NewPostgresDB(environment.Setting.Postgres.DBUser, environment.Setting.Postgres.DBPassword,
		environment.Setting.Postgres.DBName, environment.Setting.Postgres.DBHost, environment.Setting.Postgres.DBPort, loggerSugar)

	coursePostgresDB := database.NewCoursePostgresDB(postgresConnectionDB, loggerSugar)
	studentPostgresDB := database.NewStudentPostgresDB(postgresConnectionDB, loggerSugar)

	genericHandler := &handler.Generic{
		LoggerSugar: loggerSugar,
	}

	courseService := service.CourseService{
		LoggerSugar:                    loggerSugar,
		CourseDomainDataBaseRepository: &coursePostgresDB,
		CourseDomainCacheRepository:    &redisCache,
	}

	courseHandler := &handler.Course{
		CourseService: &courseService,
		LoggerSugar:   loggerSugar,
	}

	studentService := service.StudentService{
		LoggerSugar:                     loggerSugar,
		StudentDomainDataBaseRepository: &studentPostgresDB,
		StudentDomainCacheRepository:    &redisCache,
	}

	studentHandler := &handler.Student{
		StudentService: &studentService,
		LoggerSugar:    loggerSugar,
	}

	contextPath := environment.Setting.Server.Context
	newRouter := adpterHttpInput.GetNewRouter(loggerSugar)
	newRouter.GetChiRouter().Route(fmt.Sprintf("/%s", contextPath), func(r chi.Router) {
		r.NotFound(genericHandler.NotFound)
		r.Group(newRouter.AddGroupHandlerHealthCheck(genericHandler))
		r.Group(newRouter.AddGroupHandlerCourse(courseHandler))
		r.Group(newRouter.AddGroupHandlerStudent(studentHandler))
	})

	serverHttp := &http.Server{
		Addr:           fmt.Sprintf(":%s", environment.Setting.Server.Port),
		Handler:        newRouter.GetChiRouter(),
		ReadTimeout:    environment.Setting.Server.ReadTimeout,
		WriteTimeout:   environment.Setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	loggerSugar.Infow("server started", "port", serverHttp.Addr,
		"contextPath", contextPath)

	if err := serverHttp.ListenAndServe(); err != nil {
		loggerSugar.Errorw("error to listen and starts server", "port", serverHttp.Addr,
			"contextPath", contextPath, "err", err.Error())
		panic(err.Error())
	}

}
