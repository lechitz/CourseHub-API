package output

import (
	"github.com/lechitz/CourseHub-API/application/domain"
	"time"
)

type CourseDomainDataBaseRepositoryMock struct {
	SaveMock    func(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error)
	GetByIDMock func(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, error)
}

type CourseDomainCacheRepositoryMock struct {
	SetMock    func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	GetMock    func(contextControl domain.ContextControl, key string) (string, error)
	DeleteMock func(contextControl domain.ContextControl, key string) error
}

func (c CourseDomainDataBaseRepositoryMock) Save(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error) {
	if c.SaveMock != nil {
		return c.SaveMock(contextControl, course)
	}
	return domain.CourseDomain{}, nil
}

func (c CourseDomainDataBaseRepositoryMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.CourseDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.CourseDomain{}, nil
}

func (c CourseDomainCacheRepositoryMock) Delete(contextControl domain.ContextControl, key string) error {
	if c.DeleteMock != nil {
		return c.DeleteMock(contextControl, key)
	}
	return nil
}

func (c CourseDomainCacheRepositoryMock) Get(contextControl domain.ContextControl, key string) (string, error) {
	if c.GetMock != nil {
		return c.GetMock(contextControl, key)
	}
	return "", nil
}

func (c CourseDomainCacheRepositoryMock) Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
	if c.SetMock != nil {
		return c.SetMock(contextControl, key, hash, expirationTime)
	}
	return nil
}
