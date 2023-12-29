package handler

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/lechitz/CourseHub-API/application/domain"
	"github.com/lechitz/CourseHub-API/application/port/input"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	SuccessToCreateCourse = "success to create course"
	ErrorToCreateCourse   = "error to create and process the request"
)

type Course struct {
	CourseService input.ICourseService
	LoggerSugar   *zap.SugaredLogger
}

type CourseRequest struct {
	ID               int64     `json:"id"`
	Description      string    `json:"description"`
	Outline          string    `json:"outline"`
	RegistrationDate time.Time `json:"registration_date"`
}

type CourseResponse struct {
	ID               int64     `json:"id"`
	Description      string    `json:"description"`
	Outline          string    `json:"outline"`
	RegistrationDate time.Time `json:"registration_date"`
}

func (c *Course) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var courseRequest CourseRequest
	json.NewDecoder(r.Body).Decode(&courseRequest)

	var courseDomain domain.CourseDomain
	copier.Copy(&courseDomain, &courseRequest)

	courseDomain, err := c.CourseService.Create(contextControl, courseDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreateCourse, "error", err.Error())
		response := objectResponse(ErrorToCreateCourse, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var courseResponse CourseResponse
	copier.Copy(&courseResponse, &courseDomain)
	response := objectResponse(courseResponse, SuccessToCreateCourse)
	responseReturn(w, http.StatusCreated, response.Bytes())
}
