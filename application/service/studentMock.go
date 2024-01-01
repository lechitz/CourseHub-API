package service

import "github.com/lechitz/CourseHub-API/application/domain"

type StudentMock struct {
	CreateMock  func(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error)
	GetByIDMock func(ID int64) (domain.StudentDomain, error)
}

func (c StudentMock) Create(contextControl domain.ContextControl, student domain.StudentDomain) (domain.StudentDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, student)
	}
	return domain.StudentDomain{}, nil
}

func (c StudentMock) GetByID(ID int64) (domain.StudentDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(ID)
	}
	return domain.StudentDomain{}, nil
}
