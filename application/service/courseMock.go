package service

import "github.com/lechitz/CourseHub-API/application/domain"

type CourseMock struct {
	CreateMock  func(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error)
	GetByIDMock func(ID int64) (domain.CourseDomain, error)
}

func (c CourseMock) Create(contextControl domain.ContextControl, course domain.CourseDomain) (domain.CourseDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, course)
	}
	return domain.CourseDomain{}, nil
}

func (c CourseMock) GetByID(ID int64) (domain.CourseDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(ID)
	}
	return domain.CourseDomain{}, nil
}
