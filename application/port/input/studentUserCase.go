package input

import "github.com/lechitz/CourseHub-API/application/domain"

type IStudentService interface {
	Create(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error)
	GetStudents(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error)
}
