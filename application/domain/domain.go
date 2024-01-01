package domain

import "time"

type CourseDomain struct {
	ID               int64
	Description      string
	Outline          string
	RegistrationDate time.Time
}

type StudentDomain struct {
	ID               int64
	Name             string
	RegistrationDate time.Time
}
