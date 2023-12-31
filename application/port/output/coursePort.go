package output

import (
	"github.com/lechitz/CourseHub-API/application/domain"
	"time"
)

type ICourseDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, bool, error)
}

type ICourseDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
