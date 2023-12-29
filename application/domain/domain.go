package domain

import "time"

type CourseDomain struct {
	ID               int64
	Description      string
	Outline          string
	RegistrationDate time.Time
}
