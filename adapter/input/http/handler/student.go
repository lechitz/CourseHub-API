package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jinzhu/copier"
	"strconv"

	"github.com/lechitz/CourseHub-API/application/domain"
	"github.com/lechitz/CourseHub-API/application/port/input"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	SuccessToCreateStudent = "student created with success"
	SuccessToGetStudent    = "student found with success"
	ErrorToCreateStudent   = "error to create and process the request"
	ErrorToGetStudent      = "error to get student by id"
	StudentNotFound        = "student with id %d wasn´t found"
	// FieldNameError           = "Field name cannot be empty"
	ErrorToGetListOfStudents = "error to get list of students"
	SuccessToListStudents    = "list of students retrieved successfully"
)

type Student struct {
	StudentService input.IStudentService
	LoggerSugar    *zap.SugaredLogger
}

type StudentRequest struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	RegistrationDate time.Time `json:"registration_date"`
}

type StudentResponse struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	RegistrationDate time.Time `json:"registration_date"`
}

func (c *Student) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var studentRequest StudentRequest
	json.NewDecoder(r.Body).Decode(&studentRequest)

	var studentDomain domain.StudentDomain
	copier.Copy(&studentDomain, &studentRequest)

	if studentRequest.Name == "" {
		c.LoggerSugar.Errorw(ErrorToCreateCourse, "error", FieldOutlineError)
		response := objectResponse(ErrorToCreateCourse, FieldOutlineError)
		responseReturn(w, http.StatusBadRequest, response.Bytes())
		return
	}

	studentDomain, err := c.StudentService.Create(contextControl, studentDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreateStudent, "error", err.Error())
		response := objectResponse(ErrorToCreateStudent, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var studentResponse StudentResponse
	copier.Copy(&studentResponse, &studentDomain)
	response := objectResponse(studentResponse, SuccessToCreateStudent)
	responseReturn(w, http.StatusCreated, response.Bytes())
}

func (c *Student) GetByID(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var IDRequest, err = strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetStudent, "error", err.Error())
		response := objectResponse(ErrorToGetStudent, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	studentDomain, exists, err := c.StudentService.GetByID(contextControl, IDRequest)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetStudent, "error", err.Error())
		response := objectResponse(ErrorToGetStudent, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if !exists {
		c.LoggerSugar.Errorw(StudentNotFound)
		response := objectResponse(StudentNotFound, fmt.Sprintf(StudentNotFound, IDRequest))
		responseReturn(w, http.StatusNotFound, response.Bytes())
		return
	}

	var studentResponse StudentResponse
	copier.Copy(&studentResponse, &studentDomain)
	response := objectResponse(studentResponse, SuccessToGetStudent)
	responseReturn(w, http.StatusOK, response.Bytes())
}

func (c *Student) GetStudents(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var students []domain.StudentDomain

	students, err := c.StudentService.GetStudents(contextControl, students)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetListOfStudents, "error", err.Error())
		response := objectResponse(ErrorToGetListOfStudents, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var studentResponses []StudentResponse
	for _, studentDomain := range students {
		var studentResponse StudentResponse
		copier.Copy(&studentResponse, &studentDomain)
		studentResponses = append(studentResponses, studentResponse)
	}

	response := objectResponse(studentResponses, SuccessToListStudents)
	responseReturn(w, http.StatusOK, response.Bytes())
}
