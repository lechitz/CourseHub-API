package input

import "github.com/lechitz/CourseHub-API/application/domain"

type ICourseService interface {
	Create(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, bool, error)
	GetCourses(contextControl domain.ContextControl, courses []domain.CourseDomain) ([]domain.CourseDomain, error)
}
