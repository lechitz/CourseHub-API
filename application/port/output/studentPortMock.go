package output

import (
	"github.com/lechitz/CourseHub-API/application/domain"
	"time"
)

type StudentDomainDataBaseRepositoryMock struct {
	SaveMock        func(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error)
	GetByIDMock     func(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error)
	GetStudentsMock func(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error)
}

type StudentDomainCacheRepositoryMock struct {
	SetMock    func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	GetMock    func(contextControl domain.ContextControl, key string) (string, error)
	DeleteMock func(contextControl domain.ContextControl, key string) error
}

func (c StudentDomainDataBaseRepositoryMock) Save(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error) {
	if c.SaveMock != nil {
		return c.SaveMock(contextControl, student)
	}
	return domain.StudentDomain{}, nil
}

func (c StudentDomainDataBaseRepositoryMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.StudentDomain, bool, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.StudentDomain{}, false, nil
}

func (c StudentDomainDataBaseRepositoryMock) GetStudents(contextControl domain.ContextControl, students []domain.StudentDomain) ([]domain.StudentDomain, error) {
	if c.GetStudentsMock != nil {
		return c.GetStudentsMock(contextControl, students)
	}
	return []domain.StudentDomain{}, nil
}

func (c StudentDomainCacheRepositoryMock) Delete(contextControl domain.ContextControl, key string) error {
	if c.DeleteMock != nil {
		return c.DeleteMock(contextControl, key)
	}
	return nil
}

func (c StudentDomainCacheRepositoryMock) Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
	if c.SetMock != nil {
		return c.SetMock(contextControl, key, hash, expirationTime)
	}
	return nil
}

func (c StudentDomainCacheRepositoryMock) Get(contextControl domain.ContextControl, key string) (string, error) {
	if c.GetMock != nil {
		return c.GetMock(contextControl, key)
	}
	return "", nil
}
