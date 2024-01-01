package http

import (
	"github.com/go-chi/chi"
	"github.com/lechitz/CourseHub-API/adapter/input/http/handler"
	"go.uber.org/zap"
)

type Router struct {
	ContextPath string
	chiRouter   chi.Router
	LoggerSugar *zap.SugaredLogger
}

func GetNewRouter(loggerSugar *zap.SugaredLogger) Router {
	router := chi.NewRouter()
	return Router{
		chiRouter:   router,
		LoggerSugar: loggerSugar,
	}
}

func (router Router) GetChiRouter() chi.Router {
	return router.chiRouter
}

func (router Router) AddGroupHandlerHealthCheck(ah *handler.Generic) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/health-check", func(r chi.Router) {
			r.Get("/", ah.HealthCheck)
		})
	}
}

func (router Router) AddGroupHandlerCourse(ah *handler.Course) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/course", func(r chi.Router) {
			r.Post("/create", ah.Create)
			r.Get("/{id}", ah.GetByID)
			r.Get("/courses", ah.GetCourses)
		})
	}
}

func (router Router) AddGroupHandlerStudent(ah *handler.Student) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/student", func(r chi.Router) {
			r.Post("/create", ah.Create)
			r.Get("/{id}", ah.GetByID)
			r.Get("/students", ah.GetStudents)
		})
	}
}
