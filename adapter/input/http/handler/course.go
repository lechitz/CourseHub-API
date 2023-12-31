package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jinzhu/copier"
	"github.com/lechitz/CourseHub-API/application/domain"
	"github.com/lechitz/CourseHub-API/application/port/input"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const (
	SuccessToCreateCourse = "course created with success"
	SuccessToGetCourse    = "course found with success"
	ErrorToCreateCourse   = "error to create and process the request"
	ErrorToGetCourse      = "error to get course by id"
	CourseNotFound        = "course with id %d wasn´t found"
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

func (c *Course) GetByID(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var IDRequest, err = strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetCourse, "error", err.Error())
		response := objectResponse(ErrorToGetCourse, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	courseDomain, exists, err := c.CourseService.GetByID(contextControl, IDRequest)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetCourse, "error", err.Error())
		response := objectResponse(ErrorToGetCourse, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if !exists {
		c.LoggerSugar.Errorw(CourseNotFound)
		response := objectResponse(CourseNotFound, fmt.Sprintf(CourseNotFound, IDRequest))
		responseReturn(w, http.StatusNotFound, response.Bytes())
		return
	}

	var courseResponse CourseResponse
	copier.Copy(&courseResponse, &courseDomain)
	response := objectResponse(courseResponse, SuccessToGetCourse)
	responseReturn(w, http.StatusOK, response.Bytes())
}

func (c *Course) GetCourses(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	// Criar um slice de cursos para passar para o serviço
	var courses []domain.CourseDomain

	// Chamar o método GetCourses do serviço corrigido
	courses, err := c.CourseService.GetCourses(contextControl, courses)
	if err != nil {
		c.LoggerSugar.Errorw("error to get list of courses", "error", err.Error())
		response := objectResponse("error to get list of courses", err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var courseResponses []CourseResponse
	for _, courseDomain := range courses {
		var courseResponse CourseResponse
		copier.Copy(&courseResponse, &courseDomain)
		courseResponses = append(courseResponses, courseResponse)
	}

	response := objectResponse(courseResponses, "list of courses retrieved successfully")
	responseReturn(w, http.StatusOK, response.Bytes())
}
