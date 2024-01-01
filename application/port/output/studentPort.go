package output

import (
	"github.com/lechitz/CourseHub-API/application/domain"
	"time"
)

type IStudentDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, course domain.StudentDomain) (domain.StudentDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error)
	GetStudents(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error)
}

type IStudentDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
