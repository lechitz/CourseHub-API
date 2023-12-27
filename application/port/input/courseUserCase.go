package input

import "github.com/lechitz/CourseHub-API/application/domain"

type ICourseService interface {
	Create(contextControl domain.ContextControl, course domain.Course) (domain.Course, error)
}
